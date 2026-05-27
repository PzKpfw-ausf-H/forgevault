package service

import (
	"context"
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
}

func NewFileService(repo repo.AssetFilesRepo, bucket string, s3 *minio.Client, ttl time.Duration) *FileService {
	return &FileService{
		repo:   repo,
		bucket: bucket,
		s3:     s3,
		ttl:    ttl,
	}
}

func (fs *FileService) GetUploadURL(ctx context.Context, userID domain.UserID, assetID domain.AssetID, filename, contentType string) (version int, storageKey, uploadURL string, expiresInSec int64, err error) {
	max, err := fs.repo.GetMaxVersion(ctx, assetID)
	if err != nil {
		return 0, "", "", 0, err
	}

	version = max + 1

	storageKey = "assets/" + string(assetID) + "/v" + strconv.Itoa(version) + "/" + filename

	uploadURL, err = fs.s3.PresignGet(ctx, fs.bucket, storageKey, fs.ttl)
	if err != nil {
		return 0, "", "", 0, err
	}

	return version, storageKey, uploadURL, int64(fs.ttl.Minutes() * 60), nil
}

func (fs *FileService) ConfirmUpload(ctx context.Context, userID domain.UserID, assetID domain.AssetID, version int,
	filename, contentType string, sizeBytes int64, storageKey, checksum string) (domain.AssetFile, error) {
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

	err := fs.repo.Create(ctx, file)
	if err != nil {
		return domain.AssetFile{}, nil
	}

	return file, nil
}

func (fs *FileService) GetDownloadURL(ctx context.Context, userID domain.UserID, assetID domain.AssetID, version int) (downloadURL string, expiresInSec int64, err error) {
	file, err := fs.repo.GetByAssetVersion(ctx, assetID, version)
	if err != nil {
		return "", 0, err
	}

	downloadURL, err = fs.s3.PresignGet(ctx, fs.bucket, file.StorageKey, fs.ttl)
	if err != nil {
		return "", 0, err
	}

	return downloadURL, int64(fs.ttl.Minutes() * 60), nil
}
