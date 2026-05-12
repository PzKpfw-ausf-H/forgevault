package repo

import (
	"context"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/google/uuid"
)

type RefreshSessionsRepo interface {
	Create(ctx context.Context, s domain.RefreshSession) error
	GetByHash(ctx context.Context, hash string) (domain.RefreshSession, error)
	Revoke(ctx context.Context, id uuid.UUID, revokedAt time.Time) error
	Rotate(ctx context.Context, oldID uuid.UUID, revokedAt time.Time, next domain.RefreshSession) error
}
