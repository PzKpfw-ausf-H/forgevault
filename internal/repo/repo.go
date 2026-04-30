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
	Type     *domain.AssetType `json:"type"`
	Tag      *string           `json:"tag"`
	TitleSub *string           `json:"titleSub"`
	AuthorID *domain.UserID    `json:"authorId"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}
