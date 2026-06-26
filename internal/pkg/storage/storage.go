package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"godan/internal/config"
)

type Storage interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
}

type MultipartUploader interface {
	InitUpload(ctx context.Context, key string, contentType string) (string, error)
	UploadPart(ctx context.Context, key string, uploadID string, partNumber int, reader io.Reader, size int64) error
	CompleteUpload(ctx context.Context, key string, uploadID string, parts []UploadPart) (string, error)
	AbortUpload(ctx context.Context, key string, uploadID string) error
}

type UploadPart struct {
	PartNumber int `json:"part_number"`
	PartSize   int `json:"part_size"`
}

func New(cfg *config.StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "minio":
		return newMinIO(&cfg.MinIO)
	default:
		return newLocal(&cfg.Local)
	}
}

func NewMultipartUploader(cfg *config.StorageConfig) (MultipartUploader, error) {
	switch cfg.Type {
	case "minio":
		return newMinIO(&cfg.MinIO)
	default:
		return newLocalMultipart(&cfg.Local)
	}
}

func GenVideoKey(userID uint64, filename string) string {
	return fmt.Sprintf("videos/%d/%s", userID, filename)
}

func GenAvatarKey(userID uint64, ext string) string {
	return fmt.Sprintf("avatars/%d%s", userID, ext)
}

func GenCoverKey(userID uint64, ext string) string {
	return fmt.Sprintf("covers/%d_%d%s", userID, time.Now().UnixNano(), ext)
}
