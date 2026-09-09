package domain

// FileRole block
type FileRole string

const (
	FileRoleMain          FileRole = "main"
	FileRoleTexture       FileRole = "texture"
	FileRoleAttachment    FileRole = "attachment"
	FileRolePreviewSource FileRole = "preview_source"
)

func (r FileRole) Valid() bool {
	switch r {
	case FileRoleMain, FileRoleTexture, FileRoleAttachment, FileRolePreviewSource:
		return true
	default:
		return false
	}
}

// TextureType block
type TextureType string

const (
	TextureTypeBaseColor TextureType = "base_color"
	TextureTypeNormal    TextureType = "normal"
	TextureTypeRoughness TextureType = "roughness"
	TextureTypeMetallic  TextureType = "metallic"
	TextureTypeAO        TextureType = "ao"
	TextureTypeEmissive  TextureType = "emissive"
	TextureTypeOpacity   TextureType = "opacity"
	TextureTypeHeight    TextureType = "height"
	TextureTypeORM       TextureType = "orm"
	TextureTypeOther     TextureType = "other"
)

func (t TextureType) Valid() bool {
	switch t {
	case TextureTypeBaseColor, TextureTypeNormal, TextureTypeRoughness, TextureTypeMetallic,
		TextureTypeAO, TextureTypeEmissive, TextureTypeOpacity, TextureTypeHeight,
		TextureTypeORM, TextureTypeOther:
		return true
	default:
		return false
	}
}

// ArtifactType block
type ArtifactType string

const (
	ArtifactTypeViewerModel  ArtifactType = "viewer_model"
	ArtifactTypeThumbnail    ArtifactType = "thumbnail"
	ArtifactTypePreviewImage ArtifactType = "preview_image"
)

func (t ArtifactType) Valid() bool {
	switch t {
	case ArtifactTypeViewerModel, ArtifactTypeThumbnail, ArtifactTypePreviewImage:
		return true
	default:
		return false
	}
}

// ProcessingStatus block
type ProcessingStatus string

const (
	ProcessingStatusPending    ProcessingStatus = "pending"
	ProcessingStatusProcessing ProcessingStatus = "processing"
	ProcessingStatusReady      ProcessingStatus = "ready"
	ProcessingStatusFailed     ProcessingStatus = "failed"
)

func (s ProcessingStatus) Valid() bool {
	switch s {
	case ProcessingStatusPending, ProcessingStatusProcessing, ProcessingStatusReady, ProcessingStatusFailed:
		return true
	default:
		return false
	}
}
