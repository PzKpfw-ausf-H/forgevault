package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
)

type MemRepo struct {
	assets   []domain.Asset
	assetMap map[domain.AssetID]domain.Asset
	mu       sync.RWMutex
}

func NewMemRepo() *MemRepo {
	return &MemRepo{
		assets:   []domain.Asset{},
		assetMap: make(map[domain.AssetID]domain.Asset),
	}
}

func (r *MemRepo) Create(ctx context.Context, a domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := domain.ValidateNewAsset(&a); err != nil {
		return err
	}
	if _, exists := r.assetMap[a.ID]; exists {
		return repo.ErrConflict
	}
	r.assetMap[a.ID] = a
	r.assets = append(r.assets, a)
	return nil
}

func (r *MemRepo) GetByID(ctx context.Context, id domain.AssetID) (domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	a, exists := r.assetMap[id]
	if !exists {
		return domain.Asset{}, repo.ErrNotFound
	}
	assetToReturn := a
	assetToReturn.Tags = make([]string, len(a.Tags))
	copy(assetToReturn.Tags, a.Tags)
	return assetToReturn, nil
}

func (r *MemRepo) List(ctx context.Context, f repo.AssetFilter) ([]domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	limit := f.Limit
	offset := f.Offset

	if f.Limit <= 0 {
		limit = 20
	}
	if f.Offset < 0 {
		offset = 0
	}
	skipped := 0

	assets := make([]domain.Asset, 0, len(r.assets))
	for _, a := range r.assets {
		if f.Type != nil && a.Type != *f.Type {
			continue
		}
		if f.Tag != nil {
			hasTag := false
			for _, tag := range a.Tags {
				if tag == *f.Tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		if f.TitleSub != nil && !strings.Contains(a.Title, *f.TitleSub) {
			continue
		}
		if f.AuthorID != nil && a.AuthorID != *f.AuthorID {
			continue
		}
		if skipped < offset {
			skipped++
			continue
		}
		assets = append(assets, a)

		if len(assets) >= limit {
			break
		}
	}
	return assets, nil
}

func (r *MemRepo) Update(ctx context.Context, a domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := domain.ValidateNewAsset(&a); err != nil {
		return err
	}

	existing, exists := r.assetMap[a.ID]
	if !exists {
		return repo.ErrNotFound
	}
	existing.Title = a.Title
	existing.Description = a.Description
	existing.UpdatedAt = a.UpdatedAt
	existing.Type = a.Type
	existing.Tags = make([]string, len(a.Tags))
	copy(existing.Tags, a.Tags)
	r.assetMap[a.ID] = existing

	for i, asset := range r.assets {
		if asset.ID == a.ID {
			r.assets[i] = existing
			break
		}
	}

	return nil
}

func (r *MemRepo) Delete(ctx context.Context, id domain.AssetID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	assetToDelete, exists := r.assetMap[id]
	if !exists {
		return repo.ErrNotFound
	}
	lastIdx := len(r.assets) - 1
	lastAsset := r.assets[lastIdx]
	idToDelete := getAssetIdToDelete(r.assets, id)
	if lastAsset.ID != assetToDelete.ID {
		r.assets[idToDelete] = r.assets[lastIdx]
	}
	r.assets = r.assets[:lastIdx]
	delete(r.assetMap, id)

	return nil
}

func getAssetIdToDelete(assets []domain.Asset, id domain.AssetID) int {
	idToDelete := 0
	for _, asset := range assets {
		if asset.ID == id {
			break
		}
		idToDelete++
	}

	return idToDelete
}
