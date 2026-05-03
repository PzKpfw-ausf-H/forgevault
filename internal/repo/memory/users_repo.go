package memory

import (
	"context"
	"sync"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
)

type UserMemRepo struct {
	byID    map[domain.UserID]domain.User
	byEmail map[string]domain.User
	mu      sync.RWMutex
}

func NewUserMemRepo() *UserMemRepo {
	return &UserMemRepo{
		byID:    make(map[domain.UserID]domain.User),
		byEmail: make(map[string]domain.User),
	}
}

func (ur *UserMemRepo) Create(ctx context.Context, u domain.User) error {
	ur.mu.Lock()
	defer ur.mu.Unlock()

	if err := domain.ValidateNewUser(&u); err != nil {
		return err
	}
	if _, exists := ur.byEmail[u.Email]; exists {
		return repo.ErrConflict
	}
	if _, exists := ur.byID[u.ID]; exists {
		return repo.ErrConflict
	}
	ur.byEmail[u.Email] = u
	ur.byID[u.ID] = u
	return nil
}

func (ur *UserMemRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	u, exists := ur.byEmail[email]
	if !exists {
		return domain.User{}, repo.ErrNotFound
	}
	return u, nil
}

func (ur *UserMemRepo) GetByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	ur.mu.RLock()
	defer ur.mu.RUnlock()

	u, exists := ur.byID[id]
	if !exists {
		return domain.User{}, repo.ErrNotFound
	}

	return u, nil
}
