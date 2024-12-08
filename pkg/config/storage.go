package config

type StorageConfig struct {
	MinioConfig        `mapstructure:",squash"`
	GCSConfig          `mapstructure:",squash"`
	AWSConfig          `mapstructure:",squash"`
	LocalStorageConfig `mapstructure:",squash"`
}

type LocalStorageConfig struct {
	BasePath string `mapstructure:"local_storage_basepath"`
}

type MinioConfig struct {
	BaseUrl         string `mapstructure:"minio_base_url"`
	AccessKeyId     string `mapstructure:"minio_access_key_id"`
	SecretAccessKey string `mapstructure:"minio_secret_access_key"`
	Location        string `mapstructure:"minio_location"`
}

type GCSConfig struct {
	ServiceAccountPath string `mapstructure:"gcs_service_account_path"`
	BucketName         string `mapstructure:"gcs_bucket_name"`
}

type AWSConfig struct {
	Region       string `mapstructure:"aws_region"`
	AccessKey    string `mapstructure:"aws_access_key"`
	AccessSecret string `mapstructure:"aws_access_secret"`
	BucketName   string `mapstructure:"aws_bucket_name"`
}
