package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersSQLRepo struct {
	pool *pgxpool.Pool
}

func NewUsersSQLRepo(pool *pgxpool.Pool) *UsersSQLRepo {
	return &UsersSQLRepo{
		pool: pool,
	}
}

func (ur *UsersSQLRepo) Create(ctx context.Context, u domain.User) error {
	if err := domain.ValidateNewUser(&u); err != nil {
		return repo.ErrUnauthorized
	}

	_, err := ur.pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, created_at) VALUES ($1, $2, $3, $4)`,
		u.ID,
		strings.ToLower(u.Email),
		u.PasswordHash,
		u.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("execute create user: %v", err)
	}

	return nil
}

func (ur *UsersSQLRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	row := ur.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1`,
		strings.ToLower(email),
	)

	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, repo.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("get by email: %v", err)
	}

	return user, nil
}

func (ur *UsersSQLRepo) GetByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	var user domain.User

	row := ur.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at
		FROM users
		WHERE id = $1`,
		id,
	)
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, repo.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("get user by id: %v", err)
	}

	return user, nil
}
