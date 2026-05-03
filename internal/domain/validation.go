package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var tagRe = regexp.MustCompile(`^[a-z0-9\-]+$`)

func NormalizeTags(tags []string) ([]string, error) {
	tagsMap := make(map[string]struct{})
	out := make([]string, 0, len(tags))

	for _, tag := range tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		if len(tag) < 2 || len(tag) > 32 {
			return nil, fmt.Errorf("tag: %w", ErrInvalidTagLen)
		}
		if !tagRe.MatchString(tag) {
			return nil, fmt.Errorf("tag: %w", ErrInvalidTag)
		}
		if _, ok := tagsMap[tag]; ok {
			continue
		}
		tagsMap[tag] = struct{}{}
		out = append(out, tag)
	}

	return out, nil
}

func ValidateNewUser(u *User) error {
	u.ID = UserID(strings.TrimSpace(string(u.ID)))
	if u.ID == "" {
		return ErrInvalidID
	}
	if err := ValidateEmail(u.Email); err != nil {
		return ErrInvalidEmail
	}

	return nil
}

func ValidateEmail(email string) error {
	em := strings.ToLower(strings.TrimSpace(email))
	if em == "" {
		return ErrInvalidEmail
	}
	if match, err := regexp.MatchString(`@`, em); err != nil || !match {
		return ErrInvalidEmail
	}

	return nil
}

func ValidateNewAsset(a *Asset) error {
	a.Title = strings.TrimSpace(a.Title)
	if a.Title == "" {
		return fmt.Errorf("asset title: %w", ErrInvalidTitle)
	}
	a.ID = AssetID(strings.TrimSpace(string(a.ID)))
	if a.ID == "" {
		return fmt.Errorf("ID: %w", ErrInvalidID)
	}
	if !IsValidAssetType(a.Type) {
		return fmt.Errorf("asset type: %w", ErrInvalidAssetType)
	}
	a.AuthorID = UserID(strings.TrimSpace(string(a.AuthorID)))
	if a.AuthorID == "" {
		return fmt.Errorf("author id: %w", ErrInvalidID)
	}
	validatedTags, tagsErr := NormalizeTags(a.Tags)
	if tagsErr != nil {
		return tagsErr
	}
	a.Tags = validatedTags

	return nil
}

func ValidateNewAssetFile(f *AssetFile) error {
	f.ID = FileID(strings.TrimSpace(string(f.ID)))
	if f.ID == "" {
		return fmt.Errorf("asset file id: %w", ErrInvalidID)
	}
	f.AssetID = AssetID(strings.TrimSpace(string(f.AssetID)))
	if f.AssetID == "" {
		return fmt.Errorf("asset file asset id: %w", ErrInvalidID)
	}
	if f.Version <= 0 {
		return fmt.Errorf("asset file version: %w", ErrInvalidVersion)
	}
	if strings.TrimSpace(f.Filename) == "" {
		return fmt.Errorf("asset file filename: %w", ErrInvalidFilename)
	}
	if f.SizeBytes <= 0 {
		return fmt.Errorf("asset file size bytes: %w", ErrInvalidSize)
	}
	if strings.TrimSpace(f.ContentType) == "" {
		return fmt.Errorf("asset file content type: %w", ErrInvalidContentType)
	}
	if strings.TrimSpace(f.StorageKey) == "" {
		return fmt.Errorf("asset file storage key: %w", ErrInvalidStorageKey)
	}

	return nil
}
