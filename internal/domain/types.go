package domain

type AssetID string
type UserID string
type FileID string
type AssetType string

const (
	AssetType3D    AssetType = "3d"
	AssetType2D    AssetType = "2d"
	AssetTypeAudio AssetType = "audio"
	AssetTypeVFX   AssetType = "vfx"
	AssetTypeDoc   AssetType = "doc"
)

func IsValidAssetType(t AssetType) bool {
	switch t {
	case AssetType2D, AssetType3D, AssetTypeAudio, AssetTypeVFX, AssetTypeDoc:
		return true
	default:
		return false
	}
}
