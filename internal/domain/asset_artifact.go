package domain

import "time"

type AssetArtifact struct {
	ID         AssetArtifactID
	AssetID    AssetID
	Type       ArtifactType
	MimeType   string
	Size       int64
	Checksum   string
	StorageKey string
	CreatedAt  time.Time
}
