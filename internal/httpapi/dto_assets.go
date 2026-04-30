package httpapi

import (
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
)

type CreateAssetRequest struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
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

func toAssetResponse(a domain.Asset) AssetResponse {
	var ar AssetResponse
	ar.Tags = make([]string, len(a.Tags))
	copy(ar.Tags, a.Tags)
	ar.ID = a.ID
	ar.Title = a.Title
	ar.Description = a.Description
	ar.Type = a.Type
	ar.AuthorID = a.AuthorID
	ar.CreatedAt = a.CreatedAt
	ar.UpdatedAt = a.UpdatedAt

	return ar
}
