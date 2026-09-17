package minio

import (
	"fmt"
	"net/url"

	"github.com/PzKpfw-ausf-H/forgevault/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewClient(cfg config.S3Config) (*minio.Client, error) {
	u, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("storage: minio client: endpoint parse: %w", err)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("storage: minio client: invalid host")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("storage: minio client: invalid scheme")
	}

	useSSL := false

	if u.Scheme == "https" {
		useSSL = true
	}

	client, err := minio.New(u.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: minio client: new client: %w", err)
	}

	return client, nil
}
