package postgres

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/google/uuid"
)

func validAsset() domain.Asset {
	now := time.Now().UTC().Truncate(time.Microsecond)

	return domain.Asset{
		ID:          domain.AssetID(uuid.NewString()),
		Title:       "TestAsset",
		Description: "something to test",

		AuthorName: "Student A",
		UploadedBy: domain.UserID(uuid.NewString()),

		Type:             domain.AssetType3D,
		Tags:             []string{"tag1", "tag2", "tag3"},
		ProcessingStatus: domain.ProcessingStatusReady,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

func assertAssetEqual(t *testing.T, got, want domain.Asset) {
	t.Helper()

	if got.ID != want.ID {
		t.Errorf("expected id %q but got %q", want.ID, got.ID)
	}

	if got.Title != want.Title {
		t.Errorf("expected title %q but got %q", want.Title, got.Title)
	}

	if got.Description != want.Description {
		t.Errorf("expected description %q but got %q", want.Description, got.Description)
	}

	if got.AuthorName != want.AuthorName {
		t.Errorf("expected author name %q but got %q", want.AuthorName, got.AuthorName)
	}

	if got.UploadedBy != want.UploadedBy {
		t.Errorf("expected id %q but got %q", want.UploadedBy, got.UploadedBy)
	}

	if got.Type != want.Type {
		t.Errorf("expected type %q but got %q", want.Type, got.Type)
	}

	if !slices.Equal(got.Tags, want.Tags) {
		t.Errorf("expected tags %v but got %v", want.Tags, got.Tags)
	}

	if got.ProcessingStatus != want.ProcessingStatus {
		t.Errorf("expected processing status %q but got %q", want.ProcessingStatus, got.ProcessingStatus)
	}

	if !got.CreatedAt.Equal(want.CreatedAt) {
		t.Errorf("expected created at %v but got %v", want.CreatedAt, got.CreatedAt)
	}

	if !got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Errorf("expected updated at %v but got %v", want.UpdatedAt, got.UpdatedAt)
	}
}

func TestAssetRepository_CreateAndGetByID(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	// copying the createTestAsset logic instead of calling it because i need the asset repo
	userRepo := NewUserRepository(db)
	assetRepo := NewAssetRepository(db)

	user := validUser()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	want := validAsset()
	want.UploadedBy = user.ID
	if err := assetRepo.Create(ctx, want); err != nil {
		t.Fatalf("create asset: %v", err)
	}

	got, err := assetRepo.GetByID(ctx, want.ID)
	if err != nil {
		t.Fatalf("get asset by id: %v", err)
	}

	assertAssetEqual(t, got, want)
}

func TestAssetRepository_GetByID_NotFound(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	assetRepo := NewAssetRepository(db)

	_, err := assetRepo.GetByID(ctx, "99999999-9999-9999-9999-999999999999")
	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("get by id not found: expected ErrNotFound but got: %v", err)
	}
}

func TestAssetRepository_Update(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	userRepo := NewUserRepository(db)
	assetRepo := NewAssetRepository(db)

	user := validUser()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	asset := validAsset()
	asset.UploadedBy = user.ID

	asset.Title = "Old title"
	asset.Description = "Old description"
	asset.Tags = []string{"old", "medieval"}

	originalUploadedBy := asset.UploadedBy
	originalCreatedAt := asset.CreatedAt

	newTitle := "New title"
	newDesc := "New desc"
	newTags := []string{"building", "new"}
	newUpdatedAt := asset.UpdatedAt.Add(time.Minute)

	if err := assetRepo.Create(ctx, asset); err != nil {
		t.Fatalf("create asset: %v", err)
	}

	asset.Title = newTitle
	asset.Description = newDesc
	asset.Tags = newTags
	asset.UpdatedAt = newUpdatedAt

	if err := assetRepo.Update(ctx, asset); err != nil {
		t.Fatalf("update asset: %v", err)
	}

	got, err := assetRepo.GetByID(ctx, asset.ID)
	if err != nil {
		t.Fatalf("get asset by id updated: %v", err)
	}

	if got.Title != newTitle {
		t.Errorf("expected title %q but got %q", newTitle, got.Title)
	}

	if got.Description != newDesc {
		t.Errorf("expected description %q but got %q", newDesc, got.Description)
	}

	if !slices.Equal(got.Tags, newTags) {
		t.Errorf("expected tags %v but got %v", newTags, got.Tags)
	}

	if !got.UpdatedAt.Equal(newUpdatedAt) {
		t.Errorf("expected updated at %v but got %v", newUpdatedAt, got.UpdatedAt)
	}

	if got.UploadedBy != originalUploadedBy {
		t.Errorf("expected uploaded by %q but got %q", originalUploadedBy, got.UploadedBy)
	}

	if !got.CreatedAt.Equal(originalCreatedAt) {
		t.Errorf("expected created at %v but got %v", originalCreatedAt, got.CreatedAt)
	}
}

func TestAssetRepository_Delete(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	userRepo := NewUserRepository(db)
	assetRepo := NewAssetRepository(db)

	user := validUser()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	want := validAsset()
	want.UploadedBy = user.ID
	if err := assetRepo.Create(ctx, want); err != nil {
		t.Fatalf("create asset: %v", err)
	}

	if err := assetRepo.Delete(ctx, want.ID); err != nil {
		t.Fatalf("delete asset: %v", err)
	}

	_, err := assetRepo.GetByID(ctx, want.ID)
	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("get by id deleted: expected ErrNotFound but got: %v", err)
	}

	// checking ON DELETE CASCADE on asset_tags table
	var count int

	err = db.QueryRow(ctx,
		`SELECT COUNT(*)
		FROM asset_tags
		WHERE asset_id = $1`,
		want.ID,
	).Scan(&count)

	if err != nil {
		t.Fatalf("query row on asset tags: %v", err)
	}

	if count != 0 {
		t.Fatalf("expected count 0 but got %d", count)
	}
}

func TestAssetRepository_Create_RollbackOnTagError(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	userRepo := NewUserRepository(db)
	assetRepo := NewAssetRepository(db)

	user := validUser()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	asset := validAsset()
	asset.UploadedBy = user.ID
	asset.Tags = []string{"tag1", "tag1"}

	err := assetRepo.Create(ctx, asset)
	if err == nil {
		t.Fatalf("expected Create() error")
	}

	_, err = assetRepo.GetByID(ctx, asset.ID)

	if !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("get by id after rollback: expected ErrNotFound but got: %v", err)
	}
}

func TestAssetRepository_List_Pagination(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()
	baseTime := time.Now().UTC().Truncate(time.Microsecond)

	userRepo := NewUserRepository(db)
	assetRepo := NewAssetRepository(db)

	user := validUser()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// Creating assets
	asset1 := validAsset()
	asset1.UploadedBy = user.ID
	asset1.CreatedAt = baseTime
	if err := assetRepo.Create(ctx, asset1); err != nil {
		t.Fatalf("create asset1: %v", err)
	}

	asset2 := validAsset()
	asset2.UploadedBy = user.ID
	asset2.CreatedAt = baseTime.Add(time.Minute)
	if err := assetRepo.Create(ctx, asset2); err != nil {
		t.Fatalf("create asset2: %v", err)
	}

	asset3 := validAsset()
	asset3.UploadedBy = user.ID
	asset3.CreatedAt = baseTime.Add(2 * time.Minute)
	if err := assetRepo.Create(ctx, asset3); err != nil {
		t.Fatalf("create asset3: %v", err)
	}

	asset4 := validAsset()
	asset4.UploadedBy = user.ID
	asset4.CreatedAt = baseTime.Add(3 * time.Minute)
	if err := assetRepo.Create(ctx, asset4); err != nil {
		t.Fatalf("create asset4: %v", err)
	}

	asset5 := validAsset()
	asset5.UploadedBy = user.ID
	asset5.CreatedAt = baseTime.Add(4 * time.Minute)
	if err := assetRepo.Create(ctx, asset5); err != nil {
		t.Fatalf("create asset5: %v", err)
	}
	// End of creating assets

	want := []domain.Asset{asset4, asset3}

	result, err := assetRepo.List(ctx, repo.AssetListParams{
		Limit:  2,
		Offset: 1,
	})
	if err != nil {
		t.Fatalf("list assets: %v", err)
	}

	if result.Total != 5 {
		t.Fatalf("expected total 5 but got %d", result.Total)
	}

	if len(result.Assets) != len(want) {
		t.Fatalf("expected assets len %d but gut %d", len(want), len(result.Assets))
	}

	for i := range result.Assets {
		assertAssetEqual(t, result.Assets[i], want[i])
	}
}

func TestAssetRepository_List_TypeFilter(t *testing.T) {
	db := openTestDB(t)
	cleanTestDB(t, db)

	ctx := context.Background()

	userRepo := NewUserRepository(db)
	assetRepo := NewAssetRepository(db)

	user := validUser()
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// Creating assets
	asset1 := validAsset()
	asset1.UploadedBy = user.ID
	asset1.Type = domain.AssetType3D
	if err := assetRepo.Create(ctx, asset1); err != nil {
		t.Fatalf("create asset1: %v", err)
	}

	asset2 := validAsset()
	asset2.UploadedBy = user.ID
	asset2.Type = domain.AssetType3D
	if err := assetRepo.Create(ctx, asset2); err != nil {
		t.Fatalf("create asset2: %v", err)
	}

	asset3 := validAsset()
	asset3.UploadedBy = user.ID
	asset3.Type = domain.AssetTypeAudio
	if err := assetRepo.Create(ctx, asset3); err != nil {
		t.Fatalf("create asset3: %v", err)
	}

	asset4 := validAsset()
	asset4.UploadedBy = user.ID
	asset4.Type = domain.AssetType2D
	if err := assetRepo.Create(ctx, asset4); err != nil {
		t.Fatalf("create asset4: %v", err)
	}
	// End of creating assets

	assetType := domain.AssetType3D

	result, err := assetRepo.List(ctx, repo.AssetListParams{
		Limit:  10,
		Offset: 0,
		Type:   &assetType,
	})

	if err != nil {
		t.Fatalf("list with params: %v", err)
	}

	if result.Total != 2 {
		t.Fatalf("expected total 2 but got %d", result.Total)
	}

	if len(result.Assets) != 2 {
		t.Fatalf("expected result assets len 2 but got %d", len(result.Assets))
	}

	for _, asset := range result.Assets {
		if asset.Type != domain.AssetType3D {
			t.Fatalf("expected type %q but got %q", domain.AssetType3D, asset.Type)
		}
	}

}
