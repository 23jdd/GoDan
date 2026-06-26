package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"godan/internal/config"
)

type local struct {
	path      string
	urlPrefix string
}

func newLocal(cfg *config.LocalStorageConfig) (*local, error) {
	absPath, err := filepath.Abs(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("local storage path error: %w", err)
	}
	if err := os.MkdirAll(absPath, 0o755); err != nil {
		return nil, fmt.Errorf("local storage mkdir error: %w", err)
	}
	return &local{path: absPath, urlPrefix: cfg.URLPrefix}, nil
}

func (l *local) Upload(_ context.Context, key string, reader io.Reader, _ int64, _ string) (string, error) {
	fullPath := filepath.Join(l.path, key)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("local upload mkdir: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("local upload create: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return "", fmt.Errorf("local upload write: %w", err)
	}

	return fmt.Sprintf("%s/%s", l.urlPrefix, key), nil
}

func (l *local) Delete(_ context.Context, key string) error {
	return os.Remove(filepath.Join(l.path, key))
}

func (l *local) Path() string {
	return l.path
}
