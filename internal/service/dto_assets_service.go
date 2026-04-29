package service

import (
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type CreateAssetRequest struct {
	Title       string
	Description string
	Type        domain.AssetType
	Tags        []string
}

type PatchAssetRequest struct {
	Title       *string
	Description *string
	Type        *domain.Asset
	Tags        *[]string
}

type AssetResponse struct {
	ID          domain.AssetID
	Title       string
	Description string
	Type        domain.AssetType
	Tags        []string
	AuthorID    domain.UserID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
