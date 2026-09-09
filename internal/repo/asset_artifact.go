package repo

import (
	"context"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type AssetArtifactRepository interface {
	Create(ctx context.Context, artifact domain.AssetArtifact) error
	GetByID(ctx context.Context, id domain.AssetArtifactID) (domain.AssetArtifact, error)
	ListByAssetID(ctx context.Context, assetID domain.AssetID) ([]domain.AssetArtifact, error)
	Delete(ctx context.Context, id domain.AssetArtifactID) error
}
