package cdn

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/util"
	"time"

	cl "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinary struct {
	cfg    *config.Config
	client *cl.Cloudinary
}

// Compile-time check to verify implements interface.
var _ CDN = (*cloudinary)(nil)

// NewCloudinaryCDN creates a new instance cloudinary cdn
func NewCloudinaryCDN(cfg *config.Config) (CDN, error) {
	client, err := cl.NewFromParams(
		cfg.CDNConfig.Cloudinary.CloudName,
		cfg.CDNConfig.Cloudinary.APIKey,
		cfg.CDNConfig.Cloudinary.APISecret,
	)
	if err != nil {
		return nil, err
	}

	return &cloudinary{cfg, client}, nil
}

// Put creates a new cloud storage object or overwrites with input byte data an existing one to SDK cloudinary.
func (c *cloudinary) Put(ctx context.Context, name string, contents []byte) (any, error) {
	timeout := time.Duration(util.SetDefaultInt(c.cfg.CDNConfig.Cloudinary.StreamTimeout, 30)) * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	useFilename := true
	obj, err := c.client.Upload.Upload(ctx, bytes.NewReader(contents), uploader.UploadParams{
		Folder:           c.cfg.CDNConfig.Cloudinary.Dir,
		UseFilename:      &useFilename,
		FilenameOverride: name,
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary.Put: %w", err)
	}
	return obj, err
}

// Delete removes an object from the SDK cloudinary.
func (c *cloudinary) Delete(ctx context.Context, publicID string) error {
	_, err := c.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("cloudinary.Delete: %w", err)
	}
	return nil
}

func (c *cloudinary) Get(ctx context.Context, identifier string) ([]byte, error) {
	panic("implement me")
}
