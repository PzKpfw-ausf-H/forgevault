package domain

import (
	"errors"
	"testing"
	"time"
)

func ptr[T any](v T) *T {
	return &v
}

func validMain() AssetFile {
	now := time.Now()

	return AssetFile{
		ID:           AssetFileID("assetFile-1"),
		AssetID:      AssetID("asset-1"),
		Role:         FileRoleMain,
		TextureType:  nil,
		OriginalName: "something_file.blend",
		MimeType:     "application/json",
		Extension:    ".blend",
		Size:         15,
		StorageKey:   "path/to/file",
		Checksum:     "something_not_empty",
		CreatedAt:    now,
	}
}

func validTexture() AssetFile {
	now := time.Now()

	return AssetFile{
		ID:           AssetFileID("assetFile-1"),
		AssetID:      AssetID("asset-1"),
		Role:         FileRoleTexture,
		TextureType:  ptr(TextureTypeBaseColor),
		OriginalName: "base_color.png",
		MimeType:     "image/png",
		Extension:    ".png",
		Size:         15,
		StorageKey:   "path/to/file",
		Checksum:     "something_not_empty",
		CreatedAt:    now,
	}
}

func TestAssetFile_Validate(t *testing.T) {
	tests := []struct {
		name      string
		assetFile AssetFile
		wantErr   bool
	}{
		{
			name:      "valid main",
			assetFile: validMain(),
			wantErr:   false,
		},
		{
			name:      "valid texture",
			assetFile: validTexture(),
			wantErr:   false,
		},
		{
			name: "empty ID",
			assetFile: AssetFile{
				ID:           "",
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty asset id",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      "",
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRole("Banana"),
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "texture without texture type",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  nil,
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "texture with invalid texture type",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureType("Banana")),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "non-texture with texture type",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleAttachment,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "blank original name",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "         ",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "blank mime type",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "       ",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "blank extension",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    "",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "negative size",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         -15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "blank storage key",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "blank checksum",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "zero created at",
			assetFile: AssetFile{
				ID:           AssetFileID("assetFile-1"),
				AssetID:      AssetID("asset-1"),
				Role:         FileRoleTexture,
				TextureType:  ptr(TextureTypeBaseColor),
				OriginalName: "base_color.png",
				MimeType:     "image/png",
				Extension:    ".png",
				Size:         15,
				StorageKey:   "path/to/file",
				Checksum:     "something_not_empty",
				CreatedAt:    time.Time{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			err := tt.assetFile.Validate()

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
