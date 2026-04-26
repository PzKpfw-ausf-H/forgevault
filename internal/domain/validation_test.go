package domain

import (
	"testing"
	"time"
)

func TestValidation_ValidateNewAsset(t *testing.T) {
	a := &Asset{
		ID:          "1",
		Title:       "TestItem",
		Description: "none yet",
		Type:        AssetType3D,
		Tags:        []string{"tag1", "Tag2", "tAg3"},
		AuthorID:    "1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := ValidateNewAsset(a); err != nil {
		t.Fatalf("err validating asset: %v", err)
	}
}

func TestValidation_NormalizeTagsToLower(t *testing.T) {
	tagsRaw := []string{"tag1", "TAG2", "tAg3"}
	newTags, err := NormalizeTags(tagsRaw)
	if err != nil {
		t.Fatalf("err normilizing tags: %v", err)
	}
	if newTags[1] != "tag2" {
		t.Fatalf("err to lower all")
	}
	if newTags[2] != "tag3" {
		t.Fatal("err to lower specific")
	}
}

func TestValidation_NormalizeTagsRemoveDups(t *testing.T) {
	tagsRaw := []string{"TAG", "tag", "tag2", "TAG2"}
	newTags, tagsErr := NormalizeTags(tagsRaw)
	if tagsErr != nil {
		t.Fatalf("err deleting duplicates: %v", tagsErr)
	}
	if len(newTags) != 2 {
		t.Fatal("err tags length")
	}
}

func TestValidation_ValidateNewAssetFile(t *testing.T) {
	f := &AssetFile{
		ID:          "1",
		AssetID:     "2",
		Version:     1,
		Filename:    "filename",
		SizeBytes:   15,
		ContentType: "contentType",
		StorageKey:  "storagekey",
		Checksum:    "checksum",
		CreatedAt:   time.Now(),
	}

	if err := ValidateNewAssetFile(f); err != nil {
		t.Fatalf("err validating asset file: %v", err)
	}
}
