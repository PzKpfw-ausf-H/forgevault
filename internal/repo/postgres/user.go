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

type UserRepository struct {
	db *pgxpool.Pool
}

var _ repo.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO users (id, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5)`,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("postgres user: create user: %w", repo.ErrAlreadyExists)
			}
		}

		return fmt.Errorf("postgres user: create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	var user domain.User

	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, role, created_at
		FROM users WHERE id = $1;`,
		id,
	)

	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("postgres user: get by id: %w", repo.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("postgres user: get by id: unexpected scan error: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1;`,
		email,
	)

	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("postgres user: get by email: %w", repo.ErrNotFound)
		}

		return domain.User{}, fmt.Errorf("postgres user: get by email: unexpected scan error: %w", err)
	}

	return user, nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID domain.UserID, newRole domain.Role) error {
	cmd, err := r.db.Exec(ctx,
		`UPDATE users
		SET role = $1
		WHERE id = $2;`,
		newRole,
		userID,
	)

	if err != nil {
		return fmt.Errorf("postgres user: update role: unexpected error: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("postgres user: update role: %w", repo.ErrNotFound)
	}

	return nil
}
