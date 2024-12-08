package config

type NoSleepConfig struct {
	FlagName string `json:"flag_name" mapstructure:"flag_name"`
	Topics   string `json:"topics" mapstructure:"topics"`
}
