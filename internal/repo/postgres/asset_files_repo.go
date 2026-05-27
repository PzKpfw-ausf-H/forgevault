package postgres

import (
	"context"
	"errors"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetFilesSQLRepo struct {
	pool *pgxpool.Pool
}

func (afr *AssetFilesSQLRepo) GetMaxVersion(ctx context.Context, assetID domain.AssetID) (int, error) {
	var v int
	row := afr.pool.QueryRow(ctx,
		`SELECT version FROM asset_files
		WHERE asset_id == $1`,
		assetID,
	)

	if err := row.Scan(v); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, repo.ErrNotFound
		}

		return 0, err
	}

	return v, nil
}

func (afr *AssetFilesSQLRepo) Create(ctx context.Context, file domain.AssetFile) error {
	if _, err := afr.pool.Exec(ctx,
		`INSERT INTO asset_files (id, asset_id, version, filename, size_bytes, content_type, storage_key, checksum, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		file.ID,
		file.AssetID,
		file.Version,
		file.Filename,
		file.SizeBytes,
		file.ContentType,
		file.StorageKey,
		file.Checksum,
		file.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (afr *AssetFilesSQLRepo) GetByAssetVersion(ctx context.Context, assetID domain.AssetID, version int) (domain.AssetFile, error) {
	var file domain.AssetFile

	row := afr.pool.QueryRow(ctx,
		`SELECT id, asset_id, version, filename, size_bytes, content_type, storage_key, checksum, created_at
		FROM asset_files
		WHERE asset_id = $1 AND version = $2`,
		assetID,
		version,
	)

	if err := row.Scan(&file.ID, &file.AssetID, &file.Version, &file.Filename, &file.SizeBytes, &file.ContentType, &file.StorageKey, &file.Checksum, &file.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.AssetFile{}, repo.ErrNotFound
		}

		return domain.AssetFile{}, err
	}

	return file, nil
}

func (afr *AssetFilesSQLRepo) ListByAsset(ctx context.Context, assetID domain.AssetID) ([]domain.AssetFile, error) {
	rows, err := afr.pool.Query(ctx,
		`SELECT id, asset_id, version, filename, size_bytes, content_type, storage_key, checksum, created_at
		FROM asset_files
		WHERE asset_id = $1`,
		assetID)
	if err != nil {
		return []domain.AssetFile{}, err
	}

	out := make([]domain.AssetFile, 0)

	for rows.Next() {
		var f domain.AssetFile

		if err := rows.Scan(&f.ID, &f.AssetID, &f.Version, &f.Filename, &f.SizeBytes, &f.ContentType, &f.StorageKey, &f.Checksum, &f.CreatedAt); err != nil {
			return []domain.AssetFile{}, err
		}

		out = append(out, f)
	}

	if err := rows.Err(); err != nil {
		return []domain.AssetFile{}, err
	}

	if len(out) == 0 {
		return []domain.AssetFile{}, repo.ErrNotFound
	}

	return out, nil
}
