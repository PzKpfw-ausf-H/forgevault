package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func validAssetArtifact() domain.AssetArtifact {
	now := time.Now().UTC().Truncate(time.Microsecond)

	return domain.AssetArtifact{
		ID:         domain.AssetArtifactID(uuid.NewString()),
		AssetID:    domain.AssetID(uuid.NewString()),
		Type:       domain.ArtifactTypePreviewImage,
		MimeType:   "image/png",
		Size:       15,
		Checksum:   "something_not_empty",
		StorageKey: "path/to/file",
		CreatedAt:  now,
	}
}

func createTestAsset(t *testing.T, pool *pgxpool.Pool) domain.Asset {
	t.Helper()

	ctx := context.Background()

	userRepo := NewUserRepository(pool)
	assetRepo := NewAssetRepository(pool)

	user := validUser()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create test user: %v", err)
	}

	asset := validAsset()
	asset.UploadedBy = user.ID
	if err := assetRepo.Create(ctx, asset); err != nil {
		t.Fatalf("create test asset: %v", err)
	}

	return asset
}

func TestAssetArtifact_CreateAndGetByID(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	asset := createTestAsset(t, db)
	ctx := context.Background()

	artifactRepo := NewAssetArtifactRepository(db)
	want := validAssetArtifact()
	want.AssetID = asset.ID
	if err := artifactRepo.Create(ctx, want); err != nil {
		t.Fatalf("Create() asset artifact: %v", err)
	}

	got, err := artifactRepo.GetByID(ctx, want.ID)
	if err != nil {
		t.Fatalf("GetByID() asset artifact: %v", err)
	}

	if got.ID != want.ID {
		t.Fatalf("want id %q but got %q", want.ID, got.ID)
	}

	if got.AssetID != want.AssetID {
		t.Fatalf("want asset id %q but got %q", want.AssetID, got.AssetID)
	}

	if got.Type != want.Type {
		t.Fatalf("want type %q but got %q", want.Type, got.Type)
	}

	if got.MimeType != want.MimeType {
		t.Fatalf("want mime type %q but got %q", want.MimeType, got.MimeType)
	}

	if got.Size != want.Size {
		t.Fatalf("want size %d but got %d", want.Size, got.Size)
	}

	if got.Checksum != want.Checksum {
		t.Fatalf("want checksum %q but got %q", want.Checksum, got.Checksum)
	}

	if got.StorageKey != want.StorageKey {
		t.Fatalf("want storage key %q but got %q", want.StorageKey, got.StorageKey)
	}

	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Fatalf("want created at %q but got %q", want.CreatedAt, got.CreatedAt)
	}
}

func TestAssetArtifact_ListByAssetID(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	asset := createTestAsset(t, db)

	ctx := context.Background()

	artifactRepo := NewAssetArtifactRepository(db)

	art1 := validAssetArtifact()
	art2 := validAssetArtifact()
	art3 := validAssetArtifact()

	art1.AssetID = asset.ID
	if err := artifactRepo.Create(ctx, art1); err != nil {
		t.Fatalf("Create() artifact1 error: %v", err)
	}

	art2.AssetID = asset.ID
	if err := artifactRepo.Create(ctx, art2); err != nil {
		t.Fatalf("Create() artifact2 error: %v", err)
	}

	art3.AssetID = asset.ID
	if err := artifactRepo.Create(ctx, art3); err != nil {
		t.Fatalf("Create() artifact3 error: %v", err)
	}

	artifacts, err := artifactRepo.ListByAssetID(ctx, asset.ID)
	if err != nil {
		t.Fatalf("ListByAssetID() error: %v", err)
	}

	if len(artifacts) != 3 {
		t.Fatalf("expected len 3 but got %d", len(artifacts))
	}

}

func TestAssetArtifact_Delete(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	asset := createTestAsset(t, db)

	ctx := context.Background()

	artifactRepo := NewAssetArtifactRepository(db)

	art := validAssetArtifact()
	art.AssetID = asset.ID

	if err := artifactRepo.Create(ctx, art); err != nil {
		t.Fatalf("Create() artifact1 error: %v", err)
	}

	if err := artifactRepo.Delete(ctx, art.ID); err != nil {
		t.Fatalf("Delete() artifact: %v", err)
	}

	_, err := artifactRepo.GetByID(ctx, art.ID)
	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("GetByID() after Delete() artifact: %v", err)
	}
}

func TestAssetArtifact_GetByID_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	artifactRepo := NewAssetArtifactRepository(db)

	_, err := artifactRepo.GetByID(ctx, domain.AssetArtifactID(uuid.NewString()))
	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("GetByID() not found: expected ErrNotFound but got: %v", err)
	}
}
