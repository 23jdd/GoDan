package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"godan/internal/config"
)

type localMultipart struct {
	basePath  string
	urlPrefix string
}

func newLocalMultipart(cfg *config.LocalStorageConfig) (*localMultipart, error) {
	absPath, err := filepath.Abs(cfg.Path)
	if err != nil {
		return nil, err
	}
	return &localMultipart{basePath: absPath, urlPrefix: cfg.URLPrefix}, nil
}

func (l *localMultipart) InitUpload(_ context.Context, key string, _ string) (string, error) {
	tmpDir := l.tmpDir(key)
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return "", fmt.Errorf("local init upload: %w", err)
	}
	return key, nil
}

func (l *localMultipart) UploadPart(_ context.Context, key string, uploadID string, partNumber int, reader io.Reader, _ int64) error {
	_ = uploadID
	partPath := filepath.Join(l.tmpDir(key), fmt.Sprintf("part_%06d", partNumber))
	f, err := os.Create(partPath)
	if err != nil {
		return fmt.Errorf("local upload part %d: %w", partNumber, err)
	}
	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("local write part %d: %w", partNumber, err)
	}
	return nil
}

func (l *localMultipart) CompleteUpload(_ context.Context, key string, uploadID string, parts []UploadPart) (string, error) {
	_ = uploadID
	tmpDir := l.tmpDir(key)

	finalPath := filepath.Join(l.basePath, key)
	finalDir := filepath.Dir(finalPath)
	if err := os.MkdirAll(finalDir, 0o755); err != nil {
		return "", err
	}

	out, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	sort.Slice(parts, func(i, j int) bool { return parts[i].PartNumber < parts[j].PartNumber })

	for _, p := range parts {
		partPath := filepath.Join(tmpDir, fmt.Sprintf("part_%06d", p.PartNumber))
		f, err := os.Open(partPath)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(out, f); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
	}

	os.RemoveAll(tmpDir)
	return fmt.Sprintf("%s/%s", l.urlPrefix, key), nil
}

func (l *localMultipart) AbortUpload(_ context.Context, key string, uploadID string) error {
	_ = uploadID
	return os.RemoveAll(l.tmpDir(key))
}

func (l *localMultipart) tmpDir(key string) string {
	return filepath.Join(l.basePath, "tmp", key)
}
