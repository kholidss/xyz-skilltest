package cdn

import (
	"context"
)

type CDN interface {
	// Put uploads an object to the storage with the specified filename.
	Put(ctx context.Context, name string, contents []byte) (any, error)

	// Delete removes an object from the storage.
	Delete(ctx context.Context, identifier string) error

	// Get retrieves the contents of an object from the storage.
	Get(ctx context.Context, identifier string) ([]byte, error)
}
