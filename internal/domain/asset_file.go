package domain

import (
	"fmt"
	"strings"
	"time"
)

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

func (a AssetFile) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("%w: asset file id is empty", ErrValidation)
	}

	if a.AssetID == "" {
		return fmt.Errorf("%w: asset file asset id is empty", ErrValidation)
	}

	if !a.Role.Valid() {
		return fmt.Errorf("%w: asset file role is not valid: %v", ErrValidation, a.Role)
	}

	switch a.Role {
	case FileRoleTexture:
		if a.TextureType == nil {
			return fmt.Errorf("%w: asset file texture type is nil", ErrValidation)
		}

		if !a.TextureType.Valid() {
			return fmt.Errorf("%w: asset file texture type is not valid: %v", ErrValidation, *a.TextureType)
		}
	default:
		if a.TextureType != nil {
			return fmt.Errorf("%w: asset file texture type must be nil for this file role: %v", ErrValidation, a.Role)
		}
	}

	if strings.TrimSpace(a.OriginalName) == "" {
		return fmt.Errorf("%w: asset file original name is empty", ErrValidation)
	}

	if strings.TrimSpace(a.MimeType) == "" {
		return fmt.Errorf("%w: asset file mime type is empty", ErrValidation)
	}

	if strings.TrimSpace(a.Extension) == "" {
		return fmt.Errorf("%w: asset file extension is empty", ErrValidation)
	}

	if a.Size < 0 {
		return fmt.Errorf("%w: asset file size can not be less then zero", ErrValidation)
	}

	if strings.TrimSpace(a.StorageKey) == "" {
		return fmt.Errorf("%w: asset file storage key is empty", ErrValidation)
	}

	if strings.TrimSpace(a.Checksum) == "" {
		return fmt.Errorf("%w: asset file checksum is empty", ErrValidation)
	}

	if a.CreatedAt.IsZero() {
		return fmt.Errorf("%w: asset file created at is zero", ErrValidation)
	}

	return nil
}
