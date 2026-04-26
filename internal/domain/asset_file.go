package domain

import "time"

type AssetFile struct {
	ID          FileID
	AssetID     AssetID
	Version     int
	Filename    string
	SizeBytes   int64
	ContentType string
	StorageKey  string
	Checksum    string
	CreatedAt   time.Time
}
