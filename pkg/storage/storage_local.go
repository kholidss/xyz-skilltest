package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"io"
	"os"
	"strings"
	"time"
)

var _ Storage = (*local)(nil)

type local struct {
	basePath string
}

func NewLocalStorage(cfg *config.Config) Storage {
	return &local{basePath: cfg.StorageConfig.LocalStorageConfig.BasePath}
}

func (l *local) Put(_ context.Context, parent, name string, contents []byte, _ bool, _ string) error {
	parent = fmt.Sprintf("%s/%s", strings.Trim(l.basePath, "/"), strings.TrimSuffix(parent, "/"))

	if err := os.MkdirAll(parent, os.ModePerm); err != nil {
		return err
	}

	fo, err := os.Create(fmt.Sprintf("%s/%s", parent, name))
	if err != nil {
		return err
	}
	defer fo.Close()

	if _, err := fo.Write(contents); err != nil {
		return err
	}

	if err := fo.Sync(); err != nil {
		return err
	}

	return nil
}

func (l *local) Delete(ctx context.Context, parent, name string) error {
	parent = fmt.Sprintf("%s/%s", strings.Trim(l.basePath, "/"), strings.TrimSuffix(parent, "/"))

	if err := os.Remove(fmt.Sprintf("%s/%s", parent, name)); err != nil && os.IsNotExist(err) {
		return nil
	}

	return nil
}

func (l *local) Get(ctx context.Context, parent, name string) ([]byte, error) {
	parent = fmt.Sprintf("%s/%s", strings.Trim(l.basePath, "/"), strings.TrimSuffix(parent, "/"))

	fo, err := os.Open(fmt.Sprintf("%s/%s", parent, name))
	if err != nil && os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	defer fo.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, fo); err != nil {
		return nil, fmt.Errorf("failed to get bytes: %w", err)
	}

	return b.Bytes(), nil
}

func (l *local) IsExist(ctx context.Context, parent, name string) (bool, error) {
	parent = fmt.Sprintf("%s/%s", strings.Trim(l.basePath, "/"), strings.TrimSuffix(parent, "/"))

	_, err := os.Open(fmt.Sprintf("%s/%s", parent, name))

	return !os.IsNotExist(err), err
}

func (l *local) GetUrl(ctx context.Context, parent, name string, expirationTime time.Duration) (string, error) {
	return strings.TrimSuffix(l.basePath, "/") + strings.TrimSuffix(parent, "/") + name, nil
}
