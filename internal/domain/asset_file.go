package domain

import "time"

// AssetFile represents an original user-uploaded file
// Generated derivatives belong to AssetArtifacts

type AssetFile struct {
	ID           AssetFileID
	AssetID      AssetID
	Role         FileRole
	TextureType  *TextureType
	OriginalName string
	MimeType     string
	Extension    string
	Size         int64
	StorageKey   string
	Checksum     string
	CreatedAt    time.Time
}
