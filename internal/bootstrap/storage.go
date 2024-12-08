package bootstrap

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/storage"
)

func RegistryLocalStorage(cfg *config.Config) storage.Storage {
	return storage.NewLocalStorage(cfg)
}

func RegistryGCS(cfg *config.Config) storage.Storage {
	sto, err := storage.NewGCS(context.Background(), cfg.StorageConfig.GCSConfig.BucketName, cfg.StorageConfig.GCSConfig.ServiceAccountPath)
	if err != nil {
		logger.Fatal(fmt.Sprintf("load gcs error %v", err), logger.Any("account_path", cfg.StorageConfig.GCSConfig.ServiceAccountPath))
	}

	return sto
}

func RegistryMinio(cfg *config.Config) storage.Storage {
	minio, err := storage.NewMinio(context.Background(), cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("load minio error %v", err))
	}
	return minio
}

// RegistryAWSSession initialize aws session
func RegistryAWSSession(appCtx *config.Config) storage.Storage {
	var awsConfig aws.Config

	fmt.Println(appCtx.AWSConfig)

	awsConfig.CredentialsChainVerboseErrors = aws.Bool(true)
	region := appCtx.AWSConfig.Region
	if len(region) != 0 {
		awsConfig.Region = aws.String(region)
	}

	accessKeyID := appCtx.AWSConfig.AccessKey
	secretAccessKey := appCtx.AWSConfig.AccessSecret
	if len(accessKeyID) != 0 && len(secretAccessKey) != 0 {
		awsConfig.Credentials = credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	}

	sess, err := session.NewSession(&awsConfig)

	if err != nil {
		logger.Fatal(err, logger.EventName("aws-session"))
	}

	return storage.NewAwsS3(sess)
}
