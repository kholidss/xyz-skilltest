package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	minio2 "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"time"
)

var _ Storage = (*minio)(nil)

type minio struct {
	client *minio2.Client
}

func NewMinio(ctx context.Context, cfg *config.Config) (Storage, error) {
	client, err := minio2.New(cfg.MinioConfig.BaseUrl, &minio2.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioConfig.AccessKeyId, cfg.MinioConfig.SecretAccessKey, ""),
		Secure: false,
		Region: cfg.MinioConfig.Location,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &minio{client: client}, nil
}

func (m *minio) Put(ctx context.Context, parent, name string, contents []byte, cacheAble bool, contentType string) error {
	data := bytes.NewReader(contents)
	options := minio2.PutObjectOptions{
		CacheControl: "max-age=600",
	}
	if contentType != "" {
		options.ContentType = contentType
	}
	_, err := m.client.PutObject(ctx, parent, name, data, data.Size(), options)
	if err != nil {
		return err
	}
	return nil
}

func (m *minio) Delete(ctx context.Context, parent, name string) error {
	options := minio2.RemoveObjectOptions{}
	err := m.client.RemoveObject(ctx, parent, name, options)
	if err != nil {
		return err
	}
	return nil
}

func (m *minio) Get(ctx context.Context, parent, name string) ([]byte, error) {
	options := minio2.GetObjectOptions{}
	data, err := m.client.GetObject(ctx, parent, name, options)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	if err != nil {
		return nil, err
	}
	var file bytes.Buffer
	_, err = io.Copy(&file, data)
	if err != nil {
		return nil, err
	}
	return file.Bytes(), nil
}

func (m *minio) IsExist(ctx context.Context, parent, name string) (bool, error) {
	options := minio2.GetObjectOptions{}
	data, err := m.client.GetObject(ctx, parent, name, options)
	if err != nil {
		return false, err
	}
	defer data.Close()
	return true, nil
}

func (m *minio) GetUrl(ctx context.Context, parent, name string, expirationTime time.Duration) (string, error) {
	preSignedURL, err := m.client.PresignedGetObject(ctx, parent, name, expirationTime, nil)
	if err != nil {
		return "", err
	}

	return preSignedURL.String(), nil
}
