package domain

import (
	"errors"
	"testing"
	"time"
)

func validAssetArtifact() AssetArtifact {
	now := time.Now()

	return AssetArtifact{
		ID:         AssetArtifactID("assetArtifact-1"),
		AssetID:    AssetID("asset-1"),
		Type:       ArtifactTypePreviewImage,
		MimeType:   "image/png",
		Size:       15,
		Checksum:   "something_not_empty",
		StorageKey: "path/to/file",
		CreatedAt:  now,
	}
}

func TestAssetArtifact_Validate(t *testing.T) {
	tests := []struct {
		name     string
		artifact AssetArtifact
		wantErr  bool
	}{
		{
			name:     "valid artifact",
			artifact: validAssetArtifact(),
			wantErr:  false,
		},
		{
			name: "empty ID",
			artifact: AssetArtifact{
				ID:         "",
				AssetID:    AssetID("asset-1"),
				Type:       ArtifactTypePreviewImage,
				MimeType:   "image/png",
				Size:       15,
				Checksum:   "something_not_empty",
				StorageKey: "path/to/file",
				CreatedAt:  time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty asset ID",
			artifact: AssetArtifact{
				ID:         AssetArtifactID("assetArtifact-1"),
				AssetID:    "",
				Type:       ArtifactTypePreviewImage,
				MimeType:   "image/png",
				Size:       15,
				Checksum:   "something_not_empty",
				StorageKey: "path/to/file",
				CreatedAt:  time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			artifact: AssetArtifact{
				ID:         AssetArtifactID("assetArtifact-1"),
				AssetID:    AssetID("asset-1"),
				Type:       ArtifactType("Banana"),
				MimeType:   "image/png",
				Size:       15,
				Checksum:   "something_not_empty",
				StorageKey: "path/to/file",
				CreatedAt:  time.Now(),
			},
			wantErr: true,
		},
		{
			name: "blank mime type",
			artifact: AssetArtifact{
				ID:         AssetArtifactID("assetArtifact-1"),
				AssetID:    AssetID("asset-1"),
				Type:       ArtifactTypePreviewImage,
				MimeType:   "       ",
				Size:       15,
				Checksum:   "something_not_empty",
				StorageKey: "path/to/file",
				CreatedAt:  time.Now(),
			},
			wantErr: true,
		},
		{
			name: "negative size",
			artifact: AssetArtifact{
				ID:         AssetArtifactID("assetArtifact-1"),
				AssetID:    AssetID("asset-1"),
				Type:       ArtifactTypePreviewImage,
				MimeType:   "image/png",
				Size:       -15,
				Checksum:   "something_not_empty",
				StorageKey: "path/to/file",
				CreatedAt:  time.Now(),
			},
			wantErr: true,
		},
		{
			name: "blank checksum",
			artifact: AssetArtifact{
				ID:         AssetArtifactID("assetArtifact-1"),
				AssetID:    AssetID("asset-1"),
				Type:       ArtifactTypePreviewImage,
				MimeType:   "image/png",
				Size:       15,
				Checksum:   "      ",
				StorageKey: "path/to/file",
				CreatedAt:  time.Now(),
			},
			wantErr: true,
		},
		{
			name: "blank storage key",
			artifact: AssetArtifact{
				ID:         AssetArtifactID("assetArtifact-1"),
				AssetID:    AssetID("asset-1"),
				Type:       ArtifactTypePreviewImage,
				MimeType:   "image/png",
				Size:       15,
				Checksum:   "something_not_empty",
				StorageKey: "",
				CreatedAt:  time.Now(),
			},
			wantErr: true,
		},
		{
			name: "zero created at",
			artifact: AssetArtifact{
				ID:         AssetArtifactID("assetArtifact-1"),
				AssetID:    AssetID("asset-1"),
				Type:       ArtifactTypePreviewImage,
				MimeType:   "image/png",
				Size:       15,
				Checksum:   "something_not_empty",
				StorageKey: "path/to/file",
				CreatedAt:  time.Time{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			err := tt.artifact.Validate()

			if err != nil && !errors.Is(err, ErrValidation) {
				t.Fatalf("expected ErrValidation, got: %v", err)
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("%s: unexpected error: %v", tt.name, err)
			}

			if tt.wantErr && err == nil {
				t.Fatalf("%s: expected error", tt.name)
			}
		})
	}
}
