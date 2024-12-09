package config

// AppConfig holds the AppConfig configuration
type AppConfig struct {
	AppId    string `mapstructure:"app_id"`
	AppName  string `mapstructure:"app_name"`
	AppEnv   string `mapstructure:"app_env"`
	AppHost  string `mapstructure:"app_host"`
	AppPort  int    `mapstructure:"app_port"`
	AppDebug bool   `mapstructure:"app_debug"`

	AppOtelTrace    bool   `mapstructure:"app_otel_trace"`
	AppOtelExporter string `mapstructure:"app_otel_exporter"`

	AppPasswordSecret         string `mapstructure:"app_password_secret"`
	APPPrivateKey             string `mapstructure:"app_private_key"`
	APPPublicKey              string `mapstructure:"app_public_key"`
	APPTokenUserExpiredInHour int    `mapstructure:"app_token_user_expired_in_hour"`

	DistributeTraceConfig `mapstructure:",squash"`
}
