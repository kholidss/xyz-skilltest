package storage

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Compile-time check to verify implements interface.
var _ Storage = (*gcs)(nil)

// gcs implements the Blob interface and provides the ability
// write files to Google Cloud Storage.
type gcs struct {
	client     *storage.Client
	JWTConfig  *jwt.Config
	bucketName string
}

// NewGCS creates a Google Cloud Storage Client
func NewGCS(ctx context.Context, bucketName, cfgJsonFIle string) (Storage, error) {
	credOpt := option.WithCredentialsFile(cfgJsonFIle)

	jsonKey, err := os.ReadFile(cfgJsonFIle)
	if err != nil {
		return nil, fmt.Errorf("storage.read File: %w", err)
	}

	jwtConfig, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		fmt.Printf("Error creating JWT config: %v\n", err)
		return nil, fmt.Errorf("storage.getting condig from json: %w", err)
	}

	client, err := storage.NewClient(ctx, credOpt)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}

	return &gcs{client, jwtConfig, bucketName}, nil
}

// Put creates a new cloud storage object or overwrites an existing one.
func (s *gcs) Put(ctx context.Context, parent, name string, contents []byte, cacheable bool, contentType string) error {
	cacheControl := "public, max-age=86400"
	if !cacheable {
		cacheControl = "no-cache, max-age=0"
	}

	wc := s.client.Bucket(s.bucketName).Object(fmt.Sprintf("%s/%s", strings.TrimSuffix(parent, "/"), name)).NewWriter(ctx)
	wc.CacheControl = cacheControl
	if contentType != "" {
		wc.ContentType = contentType
	}

	if _, err := wc.Write(contents); err != nil {
		return fmt.Errorf("storage.Writer.Write: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("storage.Writer.Close: %w", err)
	}

	return nil
}

// Delete deletes a cloud storage object, returns nil if the object was
// successfully deleted, or of the object doesn't exist.
func (s *gcs) Delete(ctx context.Context, parent, name string) error {
	if err := s.client.Bucket(s.bucketName).Object(fmt.Sprintf("%s/%s", strings.TrimSuffix(parent, "/"), name)).Delete(ctx); err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			// Object doesn't exist; presumably already deleted.
			return nil
		}
		return fmt.Errorf("storage.DeleteObject: %w", err)
	}
	return nil
}

// Get returns the contents for the given object. If the object does not
// exist, it returns ErrNotFound.
func (s *gcs) Get(ctx context.Context, parent, name string) ([]byte, error) {
	r, err := s.client.Bucket(s.bucketName).Object(fmt.Sprintf("%s/%s", strings.TrimSuffix(parent, "/"), name)).NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return nil, ErrNotFound
		}
	}
	defer r.Close()

	var b bytes.Buffer
	if _, err := io.Copy(&b, r); err != nil {
		return nil, fmt.Errorf("failed to download bytes: %w", err)
	}

	return b.Bytes(), nil
}

func (s *gcs) IsExist(ctx context.Context, parent, name string) (bool, error) {
	r, err := s.client.Bucket(s.bucketName).Object(fmt.Sprintf("%s/%s", strings.TrimSuffix(parent, "/"), name)).NewReader(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			return false, ErrNotFound
		}
		return false, err
	}
	defer r.Close()
	return true, nil
}

func (s *gcs) GetUrl(ctx context.Context, parent, name string, expirationTime time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         http.MethodGet,
		GoogleAccessID: s.JWTConfig.Email,
		PrivateKey:     s.JWTConfig.PrivateKey,
		Expires:        time.Now().Add(expirationTime),
	}

	signedURL, err := storage.SignedURL(s.bucketName, fmt.Sprintf("%s/%s", strings.TrimSuffix(parent, "/"), name), opts)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
