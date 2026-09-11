package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetRepository struct {
	db *pgxpool.Pool
}

var _ repo.AssetRepository = (*AssetRepository)(nil)

func NewAssetRepository(db *pgxpool.Pool) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(ctx context.Context, asset domain.Asset) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("postgres asset: create tx: %w", err)
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO assets (id, title, description, author_name, uploaded_by, type,
		processing_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		asset.ID, asset.Title, asset.Description, asset.AuthorName, asset.UploadedBy,
		asset.Type, asset.ProcessingStatus, asset.CreatedAt, asset.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("postgres asset: insert tx: %w", repo.ErrAlreadyExists)
			}
		}
		return fmt.Errorf("postgres asset: create: insert assets tx: %w", err)
	}

	batch := &pgx.Batch{}

	for _, tag := range asset.Tags {
		batch.Queue(
			`INSERT INTO asset_tags (asset_id, tag)
			VALUES ($1, $2)`,
			asset.ID,
			tag,
		)
	}

	br := tx.SendBatch(ctx, batch)

	for i := 0; i < len(asset.Tags); i++ {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("postgres asset: create: batch results: %w", err)
		}
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("postgres asset: create: batch close: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("postgres asset: create: tx commit: %w", err)
	}

	return nil
}

func (r *AssetRepository) GetByID(ctx context.Context, assetID domain.AssetID) (domain.Asset, error) {
	var a domain.Asset

	row := r.db.QueryRow(ctx,
		`SELECT id, title, description, author_name, uploaded_by, type,
		processing_status, created_at, updated_at
		FROM assets
		WHERE id = $1`,
		assetID,
	)

	if err := row.Scan(&a.ID, &a.Title, &a.Description, &a.AuthorName, &a.UploadedBy, &a.Type, &a.ProcessingStatus,
		&a.CreatedAt, &a.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Asset{}, fmt.Errorf("postgres asset: get by id asset: %w", repo.ErrNotFound)
		}

		return domain.Asset{}, fmt.Errorf("postgres asset: get by id asset: %w", err)
	}

	rows, err := r.db.Query(ctx,
		`SELECT tag
		FROM asset_tags
		WHERE asset_id = $1
		ORDER BY tag`,
		assetID,
	)

	if err != nil {
		return domain.Asset{}, fmt.Errorf("postgres asset: get tags by asset id: %w", err)
	}

	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return domain.Asset{}, fmt.Errorf("postgres asset: scan tag: %w", err)
		}

		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return domain.Asset{}, fmt.Errorf("postgres asset: close rows: %w", err)
	}

	a.Tags = tags

	return a, nil
}

func (r *AssetRepository) Update(ctx context.Context, asset domain.Asset) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("postgres asset: update: begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	cmd, err := tx.Exec(ctx,
		`UPDATE assets
		SET title = $1,
		description = $2,
		author_name = $3,
		type = $4,
		processing_status = $5,
		updated_at = $6
		WHERE id = $7`,
		asset.Title, asset.Description, asset.AuthorName, asset.Type,
		asset.ProcessingStatus, asset.UpdatedAt, asset.ID)
	if err != nil {
		return fmt.Errorf("postgres asset: update: update assets tx: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("postgres asset: update: %w", repo.ErrNotFound)
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM asset_tags
		WHERE asset_id = $1`,
		asset.ID)
	if err != nil {
		return fmt.Errorf("postgres asset: update: delete asset tags tx: %w", err)
	}

	batch := &pgx.Batch{}

	for _, tag := range asset.Tags {
		batch.Queue(
			`INSERT INTO asset_tags (asset_id, tag)
			VALUES ($1, $2)`,
			asset.ID,
			tag,
		)
	}

	br := tx.SendBatch(ctx, batch)

	for i := 0; i < len(asset.Tags); i++ {
		_, err := br.Exec()
		if err != nil {
			br.Close()
			return fmt.Errorf("postgres asset: update: batch result: %w", err)
		}
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("postgres asset: update: close batch result: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("postgres asset: update: tx commit: %w", err)
	}

	return nil
}

func (r *AssetRepository) Delete(ctx context.Context, id domain.AssetID) error {
	cmd, err := r.db.Exec(ctx,
		`DELETE FROM assets
		WHERE id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("postgres asset: delete: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("postgres asset: delete: %w", repo.ErrNotFound)
	}

	return nil
}

func (r *AssetRepository) List(ctx context.Context, params repo.AssetListParams) (repo.AssetListResult, error) {
	// getting the total number of matching assets
	// if got params -> where type = $1
	// if params is nil -> without where type = $1
	var total int64

	if params.Type != nil {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*)
			FROM assets
			WHERE type = $1`,
			*params.Type).Scan(&total)

		if err != nil {
			return repo.AssetListResult{}, fmt.Errorf("postgres asset: list: count: %w", err)
		}
	} else {
		err := r.db.QueryRow(ctx,
			`SELECT COUNT(*)
			FROM assets`,
		).Scan(&total)

		if err != nil {
			return repo.AssetListResult{}, fmt.Errorf("postgres asset: list: count: %w", err)
		}
	}

	// getting the assets page
	assets := make([]domain.Asset, 0)

	var (
		rows pgx.Rows
		err  error
	)

	if params.Type != nil {
		rows, err = r.db.Query(ctx,
			`SELECT id, title, description, author_name, uploaded_by, type,
			processing_status, created_at, updated_at
			FROM assets
			WHERE type = $1
			ORDER BY created_at DESC, id
			LIMIT $2 OFFSET $3`,
			*params.Type,
			params.Limit,
			params.Offset)
	} else {
		rows, err = r.db.Query(ctx,
			`SELECT id, title, description, author_name, uploaded_by, type,
			processing_status, created_at, updated_at
			FROM assets
			ORDER BY created_at DESC, id
			LIMIT $1 OFFSET $2`,
			params.Limit,
			params.Offset)
	}

	if err != nil {
		return repo.AssetListResult{}, fmt.Errorf("postgres asset: list: query assets: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a domain.Asset

		if err := rows.Scan(&a.ID, &a.Title, &a.Description, &a.AuthorName, &a.UploadedBy, &a.Type,
			&a.ProcessingStatus, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return repo.AssetListResult{}, fmt.Errorf("postgres asset: list: scan asset: %w", err)
		}

		assets = append(assets, a)
	}

	if err := rows.Err(); err != nil {
		return repo.AssetListResult{}, fmt.Errorf("postgres asset: list: rows: %w", err)
	}

	if len(assets) == 0 {
		return repo.AssetListResult{
			Assets: assets,
			Total:  total,
		}, nil
	}

	ids := make([]string, 0, len(assets))
	for _, a := range assets {
		ids = append(ids, string(a.ID))
	}

	tagRows, err := r.db.Query(ctx,
		`SELECT asset_id, tag
		FROM asset_tags
		WHERE asset_id = ANY($1)
		ORDER BY asset_id, tag`,
		ids,
	)

	if err != nil {
		return repo.AssetListResult{}, fmt.Errorf("postgres asset: list: query tags: %w", err)
	}

	defer tagRows.Close()

	tagsByAssetID := make(map[domain.AssetID][]string)

	for tagRows.Next() {
		var (
			assetID domain.AssetID
			tag     string
		)

		if err := tagRows.Scan(&assetID, &tag); err != nil {
			return repo.AssetListResult{}, fmt.Errorf("postgres asset: list: scan tag: %w", err)
		}

		tagsByAssetID[assetID] = append(tagsByAssetID[assetID], tag)
	}

	if err := tagRows.Err(); err != nil {
		return repo.AssetListResult{}, fmt.Errorf("postgres asset: list: tag rows: %w", err)
	}

	for i := range assets {
		assets[i].Tags = tagsByAssetID[assets[i].ID]
	}

	return repo.AssetListResult{
		Assets: assets,
		Total:  total,
	}, nil
}
