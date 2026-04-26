package domain

import "time"

type Asset struct {
	ID          AssetID
	Title       string
	Description string
	Type        AssetType
	Tags        []string
	AuthorID    UserID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
