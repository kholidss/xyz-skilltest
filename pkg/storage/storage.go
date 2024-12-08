package storage

import (
	"context"
	"fmt"
	"time"
)

var ErrNotFound = fmt.Errorf("storage object not found")

type Storage interface {
	Put(ctx context.Context, parent, name string, contents []byte, cacheAble bool, contentType string) error

	// Delete deletes an object or does nothing if the object doesn't exist.
	Delete(ctx context.Context, parent, name string) error

	// Get fetches the object's contents.
	Get(ctx context.Context, parent, name string) ([]byte, error)

	IsExist(ctx context.Context, parent, name string) (bool, error)
	GetUrl(ctx context.Context, parent, name string, expirationTime time.Duration) (string, error)
}
