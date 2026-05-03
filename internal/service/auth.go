package service

import (
	"context"
	"fmt"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      repo.UsersRepo
	jwtSecret []byte
	ttl       time.Duration
	issuer    string
}

func NewUserService(repo repo.UsersRepo, jwtSecret string, ttl time.Duration, issuer string) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		ttl:       ttl,
		issuer:    issuer,
	}
}

type AuthClaims struct {
	jwt.RegisteredClaims
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

	now := time.Now().UTC()
	exp := now.Add(us.ttl)

	claims := AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   string(u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    us.issuer,
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(us.jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return signed, int64(us.ttl.Seconds()), nil
}
