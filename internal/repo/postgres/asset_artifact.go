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

type AssetArtifactRepository struct {
	db *pgxpool.Pool
}

var _ repo.AssetArtifactRepository = (*AssetArtifactRepository)(nil)

func NewAssetArtifactRepository(db *pgxpool.Pool) *AssetArtifactRepository {
	return &AssetArtifactRepository{db: db}
}

func (r *AssetArtifactRepository) Create(ctx context.Context, artifact domain.AssetArtifact) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO asset_artifacts (id, asset_id, type, mime_type, size, checksum, storage_key, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		artifact.ID,
		artifact.AssetID,
		artifact.Type,
		artifact.MimeType,
		artifact.Size,
		artifact.Checksum,
		artifact.StorageKey,
		artifact.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("postgres asset artifact: create artifact: %w", repo.ErrAlreadyExists)
			}
		}

		return fmt.Errorf("postgres asset artifact: create artifact: %w", err)
	}

	return nil
}

func (r *AssetArtifactRepository) GetByID(ctx context.Context, id domain.AssetArtifactID) (domain.AssetArtifact, error) {
	var art domain.AssetArtifact

	row := r.db.QueryRow(ctx,
		`SELECT id, asset_id, type, mime_type, size, checksum, storage_key, created_at
		FROM asset_artifacts
		WHERE id = $1`,
		id,
	)

	if err := row.Scan(&art.ID, &art.AssetID, &art.Type, &art.MimeType,
		&art.Size, &art.Checksum, &art.StorageKey, &art.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.AssetArtifact{}, fmt.Errorf("postgres asset artifact: get by id: %w", repo.ErrNotFound)
		}
		return domain.AssetArtifact{}, fmt.Errorf("postgres asset artifact: get by id: %w", err)
	}

	return art, nil
}

func (r *AssetArtifactRepository) ListByAssetID(ctx context.Context, assetID domain.AssetID) ([]domain.AssetArtifact, error) {
	out := make([]domain.AssetArtifact, 0)

	rows, err := r.db.Query(ctx,
		`SELECT id, asset_id, type, mime_type, size, checksum, storage_key, created_at
		FROM asset_artifacts
		WHERE asset_id = $1`,
		assetID,
	)

	if err != nil {
		return nil, fmt.Errorf("postgres asset artifact: list by asset id: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var art domain.AssetArtifact

		if err := rows.Scan(&art.ID, &art.AssetID, &art.Type, &art.MimeType,
			&art.Size, &art.Checksum, &art.StorageKey, &art.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres asset artifact: list by asset id: %w", err)
		}

		out = append(out, art)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres asset artifact: list by asset id: rows: %w", err)
	}

	return out, nil
}

func (r *AssetArtifactRepository) Delete(ctx context.Context, id domain.AssetArtifactID) error {
	cmd, err := r.db.Exec(ctx,
		`DELETE FROM asset_artifacts
		WHERE id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("postgres asset artifact: delete: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("postgres asset artifact: delete: %w", repo.ErrNotFound)
	}

	return nil
}
