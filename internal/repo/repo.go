package repo

import (
	"context"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type AssetsRepo interface {
	Create(ctx context.Context, a domain.Asset) error
	GetByID(ctx context.Context, id domain.AssetID) (domain.Asset, error)
	List(ctx context.Context, f AssetFilter) ([]domain.Asset, error)
	Update(ctx context.Context, a domain.Asset) error
	Delete(ctx context.Context, id domain.AssetID) error
}

type AssetFilter struct {
	Type     *domain.AssetType
	Tag      *string
	TitleSub *string
	AuthorID *domain.UserID
	Limit    int
	Offset   int
}
