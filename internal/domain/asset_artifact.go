package domain

import (
	"fmt"
	"strings"
	"time"
)

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

func (a AssetArtifact) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("%w: asset artifact id is empty", ErrValidation)
	}

	if a.AssetID == "" {
		return fmt.Errorf("%w: asset artifact asset id is empty", ErrValidation)
	}

	if !a.Type.Valid() {
		return fmt.Errorf("%w: asset artifact type is not valid: %v", ErrValidation, a.Type)
	}

	if strings.TrimSpace(a.MimeType) == "" {
		return fmt.Errorf("%w: asset artifact mime type is empty")
	}

	if a.Size < 0 {
		return fmt.Errorf("%w: asset artifact size cannot be less then zero")
	}

	if strings.TrimSpace(a.Checksum) == "" {
		return fmt.Errorf("%w: asset artifact checksum is empty")
	}

	if strings.TrimSpace(a.StorageKey) == "" {
		return fmt.Errorf("%w: asset artifact storage key is empty")
	}

	if a.CreatedAt.IsZero() {
		return fmt.Errorf("%w: asset artifact created at is zero")
	}

	return nil
}
