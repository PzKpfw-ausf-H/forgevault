package domain

import (
	"fmt"
	"strings"
	"time"
)

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

func (a Asset) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("%w: asset id is empty", ErrValidation)
	}

	if strings.TrimSpace(a.Title) == "" {
		return fmt.Errorf("%w: asset title is empy", ErrValidation)
	}

	if a.UploadedBy == "" {
		return fmt.Errorf("%w: asset uploaded by is empty", ErrValidation)
	}

	if !a.Type.Valid() {
		fmt.Errorf("%w: asset type is not valid: %v", ErrValidation, a.Type)
	}

	if !a.ProcessingStatus.Valid() {
		fmt.Errorf("%w: asset processing status is not valid: %v", ErrValidation, a.ProcessingStatus)
	}

	if a.CreatedAt.IsZero() {
		fmt.Errorf("%w: asset created at is zero")
	}

	if a.UpdatedAt.IsZero() {
		fmt.Errorf("%w: asset updated at is zero")
	}

	return nil
}
