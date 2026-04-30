package service

import (
	"context"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/google/uuid"
)

type AssetService struct {
	repo repo.AssetsRepo
	now  func() time.Time
}

func NewAssetService(repo repo.AssetsRepo) *AssetService {
	return &AssetService{
		repo: repo,
		now:  time.Now().UTC,
	}
}

func (s *AssetService) Create(ctx context.Context, input CreateAssetRequest) (domain.Asset, error) {
	asset := domain.Asset{
		ID:          domain.AssetID(uuid.New().String()),
		Title:       input.Title,
		Description: input.Description,
		Type:        input.Type,
		Tags:        make([]string, len(input.Tags)),
		AuthorID:    "demo-user",
		CreatedAt:   s.now().UTC(),
		UpdatedAt:   s.now().UTC(),
	}
	copy(asset.Tags, input.Tags)
	if err := s.repo.Create(ctx, asset); err != nil {
		return domain.Asset{}, err
	}

	return asset, nil
}

func (s *AssetService) GetByID(ctx context.Context, id domain.AssetID) (domain.Asset, error) {
	asset, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Asset{}, err
	}

	return asset, nil
}

func (s *AssetService) List(ctx context.Context, filter repo.AssetFilter) ([]domain.Asset, error) {
	assets, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (s *AssetService) Patch(ctx context.Context, id domain.AssetID, patch PatchAssetRequest) (domain.Asset, error) {
	assetToUpdate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Asset{}, err
	}
	if patch.Title != nil {
		assetToUpdate.Title = *patch.Title
	}
	if patch.Description != nil {
		assetToUpdate.Description = *patch.Description
	}
	if patch.Type != nil {
		assetToUpdate.Type = *patch.Type
	}
	if patch.Tags != nil {
		assetToUpdate.Tags = make([]string, len(*patch.Tags))
		copy(assetToUpdate.Tags, *patch.Tags)
	}
	assetToUpdate.UpdatedAt = s.now()

	if err := s.repo.Update(ctx, assetToUpdate); err != nil {
		return domain.Asset{}, err
	}

	return assetToUpdate, nil
}

func (s *AssetService) Delete(ctx context.Context, id domain.AssetID) error {
	return s.repo.Delete(ctx, id)
}
