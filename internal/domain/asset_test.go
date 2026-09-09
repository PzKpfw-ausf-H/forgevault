package domain

import (
	"errors"
	"testing"
	"time"
)

func validAsset() Asset {
	now := time.Now()

	return Asset{
		ID:               AssetID("asset-1"),
		Title:            "Tank",
		UploadedBy:       UserID("user-1"),
		Type:             AssetType3D,
		ProcessingStatus: ProcessingStatusReady,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func TestAsset_Validate(t *testing.T) {
	tests := []struct {
		name    string
		asset   Asset
		wantErr bool
	}{
		{
			name:    "validAsset",
			asset:   validAsset(),
			wantErr: false,
		},
		{
			name: "empty ID",
			asset: Asset{
				ID:               "",
				Title:            "Tank",
				UploadedBy:       UserID("user-1"),
				Type:             AssetType3D,
				ProcessingStatus: ProcessingStatusReady,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "blank title",
			asset: Asset{
				ID:               AssetID("asset-1"),
				Title:            "         ",
				UploadedBy:       UserID("user-1"),
				Type:             AssetType3D,
				ProcessingStatus: ProcessingStatusReady,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "empty uploaded by",
			asset: Asset{
				ID:               AssetID("asset-1"),
				Title:            "Tank",
				UploadedBy:       "",
				Type:             AssetType3D,
				ProcessingStatus: ProcessingStatusReady,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid asset type",
			asset: Asset{
				ID:               AssetID("asset-1"),
				Title:            "Tank",
				UploadedBy:       UserID("user-1"),
				Type:             AssetType("Banana"),
				ProcessingStatus: ProcessingStatusReady,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid processing status",
			asset: Asset{
				ID:               AssetID("asset-1"),
				Title:            "Tank",
				UploadedBy:       UserID("user-1"),
				Type:             AssetType3D,
				ProcessingStatus: ProcessingStatus("Banana"),
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "zero created at",
			asset: Asset{
				ID:               AssetID("asset-1"),
				Title:            "Tank",
				UploadedBy:       UserID("user-1"),
				Type:             AssetType3D,
				ProcessingStatus: ProcessingStatusReady,
				CreatedAt:        time.Time{},
				UpdatedAt:        time.Now(),
			},
			wantErr: true,
		},
		{
			name: "zero updated at",
			asset: Asset{
				ID:               AssetID("asset-1"),
				Title:            "Tank",
				UploadedBy:       UserID("user-1"),
				Type:             AssetType3D,
				ProcessingStatus: ProcessingStatusReady,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Time{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			err := tt.asset.Validate()

			if err != nil && !errors.Is(err, ErrValidation) {
				t.Fatalf("expected ErrValidation, got: %v", err)
			}

			if tt.wantErr && err == nil {
				t.Fatalf("%s: expected error", tt.name)
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("%s: unexpected error: %v", tt.name, err)
			}
		})
	}
}
