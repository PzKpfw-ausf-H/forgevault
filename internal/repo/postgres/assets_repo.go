package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetsSQLRepo struct {
	pool *pgxpool.Pool
}

func NewAssetsSQLRepo(pool *pgxpool.Pool) *AssetsSQLRepo {
	return &AssetsSQLRepo{
		pool: pool,
	}
}

func (ar *AssetsSQLRepo) Create(ctx context.Context, a domain.Asset) error {
	if err := domain.ValidateNewAsset(&a); err != nil {
		return err
	}
	tx, err := ar.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	_, txErr := tx.Exec(ctx,
		`INSERT INTO assets (id, title, description, type, author_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		a.ID,
		a.Title,
		a.Description,
		a.Type,
		a.AuthorID,
		a.CreatedAt,
		a.UpdatedAt,
	)
	if txErr != nil {
		return fmt.Errorf("insert into assets: %v", txErr)
	}

	if len(a.Tags) > 0 {
		b := &pgx.Batch{}
		for _, tag := range a.Tags {
			b.Queue(
				`INSERT INTO asset_tags (asset_id, tag) VALUES ($1, $2)`,
				a.ID,
				tag,
			)
		}

		br := tx.SendBatch(ctx, b)
		if err := br.Close(); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (ar *AssetsSQLRepo) GetByID(ctx context.Context, id domain.AssetID) (domain.Asset, error) {
	var asset *domain.Asset

	rows, err := ar.pool.Query(ctx,
		`SELECT a.id, a.title, a.description, a.type, a.author_id, a.created_at, a.updated_at, at.tag
		FROM assets a
		LEFT JOIN asset_tags at ON a.id = at.asset_id
		WHERE a.id = $1
		ORDER BY at.tag`,
		id,
	)
	if err != nil {
		return domain.Asset{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id          domain.AssetID
			title       string
			description string
			assetType   domain.AssetType
			author_id   domain.UserID
			created_at  time.Time
			updated_at  time.Time
			tagRaw      *string
		)

		if err := rows.Scan(&id, &title, &description, &assetType, &author_id, &created_at, &updated_at, &tagRaw); err != nil {
			return domain.Asset{}, fmt.Errorf("get asset scan: %v", err)
		}

		if asset == nil {
			asset = &domain.Asset{
				ID:          id,
				Title:       title,
				Description: description,
				Type:        assetType,
				AuthorID:    author_id,
				CreatedAt:   created_at,
				UpdatedAt:   updated_at,
				Tags:        []string{},
			}
		}

		if tagRaw != nil {
			asset.Tags = append(asset.Tags, *tagRaw)
		}
	}

	if err := rows.Err(); err != nil {
		return domain.Asset{}, fmt.Errorf("get asset rows: %v", err)
	}

	if asset == nil {
		return domain.Asset{}, repo.ErrNotFound
	}

	return *asset, nil
}

func (ar *AssetsSQLRepo) List(ctx context.Context, f repo.AssetFilter) ([]domain.Asset, error) {
	conds := make([]string, 0, 4)
	args := make([]any, 0, 8)
	n := 1

	base := `SELECT a.id, a.title, a.description, a.type, a.author_id, a.created_at, a.updated_at
	         FROM assets a`

	if f.Type != nil {
		conds = append(conds, fmt.Sprintf("a.type = $%d", n))
		args = append(args, string(*f.Type))
		n++
	}
	if f.AuthorID != nil {
		conds = append(conds, fmt.Sprintf("a.author_id = $%d", n))
		args = append(args, *f.AuthorID)
		n++
	}
	if f.TitleSub != nil {
		conds = append(conds, fmt.Sprintf("a.title ILIKE '%%' || $%d || '%%'", n))
		args = append(args, *f.TitleSub)
		n++
	}
	if f.Tag != nil {
		conds = append(conds, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM asset_tags at WHERE at.asset_id = a.id AND at.tag = $%d)", n,
		))
		args = append(args, *f.Tag)
		n++
	}

	sql := base
	if len(conds) > 0 {
		sql += " WHERE " + strings.Join(conds, " AND ")
	}
	sql += " ORDER BY a.created_at DESC, a.id DESC"

	limit := f.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}

	sql += fmt.Sprintf(" LIMIT $%d OFFSET $%d", n, n+1)
	args = append(args, limit, offset)

	rows, err := ar.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("list assets query: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Asset, 0, limit)
	ids := make([]domain.AssetID, 0, limit)

	for rows.Next() {
		var a domain.Asset
		if err := rows.Scan(&a.ID, &a.Title, &a.Description, &a.Type, &a.AuthorID, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("list assets scan: %w", err)
		}
		a.Tags = []string{}
		out = append(out, a)
		ids = append(ids, a.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list assets rows: %w", err)
	}

	if len(out) == 0 {
		return out, nil
	}

	tagRows, err := ar.pool.Query(ctx,
		`SELECT asset_id, tag
		   FROM asset_tags
		  WHERE asset_id = ANY($1)
		  ORDER BY asset_id, tag`,
		ids,
	)
	if err != nil {
		return nil, fmt.Errorf("list tags query: %w", err)
	}
	defer tagRows.Close()

	tagsByID := make(map[domain.AssetID][]string, len(ids))
	for tagRows.Next() {
		var assetID domain.AssetID
		var tag string
		if err := tagRows.Scan(&assetID, &tag); err != nil {
			return nil, fmt.Errorf("list tags scan: %w", err)
		}
		tagsByID[assetID] = append(tagsByID[assetID], tag)
	}
	if err := tagRows.Err(); err != nil {
		return nil, fmt.Errorf("list tags rows: %w", err)
	}

	for i := range out {
		out[i].Tags = tagsByID[out[i].ID]
	}

	return out, nil
}

func (ar *AssetsSQLRepo) Delete(ctx context.Context, id domain.AssetID) error {
	row, err := ar.pool.Exec(ctx,
		`DELETE FROM assets WHERE id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("delete asset: %v", err)
	}

	if row.RowsAffected() <= 0 {
		return repo.ErrNotFound
	}

	return nil
}

func (ar *AssetsSQLRepo) Update(ctx context.Context, a domain.Asset) error {
	if err := domain.ValidateNewAsset(&a); err != nil {
		return err
	}

	tx, err := ar.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	cmd, err := tx.Exec(ctx,
		`UPDATE assets SET title=$1, description=$2, type=$3, updated_at=$4 WHERE id=$5`,
		a.Title, a.Description, a.Type, a.UpdatedAt, a.ID,
	)
	if err != nil {
		return fmt.Errorf("update asset: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return repo.ErrNotFound
	}

	_, err = tx.Exec(ctx, `DELETE FROM asset_tags WHERE asset_id=$1`, a.ID)
	if err != nil {
		return fmt.Errorf("delete asset tags: %w", err)
	}

	if len(a.Tags) > 0 {
		b := &pgx.Batch{}
		for _, tag := range a.Tags {
			b.Queue(`INSERT INTO asset_tags (asset_id, tag) VALUES ($1,$2)`, a.ID, tag)
		}
		br := tx.SendBatch(ctx, b)
		if err := br.Close(); err != nil {
			return fmt.Errorf("insert asset tags: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update: %w", err)
	}
	return nil
}
