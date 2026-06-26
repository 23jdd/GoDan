package storage

import (
	"context"
	"fmt"
	"io"

	"godan/internal/config"
)

type Storage interface {
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
}

func New(cfg *config.StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "minio":
		return newMinIO(&cfg.MinIO)
	default:
		return newLocal(&cfg.Local)
	}
}

func GenAvatarKey(userID uint64, ext string) string {
	return fmt.Sprintf("avatars/%d%s", userID, ext)
}
