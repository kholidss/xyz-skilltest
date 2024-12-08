package config

type CDNConfig struct {
	Provider string `mapstructure:"cdn_provider"`

	Cloudinary `mapstructure:",squash"`
}

type Cloudinary struct {
	CloudName     string `mapstructure:"cloudinary_cloud_name"`
	APIKey        string `mapstructure:"cloudinary_api_key"`
	APISecret     string `mapstructure:"cloudinary_api_secret"`
	Dir           string `mapstructure:"cloudinary_dir"`
	StreamTimeout int    `mapstructure:"cloudinary_stream_timeout"`
}
