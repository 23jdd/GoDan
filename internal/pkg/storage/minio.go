package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"godan/internal/config"
)

type minioStore struct {
	client    *minio.Client
	bucket    string
	urlExpire time.Duration
}

func newMinIO(cfg *config.MinIOConfig) (*minioStore, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio connect error: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("minio bucket check error: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio create bucket error: %w", err)
		}
	}

	return &minioStore{
		client:    client,
		bucket:    cfg.Bucket,
		urlExpire: time.Duration(cfg.URLExpire) * time.Second,
	}, nil
}

func (m *minioStore) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := m.client.PutObject(ctx, m.bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio upload error: %w", err)
	}

	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucket, key, m.urlExpire, reqParams)
	if err != nil {
		return "", fmt.Errorf("minio presigned url error: %w", err)
	}

	return presignedURL.String(), nil
}

func (m *minioStore) Delete(ctx context.Context, key string) error {
	return m.client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
}
