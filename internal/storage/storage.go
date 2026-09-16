package storage

import (
	"context"
	"time"
)

type PresignUploadParams struct {
	Key         string
	ContentType string
}

type ObjectStorage interface {
	PresignDownload(ctx context.Context, key string, ttl time.Duration) (string, error)
	PresignUpload(ctx context.Context, params PresignUploadParams, ttl time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
