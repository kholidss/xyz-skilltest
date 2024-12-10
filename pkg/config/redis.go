package config

// RedisConfig holds the RedisConfig configuration
type RedisConfig struct {
	Host               string `mapstructure:"redis_host"`
	Port               int    `mapstructure:"redis_port"`
	DB                 int    `mapstructure:"redis_db"`
	Username           string `mapstructure:"redis_username"`
	Password           string `mapstructure:"redis_password"`
	PoolSize           int    `mapstructure:"redis_pool_size"`
	ReadTimeout        int    `mapstructure:"redis_read_timeout"`
	WriteTimeout       int    `mapstructure:"redis_write_timeout"`
	PoolTimeout        int    `mapstructure:"redis_pool_timeout"`
	MinIdleConn        int    `mapstructure:"redis_min_idle_conn"`
	IdleTimeout        int    `mapstructure:"redis_idle_timeout"`
	TLS                bool   `mapstructure:"redis_tls"`
	InsecureSkipVerify bool   `mapstructure:"redis_insecure_skip_verify"`
}
