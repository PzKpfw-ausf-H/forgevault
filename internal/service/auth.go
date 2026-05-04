package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/auth"
	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repo.UsersRepo
	tm   *auth.TokenManager
}

func NewUserService(repo repo.UsersRepo, tm *auth.TokenManager) *UserService {
	return &UserService{
		repo: repo,
		tm:   tm,
	}
}

func (us *UserService) Register(ctx context.Context, email, password string) (domain.User, error) {
	user := domain.User{
		ID:           domain.UserID(uuid.New().String()),
		Email:        email,
		PasswordHash: "",
		CreatedAt:    time.Now(),
	}
	if len(password) < 8 {
		return domain.User{}, repo.ErrInvalidPassword
	}
	passwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("err denerating hash: %v", err)
	}
	user.PasswordHash = string(passwdHash)
	if err := us.repo.Create(ctx, user); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (us *UserService) Login(ctx context.Context, email, password string) (token string, expiresInSec int64, err error) {
	u, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", 0, err
	}
	if err2 := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err2 != nil {
		return "", 0, repo.ErrUnauthorized
	}

	return us.tm.New(u.ID)
}
