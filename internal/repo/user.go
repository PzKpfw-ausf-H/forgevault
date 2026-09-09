package repo

import (
	"context"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id domain.UserID) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateRole(ctx context.Context, userID domain.UserID, newRole domain.Role) error
}
