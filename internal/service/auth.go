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
	repo        repo.UsersRepo
	tm          *auth.TokenManager
	refreshRepo repo.RefreshSessionsRepo
	refreshTTL  time.Duration
}

func NewUserService(repo repo.UsersRepo, tm *auth.TokenManager, refreshRepo repo.RefreshSessionsRepo, refreshTTL time.Duration) *UserService {
	return &UserService{
		repo:        repo,
		tm:          tm,
		refreshRepo: refreshRepo,
		refreshTTL:  refreshTTL,
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

func (us *UserService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, expiresInSec, refreshExpiresIn int64, err error) {
	u, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", 0, 0, repo.ErrUnauthorized
	}
	if err2 := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err2 != nil {
		return "", "", 0, 0, repo.ErrUnauthorized
	}
	refreshToken, refErr := GenerateRefreshToken()
	if refErr != nil {
		return "", "", 0, 0, repo.ErrUnauthorized
	}
	hash := HashRefreshToken(refreshToken)
	accessToken, expiresInSec, err = us.tm.New(u.ID)

	now := time.Now().UTC()
	rs := domain.RefreshSession{
		ID:         uuid.New(),
		UserID:     u.ID,
		TokenHash:  hash,
		ExpiresAt:  now.Add(us.refreshTTL),
		CreatedAt:  now,
		RevokedAt:  time.Time{},
		ReplacedBy: nil,
	}

	if err := us.refreshRepo.Create(ctx, rs); err != nil {
		return "", "", 0, 0, err
	}

	refreshExpiresIn = int64(us.refreshTTL.Seconds())
	return accessToken, refreshToken, expiresInSec, refreshExpiresIn, nil
}

func (us *UserService) Refresh(ctx context.Context, refreshToken string) (newAccess, newRefresh string, expiresIn, refresExpiresIn int64, err error) {
	hash := HashRefreshToken(refreshToken)
	session, err := us.refreshRepo.GetByHash(ctx, hash)
	if err != nil {
		return "", "", 0, 0, err
	}
	now := time.Now().UTC()
	if now.After(session.ExpiresAt) {
		return "", "", 0, 0, repo.ErrUnauthorized
	}

	if !session.RevokedAt.IsZero() && session.RevokedAt.Before(now) {
		return "", "", 0, 0, repo.ErrUnauthorized
	}

	newRefreshToken, newErr := GenerateRefreshToken()
	if newErr != nil {
		return "", "", 0, 0, newErr
	}
	newHash := HashRefreshToken(newRefreshToken)
	nextSession := domain.RefreshSession{
		ID:         uuid.New(),
		UserID:     session.UserID,
		TokenHash:  newHash,
		ExpiresAt:  now.Add(us.refreshTTL),
		CreatedAt:  now,
		RevokedAt:  time.Time{},
		ReplacedBy: nil,
	}

	newAccess, expiresIn, err = us.tm.New(session.UserID)

	if err := us.refreshRepo.Rotate(ctx, session.ID, now, nextSession); err != nil {
		return "", "", 0, 0, err
	}

	refresExpiresIn = int64(us.refreshTTL.Seconds())
	return newAccess, newRefreshToken, expiresIn, refresExpiresIn, err
}

func (us *UserService) Logout(ctx context.Context, refreshToken string) error {
	hash := HashRefreshToken(refreshToken)
	session, err := us.refreshRepo.GetByHash(ctx, hash)
	if err != nil {
		return err
	}

	if err := us.refreshRepo.Revoke(ctx, session.ID, session.RevokedAt); err != nil {
		return err
	}

	return nil
}
