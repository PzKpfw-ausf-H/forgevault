package minio

import (
	"context"
	"net/url"
	"time"
)

func (c *Client) PresignPut(ctx context.Context, bucket, key string, ttl time.Duration) (string, error) {
	u, err := c.mc.PresignedPutObject(ctx, bucket, key, ttl)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (c *Client) PresignGet(ctx context.Context, bucket, key string, ttl time.Duration) (string, error) {
	u, err := c.mc.PresignedGetObject(ctx, bucket, key, ttl, url.Values{})
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
