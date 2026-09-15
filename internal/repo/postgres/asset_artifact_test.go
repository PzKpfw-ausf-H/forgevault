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

func validAssetArtifact(assetID domain.AssetID) domain.AssetArtifact {
	now := time.Now().UTC().Truncate(time.Microsecond)

	return domain.AssetArtifact{
		ID:         domain.AssetArtifactID(uuid.NewString()),
		AssetID:    assetID,
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

func assertArtifactEqual(t *testing.T, got, want domain.AssetArtifact) {
	t.Helper()

	if got.ID != want.ID {
		t.Errorf("want id %q but got %q", want.ID, got.ID)
	}

	if got.AssetID != want.AssetID {
		t.Errorf("want asset id %q but got %q", want.AssetID, got.AssetID)
	}

	if got.Type != want.Type {
		t.Errorf("want type %q but got %q", want.Type, got.Type)
	}

	if got.MimeType != want.MimeType {
		t.Errorf("want mime type %q but got %q", want.MimeType, got.MimeType)
	}

	if got.Size != want.Size {
		t.Errorf("want size %d but got %d", want.Size, got.Size)
	}

	if got.Checksum != want.Checksum {
		t.Errorf("want checksum %q but got %q", want.Checksum, got.Checksum)
	}

	if got.StorageKey != want.StorageKey {
		t.Errorf("want storage key %q but got %q", want.StorageKey, got.StorageKey)
	}

	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("want created at %q but got %q", want.CreatedAt, got.CreatedAt)
	}
}

func TestAssetArtifact_CreateAndGetByID(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	asset := createTestAsset(t, db)
	ctx := context.Background()

	artifactRepo := NewAssetArtifactRepository(db)
	want := validAssetArtifact(asset.ID)
	if err := artifactRepo.Create(ctx, want); err != nil {
		t.Fatalf("Create() asset artifact: %v", err)
	}

	got, err := artifactRepo.GetByID(ctx, want.ID)
	if err != nil {
		t.Fatalf("GetByID() asset artifact: %v", err)
	}

	assertArtifactEqual(t, got, want)
}

func TestAssetArtifact_ListByAssetID(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	baseTime := time.Now().UTC().Truncate(time.Microsecond)

	ctx := context.Background()

	// Creating assets
	userRepo := NewUserRepository(db)
	assetRepo := NewAssetRepository(db)

	user := validUser()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create test user: %v", err)
	}

	asset := validAsset()
	asset.UploadedBy = user.ID
	if err := assetRepo.Create(ctx, asset); err != nil {
		t.Fatalf("create test asset: %v", err)
	}

	asset2 := validAsset()
	asset2.UploadedBy = user.ID
	if err := assetRepo.Create(ctx, asset2); err != nil {
		t.Fatalf("create test asset: %v", err)
	}
	// End of creating assets

	// Creating asset artifacts
	artifactRepo := NewAssetArtifactRepository(db)

	art1 := validAssetArtifact(asset.ID)
	art1.CreatedAt = baseTime

	art2 := validAssetArtifact(asset.ID)
	art2.CreatedAt = baseTime.Add(time.Second)

	art3 := validAssetArtifact(asset.ID)
	art3.CreatedAt = baseTime.Add(2 * time.Second)

	art4 := validAssetArtifact(asset2.ID)

	if err := artifactRepo.Create(ctx, art1); err != nil {
		t.Fatalf("Create() artifact1 error: %v", err)
	}

	if err := artifactRepo.Create(ctx, art2); err != nil {
		t.Fatalf("Create() artifact2 error: %v", err)
	}

	if err := artifactRepo.Create(ctx, art3); err != nil {
		t.Fatalf("Create() artifact3 error: %v", err)
	}

	if err := artifactRepo.Create(ctx, art4); err != nil {
		t.Fatalf("Create() artifact4 error: %v", err)
	}

	// End of creating asset artifacts

	want1 := []domain.AssetArtifact{
		art1, art2, art3,
	}
	want2 := []domain.AssetArtifact{art4}

	got1, err := artifactRepo.ListByAssetID(ctx, asset.ID)
	if err != nil {
		t.Fatalf("ListByAssetID() got1 error: %v", err)
	}

	if len(got1) != 3 {
		t.Fatalf("got1: expected len 3 but got %d", len(got1))
	}

	for i := range got1 {
		assertArtifactEqual(t, got1[i], want1[i])
	}

	got2, err := artifactRepo.ListByAssetID(ctx, asset2.ID)
	if err != nil {
		t.Fatalf("ListByAssetID() got2 error: %v", err)
	}

	if len(got2) != 1 {
		t.Fatalf("got2: expected len 1 but got %d", len(got2))
	}

	assertArtifactEqual(t, got2[0], want2[0])
}

func TestAssetArtifact_Delete(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	asset := createTestAsset(t, db)

	ctx := context.Background()

	artifactRepo := NewAssetArtifactRepository(db)

	art := validAssetArtifact(asset.ID)

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

func TestAssetArtifact_Delete_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	artifactRepo := NewAssetArtifactRepository(db)

	err := artifactRepo.Delete(ctx, "99999999-9999-9999-9999-999999999999")

	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("Delete() not found: expected ErrNotFound but got: %v", err)
	}
}
