package repo

import (
	"context"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type AssetListParams struct {
	Limit  int
	Offset int

	Type *domain.AssetType
}

type AssetListResult struct {
	Assets []domain.Asset
	Total  int64
}

type AssetRepository interface {
	Create(ctx context.Context, asset domain.Asset) error
	GetByID(ctx context.Context, assetID domain.AssetID) (domain.Asset, error)
	List(ctx context.Context, params AssetListParams) (AssetListResult, error)
	Update(ctx context.Context, asset domain.Asset) error
	Delete(ctx context.Context, id domain.AssetID) error
}
