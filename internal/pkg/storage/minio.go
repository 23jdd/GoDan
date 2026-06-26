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
	core      *minio.Core
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

	core, err := minio.NewCore(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio core init error: %w", err)
	}

	return &minioStore{
		client:    client,
		core:      core,
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

func (m *minioStore) InitUpload(ctx context.Context, key string, contentType string) (string, error) {
	uploadID, err := m.core.NewMultipartUpload(ctx, m.bucket, key, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio init multipart: %w", err)
	}
	return uploadID, nil
}

func (m *minioStore) UploadPart(ctx context.Context, key string, uploadID string, partNumber int, reader io.Reader, size int64) error {
	_, err := m.core.PutObjectPart(ctx, m.bucket, key, uploadID, partNumber, reader, size, minio.PutObjectPartOptions{})
	if err != nil {
		return fmt.Errorf("minio upload part %d: %w", partNumber, err)
	}
	return nil
}

func (m *minioStore) CompleteUpload(ctx context.Context, key string, uploadID string, parts []UploadPart) (string, error) {
	var completeParts []minio.CompletePart
	for _, p := range parts {
		completeParts = append(completeParts, minio.CompletePart{
			PartNumber: p.PartNumber,
		})
	}

	_, err := m.core.CompleteMultipartUpload(ctx, m.bucket, key, uploadID, completeParts, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("minio complete multipart: %w", err)
	}

	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucket, key, m.urlExpire, reqParams)
	if err != nil {
		return "", fmt.Errorf("minio presigned url: %w", err)
	}

	return presignedURL.String(), nil
}

func (m *minioStore) AbortUpload(ctx context.Context, key string, uploadID string) error {
	return m.core.AbortMultipartUpload(ctx, m.bucket, key, uploadID)
}
