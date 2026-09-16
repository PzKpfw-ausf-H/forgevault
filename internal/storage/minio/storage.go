package minio

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/PzKpfw-ausf-H/forgevault/internal/storage"
	"github.com/minio/minio-go/v7"
)

type Storage struct {
	client *minio.Client
	bucket string
}

var _ storage.ObjectStorage = (*Storage)(nil)

func NewStorage(client *minio.Client, bucket string) *Storage {
	return &Storage{
		client: client,
		bucket: bucket,
	}
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	if err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("storage: minio: delete object: %w", err)
	}

	return nil
}

func (s *Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		resp := minio.ToErrorResponse(err)

		switch resp.Code {
		case "NoSuchKey":
			return false, nil
		default:
			return false, fmt.Errorf("storage: minio: exists: %w", err)
		}
	}

	return true, nil
}

func (s *Storage) PresignDownload(ctx context.Context, key string, ttl time.Duration) (string, error) {
	presigned, err := s.client.PresignedGetObject(ctx, s.bucket, key, ttl, url.Values{})
	if err != nil {
		return "", fmt.Errorf("storage: minio: presign download: %w", err)
	}

	return presigned.String(), nil
}

func (s *Storage) PresignUpload(ctx context.Context, key string, ttl time.Duration) (string, error) {
	presigned, err := s.client.PresignedPutObject(ctx, s.bucket, key, ttl)
	if err != nil {
		return "", fmt.Errorf("storage: minio: presign upload: %w", err)
	}

	return presigned.String(), nil
}
