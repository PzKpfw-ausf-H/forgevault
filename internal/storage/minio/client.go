package minio

import (
	"context"
	"net/url"
	"os"
	"strconv"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	mc *minio.Client
}

func NewFromEnv() (*Client, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		return nil, ErrMissingEndpoint
	}
	accessKey := os.Getenv("S3_ACCESS_KEY")
	if accessKey == "" {
		return nil, ErrMissingAccess
	}
	secretKey := os.Getenv("S3_SECRET_KEY")
	if secretKey == "" {
		return nil, ErrMissingSecret
	}
	useSSLStr := os.Getenv("S3_USE_SSL")

	useSSL, _ := strconv.ParseBool(useSSLStr)

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	mc, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	return &Client{mc: mc}, nil
}

func (c *Client) EnsureBucket(ctx context.Context, bucket string) error {
	exists, err := c.mc.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		err := c.mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
