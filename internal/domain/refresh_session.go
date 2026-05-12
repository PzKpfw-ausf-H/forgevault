package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	ID         uuid.UUID
	UserID     UserID
	TokenHash  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	RevokedAt  time.Time
	ReplacedBy *uuid.UUID
}
