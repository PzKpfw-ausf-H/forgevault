package service

import (
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type CreateAssetRequest struct {
	Title       string           `json:"title"`
	Description string           `json:"desciption"`
	Type        domain.AssetType `json:"type"`
	Tags        []string         `json:"tags"`
}

type PatchAssetRequest struct {
	Title       *string           `json:"title"`
	Description *string           `json:"description"`
	Type        *domain.AssetType `json:"type"`
	Tags        *[]string         `json:"tags"`
}

type AssetResponse struct {
	ID          domain.AssetID   `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Type        domain.AssetType `json:"type"`
	Tags        []string         `json:"tags"`
	AuthorID    domain.UserID    `json:"authorId"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}
