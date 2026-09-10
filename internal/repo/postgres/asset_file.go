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

type AssetFileRepository struct {
	db *pgxpool.Pool
}

var _ repo.AssetFileRepository = (*AssetFileRepository)(nil)

func NewAssetFileRepository(db *pgxpool.Pool) *AssetFileRepository {
	return &AssetFileRepository{db: db}
}

func (r *AssetFileRepository) Create(ctx context.Context, af domain.AssetFile) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO asset_files (id, asset_id, role, texture_type, original_name, mime_type, extension,
		size, storage_key, checksum, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		af.ID, af.AssetID, af.Role, af.TextureType, af.OriginalName, af.MimeType, af.Extension,
		af.Size, af.StorageKey, af.Checksum, af.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("postgres asset file: create: %w", repo.ErrAlreadyExists)
			}
		}

		return fmt.Errorf("postgres asset file: create: %w", err)
	}

	return nil
}

func (r *AssetFileRepository) GetByID(ctx context.Context, id domain.AssetFileID) (domain.AssetFile, error) {
	var af domain.AssetFile

	row := r.db.QueryRow(ctx,
		`SELECT id, asset_id, role, texture_type, original_name, mime_type, extension,
		size, storage_key, checksum, created_at
		FROM asset_files
		WHERE id = $1`,
		id,
	)

	if err := row.Scan(&af.ID, &af.AssetID, &af.Role, &af.TextureType, &af.OriginalName, &af.MimeType, &af.Extension,
		&af.Size, &af.StorageKey, &af.Checksum, &af.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.AssetFile{}, fmt.Errorf("postgres asset file: get by id: %w", repo.ErrNotFound)
		}

		return domain.AssetFile{}, fmt.Errorf("postgres asset file: get by id: %w", err)
	}

	return af, nil
}

func (r *AssetFileRepository) ListByAssetID(ctx context.Context, assetID domain.AssetID) ([]domain.AssetFile, error) {
	out := make([]domain.AssetFile, 0)

	rows, err := r.db.Query(ctx,
		`SELECT id, asset_id, role, texture_type, original_name, mime_type, extension,
		size, storage_key, checksum, created_at
		FROM asset_files
		WHERE asset_id = $1
		ORDER BY created_at, id`,
		assetID,
	)

	if err != nil {
		return nil, fmt.Errorf("postgres asset file: list by asset id: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var af domain.AssetFile

		if err := rows.Scan(&af.ID, &af.AssetID, &af.Role, &af.TextureType, &af.OriginalName, &af.MimeType, &af.Extension,
			&af.Size, &af.StorageKey, &af.Checksum, &af.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres asset file: list by asset id: %w", err)
		}

		out = append(out, af)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres asset file: list by asset id: %w", err)
	}

	return out, nil
}

func (r *AssetFileRepository) Delete(ctx context.Context, id domain.AssetFileID) error {
	cmd, err := r.db.Exec(ctx,
		`DELETE FROM asset_files
	WHERE id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("postgres asset file: delete: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("postgres asset file: delete: %w", repo.ErrNotFound)
	}

	return nil
}
