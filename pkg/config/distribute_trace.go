package config

type DistributeTraceConfig struct {
	JaegerConfig `mapstructure:",squash"`
	TempoConfig  `mapstructure:",squash"`
}

type TempoConfig struct {
	TempoHost string `mapstructure:"tempo_host"`
	TempoPort string `mapstructure:"tempo_port"`
}

type JaegerConfig struct {
	JaegerHost string `mapstructure:"jaeger_host"`
	JaegerPort int    `mapstructure:"jaeger_port"`
}
