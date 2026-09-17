package minio

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

//TEST_S3_ENDPOINT, TEST_S3_BUCKET, TEST_S3_ACCESS_KEY, TEST_S3_SECRET_KEY

func openTestStorage(t *testing.T) *Storage {
	t.Helper()

	endpoint := os.Getenv("TEST_S3_ENDPOINT")
	if endpoint == "" {
		t.Skip("TEST_S3_ENDPOINT is not set")
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		t.Fatalf("open test storage: parse endpoint: %v", err)
	}

	bucket := os.Getenv("TEST_S3_BUCKET")
	if bucket == "" {
		t.Skip("TEST_S3_BUCKET is not set")
	}

	access := os.Getenv("TEST_S3_ACCESS_KEY")
	if access == "" {
		t.Skip("TEST_S3_ACCESS_KEY is not set")
	}

	secret := os.Getenv("TEST_S3_SECRET_KEY")
	if secret == "" {
		t.Skip("TEST_S3_SECRET_KEY is not set")
	}

	useSSL := false

	client, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(access, secret, ""),
		Secure: useSSL,
	})
	if err != nil {
		t.Fatalf("open test storage: new client: %v", err)
	}

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatalf("open test storage: bucket exists: %v", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			t.Fatalf("open test storage: make bucket: %v", err)
		}
	}

	return NewStorage(client, bucket)
}
