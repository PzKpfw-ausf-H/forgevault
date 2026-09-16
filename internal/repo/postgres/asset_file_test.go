package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/google/uuid"
)

func ptr[T any](v T) *T {
	return &v
}

func validMainAssetFile(assetID domain.AssetID) domain.AssetFile {
	now := time.Now().UTC().Truncate(time.Microsecond)
	id := domain.AssetFileID(uuid.NewString())

	return domain.AssetFile{
		ID:           id,
		AssetID:      assetID,
		Role:         domain.FileRoleMain,
		TextureType:  nil,
		OriginalName: "main_model.blend",
		MimeType:     "application/json",
		Extension:    ".blend",
		Size:         160,
		StorageKey:   "files/" + string(id),
		Checksum:     "not empty",
		CreatedAt:    now,
	}
}

func validTextureAssetFile(assetID domain.AssetID) domain.AssetFile {
	now := time.Now().UTC().Truncate(time.Microsecond)

	return domain.AssetFile{
		ID:           domain.AssetFileID(uuid.NewString()),
		AssetID:      assetID,
		Role:         domain.FileRoleTexture,
		TextureType:  ptr(domain.TextureTypeBaseColor),
		OriginalName: "base_color.png",
		MimeType:     "image/png",
		Extension:    ".png",
		Size:         10,
		StorageKey:   "path/to/file",
		Checksum:     "not empty",
		CreatedAt:    now,
	}
}

func assertAssetFileEqual(t *testing.T, got, want domain.AssetFile) {
	t.Helper()

	if got.ID != want.ID {
		t.Errorf("want id %q but got %q", want.ID, got.ID)
	}

	if got.AssetID != want.AssetID {
		t.Errorf("want asset id %q but got %q", want.AssetID, got.AssetID)
	}

	if got.Role != want.Role {
		t.Errorf("want role %q but got %q", want.Role, got.Role)
	}

	switch {
	case got.TextureType == nil && want.TextureType == nil:

	case got.TextureType == nil || want.TextureType == nil:
		t.Errorf("want texture type %v but got %v", want.TextureType, got.TextureType)
	case *got.TextureType != *want.TextureType:
		t.Errorf("want texture type %q but got %q", *want.TextureType, *got.TextureType)
	}

	if got.OriginalName != want.OriginalName {
		t.Errorf("want original name %q but got %q", want.OriginalName, got.OriginalName)
	}

	if got.MimeType != want.MimeType {
		t.Errorf("want mime type %q but got %q", want.MimeType, got.MimeType)
	}

	if got.Extension != want.Extension {
		t.Errorf("want extension %q but got %q", want.Extension, got.Extension)
	}

	if got.Size != want.Size {
		t.Errorf("want size %d but got %d", want.Size, got.Size)
	}

	if got.StorageKey != want.StorageKey {
		t.Errorf("want storage key %q but got %q", want.StorageKey, got.StorageKey)
	}

	if got.Checksum != want.Checksum {
		t.Errorf("want checksum %q but got %q", want.Checksum, got.Checksum)
	}

	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("want created at %q but got %q", want.CreatedAt, got.CreatedAt)
	}
}

func TestAssetFile_CreateAndGetByID(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	asset := createTestAsset(t, db)
	ctx := context.Background()

	assetFileRepo := NewAssetFileRepository(db)

	want := validMainAssetFile(asset.ID)
	if err := assetFileRepo.Create(ctx, want); err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	got, err := assetFileRepo.GetByID(ctx, want.ID)
	if err != nil {
		t.Fatalf("GetByID() error: %v", err)
	}

	assertAssetFileEqual(t, got, want)
}

func TestAssetFile_ListByAssetID(t *testing.T) {
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

	asset1 := validAsset()
	asset1.UploadedBy = user.ID
	if err := assetRepo.Create(ctx, asset1); err != nil {
		t.Fatalf("create test asset: %v", err)
	}

	asset2 := validAsset()
	asset2.UploadedBy = user.ID
	if err := assetRepo.Create(ctx, asset2); err != nil {
		t.Fatalf("create test asset: %v", err)
	}
	// End of creating assets

	// Creating asset files
	assetFileRepo := NewAssetFileRepository(db)

	mainFile := validMainAssetFile(asset1.ID)
	mainFile.CreatedAt = baseTime
	if err := assetFileRepo.Create(ctx, mainFile); err != nil {
		t.Fatalf("create main file: %v", err)
	}

	textureFile := validTextureAssetFile(asset1.ID)
	textureFile.CreatedAt = baseTime.Add(time.Second)
	if err := assetFileRepo.Create(ctx, textureFile); err != nil {
		t.Fatalf("create texture file: %v", err)
	}

	mainFile2 := validMainAssetFile(asset2.ID)
	if err := assetFileRepo.Create(ctx, mainFile2); err != nil {
		t.Fatalf("create main file 2: %v", err)
	}
	// End of creating asset files

	want1 := []domain.AssetFile{mainFile, textureFile}
	want2 := []domain.AssetFile{mainFile2}

	got1, err := assetFileRepo.ListByAssetID(ctx, asset1.ID)
	if err != nil {
		t.Fatalf("list by id: got1: %v", err)
	}

	got2, err := assetFileRepo.ListByAssetID(ctx, asset2.ID)
	if err != nil {
		t.Fatalf("list by id: got2: %v", err)
	}

	if len(want1) != len(got1) {
		t.Fatalf("got1: expected len %d but got %d", len(want1), len(got1))
	}

	for i := range got1 {
		assertAssetFileEqual(t, got1[i], want1[i])
	}

	if len(got2) != 1 {
		t.Fatalf("got2: expected len 1 but got %d", len(got2))
	}

	assertAssetFileEqual(t, got2[0], want2[0])
}

func TestAssetFile_Delete(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	asset := createTestAsset(t, db)
	ctx := context.Background()

	assetFileRepo := NewAssetFileRepository(db)

	file := validMainAssetFile(asset.ID)
	if err := assetFileRepo.Create(ctx, file); err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	if err := assetFileRepo.Delete(ctx, file.ID); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	_, err := assetFileRepo.GetByID(ctx, file.ID)
	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("delete: expected ErrNotFound but got: %v", err)
	}
}

func TestAssetFile_GetByID_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	assetFileRepo := NewAssetFileRepository(db)

	_, err := assetFileRepo.GetByID(ctx, domain.AssetFileID(uuid.NewString()))
	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("get by id not found: expected ErrNotFound but got %v", err)
	}
}

func TestAssetFile_Delete_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	assetFileRepo := NewAssetFileRepository(db)

	err := assetFileRepo.Delete(ctx, domain.AssetFileID(uuid.NewString()))
	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("delete not found: expected ErrNotFound but got %v", err)
	}
}

func TestAssetFile_TextureType_RoundTrip(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	asset := createTestAsset(t, db)

	ctx := context.Background()

	assetFileRepo := NewAssetFileRepository(db)
	want := validTextureAssetFile(asset.ID)
	if err := assetFileRepo.Create(ctx, want); err != nil {
		t.Fatalf("create() error: %v", err)
	}

	got, err := assetFileRepo.GetByID(ctx, want.ID)
	if err != nil {
		t.Fatalf("get by id(): %v", err)
	}

	if got.TextureType == nil {
		t.Fatalf("expected texture type to be not nil")
	}

	assertAssetFileEqual(t, got, want)

}

func TestAssetFile_SecondMainRejected(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	asset := createTestAsset(t, db)

	assetFileRepo := NewAssetFileRepository(db)

	main1 := validMainAssetFile(asset.ID)
	main2 := validMainAssetFile(asset.ID)

	if err := assetFileRepo.Create(ctx, main1); err != nil {
		t.Fatalf("Create main1: %v", err)
	}

	err := assetFileRepo.Create(ctx, main2)
	if !errors.Is(err, repo.ErrAlreadyExists) {
		t.Fatalf("create main2: expected ErrAlreadyExists but got: %v", err)
	}
}
