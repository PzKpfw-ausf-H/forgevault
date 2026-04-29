package http

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
