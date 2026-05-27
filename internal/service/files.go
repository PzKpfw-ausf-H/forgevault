package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/domain"
	"github.com/PzKpfw-ausf-H/forgevault/internal/repo"
	"github.com/PzKpfw-ausf-H/forgevault/internal/storage/minio"
	"github.com/google/uuid"
)

type FileService struct {
	repo   repo.AssetFilesRepo
	bucket string
	s3     *minio.Client
	ttl    time.Duration
	assets repo.AssetsRepo
}

func NewFileService(repo repo.AssetFilesRepo, bucket string, s3 *minio.Client, ttl time.Duration, assets repo.AssetsRepo) *FileService {
	return &FileService{
		repo:   repo,
		bucket: bucket,
		s3:     s3,
		ttl:    ttl,
		assets: assets,
	}
}

func (fs *FileService) GetUploadURL(ctx context.Context, userID domain.UserID, assetID domain.AssetID, filename, contentType string) (version int, storageKey, uploadURL string, expiresInSec int64, err error) {
	a, err := fs.assets.GetByID(ctx, assetID)
	if err != nil {
		return 0, "", "", 0, err
	}
	if a.AuthorID != userID {
		return 0, "", "", 0, repo.ErrUnauthorized
	}

	if filename == "" {
		return 0, "", "", 0, fmt.Errorf("filename: %w", repo.ErrValidation)
	}
	if contentType == "" {
		return 0, "", "", 0, fmt.Errorf("content type: %w", repo.ErrValidation)
	}

	max, err := fs.repo.GetMaxVersion(ctx, assetID)
	if err != nil {
		return 0, "", "", 0, err
	}

	version = max + 1

	storageKey = "assets/" + string(assetID) + "/v" + strconv.Itoa(version) + "/" + filename

	uploadURL, err = fs.s3.PresignPut(ctx, fs.bucket, storageKey, fs.ttl)
	if err != nil {
		return 0, "", "", 0, err
	}

	return version, storageKey, uploadURL, int64(fs.ttl.Seconds()), nil
}

func (fs *FileService) ConfirmUpload(ctx context.Context, userID domain.UserID, assetID domain.AssetID, version int,
	filename, contentType string, sizeBytes int64, storageKey, checksum string) (domain.AssetFile, error) {
	a, err := fs.assets.GetByID(ctx, assetID)
	if err != nil {
		return domain.AssetFile{}, err
	}
	if a.AuthorID != userID {
		return domain.AssetFile{}, repo.ErrUnauthorized
	}

	if version <= 0 {
		return domain.AssetFile{}, fmt.Errorf("version: %w", repo.ErrValidation)
	}
	if sizeBytes <= 0 {
		return domain.AssetFile{}, fmt.Errorf("size bytes: %w", repo.ErrValidation)
	}

	expected := "assets/" + string(assetID) + "/v" + strconv.Itoa(version) + "/" + filename

	if storageKey == "" {
		return domain.AssetFile{}, fmt.Errorf("storage key: %w", repo.ErrValidation)
	}
	if storageKey != expected {
		return domain.AssetFile{}, repo.ErrBadRequest
	}

	if filename == "" {
		return domain.AssetFile{}, fmt.Errorf("filename: %w", repo.ErrValidation)
	}
	if contentType == "" {
		return domain.AssetFile{}, fmt.Errorf("content type: %w", repo.ErrValidation)
	}

	fileID := uuid.New()
	file := domain.AssetFile{
		ID:          domain.FileID(fileID.String()),
		AssetID:     assetID,
		Version:     version,
		Filename:    filename,
		SizeBytes:   sizeBytes,
		ContentType: contentType,
		StorageKey:  storageKey,
		Checksum:    checksum,
		CreatedAt:   time.Now().UTC(),
	}

	err = fs.repo.Create(ctx, file)
	if err != nil {
		return domain.AssetFile{}, err
	}

	return file, nil
}

func (fs *FileService) GetDownloadURL(ctx context.Context, userID domain.UserID, assetID domain.AssetID, version int) (downloadURL string, expiresInSec int64, err error) {
	a, err := fs.assets.GetByID(ctx, assetID)
	if err != nil {
		return "", 0, err
	}
	if a.AuthorID != userID {
		return "", 0, repo.ErrUnauthorized
	}

	file, err := fs.repo.GetByAssetVersion(ctx, assetID, version)
	if err != nil {
		return "", 0, err
	}

	downloadURL, err = fs.s3.PresignGet(ctx, fs.bucket, file.StorageKey, fs.ttl)
	if err != nil {
		return "", 0, err
	}

	return downloadURL, int64(fs.ttl.Seconds()), nil
}
