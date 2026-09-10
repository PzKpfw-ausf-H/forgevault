package repo

import (
	"context"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type AssetFileRepository interface {
	Create(ctx context.Context, af domain.AssetFile) error
	GetByID(ctx context.Context, id domain.AssetFileID) (domain.AssetFile, error)
	ListByAssetID(ctx context.Context, assetID domain.AssetID) ([]domain.AssetFile, error)
	Delete(ctx context.Context, id domain.AssetFileID) error
}
