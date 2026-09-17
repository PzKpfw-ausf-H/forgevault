package config

import (
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	Database DatabaseConfig
	S3       S3Config
}

type DatabaseConfig struct {
	URL string
}

type S3Config struct {
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
}

func Load() (Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, fmt.Errorf("config: DATABASE_URL: %w", ErrNotSet)
	}

	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		return Config{}, fmt.Errorf("config: S3_ENDPOINT: %w", ErrNotSet)
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return Config{}, fmt.Errorf("config: parse endpoint: %w", err)
	}

	if u.Host == "" || u.Scheme == "" {
		return Config{}, fmt.Errorf("config: S3_ENDPOINT: invalid URL")
	}

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		return Config{}, fmt.Errorf("config: S3_BUCKET: %w", ErrNotSet)
	}

	accessKey := os.Getenv("S3_ACCESS_KEY")
	if accessKey == "" {
		return Config{}, fmt.Errorf("config: S3_ACCESS_KEY: %w", ErrNotSet)
	}

	secretKey := os.Getenv("S3_SECRET_KEY")
	if secretKey == "" {
		return Config{}, fmt.Errorf("config: S3_SECRET_KEY: %w", ErrNotSet)
	}

	db := DatabaseConfig{
		URL: dbURL,
	}

	s3 := S3Config{
		Endpoint:  endpoint,
		Bucket:    bucket,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	return Config{
		Database: db,
		S3:       s3,
	}, nil
}
