package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/google/uuid"
)

func validUser() domain.User {
	now := time.Now().UTC().Truncate(time.Microsecond)

	return domain.User{
		ID:           domain.UserID(uuid.NewString()),
		Email:        "student@gmail.com",
		Role:         domain.RoleStudent,
		PasswordHash: "hash",
		CreatedAt:    now,
	}
}

func TestUserRepository_CreateAndGetByID(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	r := NewUserRepository(db)
	ctx := context.Background()

	want := validUser()

	if err := r.Create(ctx, want); err != nil {
		t.Fatalf("create() error: %v", err)
	}

	got, err := r.GetByID(ctx, want.ID)
	if err != nil {
		t.Fatalf("GetByID() error: %v", err)
	}

	if got.ID != want.ID {
		t.Errorf("want id %q but got %q", want.ID, got.ID)
	}

	if got.Email != want.Email {
		t.Errorf("want email %s but got %s", want.Email, got.Email)
	}

	if got.PasswordHash != want.PasswordHash {
		t.Errorf("want hash %s but got %s", want.PasswordHash, got.PasswordHash)
	}

	if got.Role != want.Role {
		t.Errorf("want role %q but got %q", want.Role, got.Role)
	}

	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("want created_at %v but got %v", want.CreatedAt, got.CreatedAt)
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	now := time.Now().UTC().Truncate(time.Microsecond)
	ctx := context.Background()

	r := NewUserRepository(db)
	user1 := domain.User{
		ID:           domain.UserID(uuid.NewString()),
		Email:        "student@email.com",
		Role:         domain.RoleStudent,
		PasswordHash: "hash",
		CreatedAt:    now,
	}

	user2 := domain.User{
		ID:           domain.UserID(uuid.NewString()),
		Email:        "student@email.com",
		Role:         domain.RoleStudent,
		PasswordHash: "hash",
		CreatedAt:    now,
	}

	if err := r.Create(ctx, user1); err != nil {
		t.Fatalf("create() user1 error: %v", err)
	}

	err := r.Create(ctx, user2)
	if !errors.Is(err, repo.ErrAlreadyExists) {
		t.Fatalf("Create() error = %v, want %v", err, repo.ErrAlreadyExists)
	}
}

func TestUserRepository_GetByIDNonExisting(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	r := NewUserRepository(db)

	_, err := r.GetByID(ctx, domain.UserID("99999999-9999-9999-9999-999999999999"))

	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("GetById() nonexistst expected ErrNotFound but got %v", err)
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	want := validUser()
	r := NewUserRepository(db)

	if err := r.Create(ctx, want); err != nil {
		t.Fatalf("create() user error: %v", err)
	}

	got, err := r.GetByEmail(ctx, want.Email)
	if err != nil {
		t.Fatalf("GetByEmail() user error: %v", err)
	}

	if got.ID != want.ID {
		t.Errorf("want id %q but got %q", want.ID, got.ID)
	}

	if got.Email != want.Email {
		t.Errorf("want email %s but got %s", want.Email, got.Email)
	}

	if got.PasswordHash != want.PasswordHash {
		t.Errorf("want hash %s but got %s", want.PasswordHash, got.PasswordHash)
	}

	if got.Role != want.Role {
		t.Errorf("want role %q but got %q", want.Role, got.Role)
	}

	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("want created_at %v but got %v", want.CreatedAt, got.CreatedAt)
	}
}

func TestUserRepository_UpdateRole(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	r := NewUserRepository(db)
	user := validUser()

	if err := r.Create(ctx, user); err != nil {
		t.Fatalf("create() user error: %v", err)
	}

	if err := r.UpdateRole(ctx, user.ID, domain.RoleEmployee); err != nil {
		t.Fatalf("UpdateRole() error: %v", err)
	}

	got, err := r.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID() error: %v", err)
	}

	if got.Role != domain.RoleEmployee {
		t.Fatalf("expected role %q but got %q", domain.RoleEmployee, got.Role)
	}
}

func TestUserRepository_UpdateRole_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	r := NewUserRepository(db)

	err := r.UpdateRole(ctx, domain.UserID("99999999-9999-9999-9999-999999999999"), domain.RoleEmployee)

	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("UpdateRole() for nonexistent: expected: %v", repo.ErrNotFound)
	}
}
