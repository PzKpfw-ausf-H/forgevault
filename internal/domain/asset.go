package domain

import "time"

type AssetType string

const (
	AssetType3D    AssetType = "3d"
	AssetType2D    AssetType = "2d"
	AssetTypeAudio AssetType = "audio"
	AssetTypeVFX   AssetType = "vfx"
	AssetTypeDoc   AssetType = "doc"
)

func (t AssetType) Valid() bool {
	switch t {
	case AssetType2D, AssetType3D, AssetTypeAudio, AssetTypeVFX, AssetTypeDoc:
		return true
	default:
		return false
	}
}

type Asset struct {
	ID          AssetID
	Title       string
	Description string

	AuthorName string
	UploadedBy UserID

	Type             AssetType
	Tags             []string
	ProcessingStatus ProcessingStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}
