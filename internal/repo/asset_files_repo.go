package repo

import (
	"context"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type AssetFilesRepo interface {
	GetMaxVersion(ctx context.Context, assetID domain.AssetID) (int, error)
	Create(ctx context.Context, file domain.AssetFile) error
	GetByAssetVersion(ctx context.Context, assetID domain.AssetID, version int) (domain.AssetFile, error)
	ListByAsset(ctx context.Context, assetID domain.AssetID) ([]domain.AssetFile, error)
}
