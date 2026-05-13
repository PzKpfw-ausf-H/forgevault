package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshSessionSQLRepo struct {
	pool *pgxpool.Pool
}

func NewRefreshSessionSQLRepo(pool *pgxpool.Pool) *RefreshSessionSQLRepo {
	return &RefreshSessionSQLRepo{
		pool: pool,
	}
}

func (rr *RefreshSessionSQLRepo) Create(ctx context.Context, s domain.RefreshSession) error {
	_, err := rr.pool.Exec(ctx,
		`INSERT INTO refresh_sessions (id, user_id, token_hash, expires_at, created_at, revoked_at, replaced_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		s.ID,
		s.UserID,
		s.TokenHash,
		s.ExpiresAt,
		s.CreatedAt,
		s.RevokedAt,
		s.ReplacedBy,
	)

	if err != nil {
		return fmt.Errorf("execute create refresh sessions: %v", err)
	}

	return nil
}

func (rr *RefreshSessionSQLRepo) GetByHash(ctx context.Context, hash string) (domain.RefreshSession, error) {
	var rs domain.RefreshSession

	row := rr.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at, revoked_at, replaced_by
		FROM refresh_sessions
		WHERE token_hash = $1`,
		hash,
	)

	if err := row.Scan(&rs.ID, &rs.UserID, &rs.TokenHash, &rs.ExpiresAt, &rs.CreatedAt, &rs.RevokedAt, &rs.ReplacedBy); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RefreshSession{}, repo.ErrNotFound
		}
		return domain.RefreshSession{}, fmt.Errorf("refresh session get by hash: %v", err)
	}

	return rs, nil
}

func (rr *RefreshSessionSQLRepo) Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	cmd, err := rr.pool.Exec(ctx,
		`UPDATE refresh_sessions
		SET revoked_at = $1
		WHERE id = $2 AND revoked_at IS NULL`,
		revokedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("update refresh session revoke: %v", err)
	}

	if cmd.RowsAffected() == 0 {
		return repo.ErrUnauthorized
	}

	return nil
}

func (rr *RefreshSessionSQLRepo) Rotate(ctx context.Context, oldID uuid.UUID, revokedAt time.Time, next domain.RefreshSession) error {
	tx, err := rr.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("refresh sessions begin tx: %v", err)
	}
	defer tx.Rollback(ctx)

	cmd, err := tx.Exec(ctx,
		`UPDATE refresh_sessions
		SET revoked_at = $1, replaced_by = $2
		WHERE id = $3 AND revoked_at IS NULL`,
		revokedAt,
		next.ID,
		oldID,
	)

	if err != nil {
		return fmt.Errorf("update refresh sessions tx: %v", err)
	}

	if cmd.RowsAffected() == 0 {
		return repo.ErrUnauthorized
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO refresh_sessions (id, user_id, token_hash, expires_at, created_at, revoked_at, replaced_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		next.ID,
		next.UserID,
		next.TokenHash,
		next.ExpiresAt,
		next.CreatedAt,
		next.RevokedAt,
		next.ReplacedBy,
	)

	if err != nil {
		return fmt.Errorf("refresh sessions insert tx: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("refresh sessions commit: %v", err)
	}

	return nil
}
