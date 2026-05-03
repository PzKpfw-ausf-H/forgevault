package repo

import (
	"context"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type UsersRepo interface {
	Create(ctx context.Context, u domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id domain.UserID) (domain.User, error)
}
