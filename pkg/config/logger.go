package config

import (
	"fmt"
	"github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/kholidss/xyz-skilltest/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
)

type LoggerConfig struct {
	LogLevel      string `mapstructure:"log_level"`
	LogDriver     string `mapstructure:"log_driver"`
	LokiConfig    `mapstructure:",squash"`
	GraylogConfig `mapstructure:",squash"`
}

type LokiConfig struct {
	LogLokiHost string `mapstructure:"log_loki_host"`
	LogLokiPort int    `mapstructure:"log_loki_port"`
}

type GraylogConfig struct {
	LogGraylogHost string `mapstructure:"log_graylog_host"`
	LogGraylogPort int    `mapstructure:"log_graylog_port"`
}

func (loggerConfig *LoggerConfig) WithLokiHook(cfg *Config) logrus.Hook {
	// Configure the Loki hook
	opts := lokirus.NewLokiHookOptions().
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		//WithFormatter(&logrus.JSONFormatter{}).
		WithStaticLabels(lokirus.Labels{
			"app":         cfg.AppName,
			"environment": util.EnvironmentTransform(cfg.AppEnv),
		})

	var levels []logrus.Level
	cfgLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		cfgLevel = logrus.InfoLevel
	}

	for _, level := range logrus.AllLevels {
		if level <= cfgLevel {
			levels = append(levels, level)
		}
	}

	return lokirus.NewLokiHookWithOpts(
		fmt.Sprintf("http://%v:%v", loggerConfig.LogLokiHost, loggerConfig.LogLokiPort),
		opts,
		levels...,
	)
}

func (loggerConfig *LoggerConfig) WithGraylogHook(cfg *Config) logrus.Hook {
	hook := graylog.NewGraylogHook(fmt.Sprintf("%v:%v", loggerConfig.LogGraylogHost, loggerConfig.LogGraylogPort), map[string]interface{}{
		"app":         cfg.AppName,
		"environment": util.EnvironmentTransform(cfg.AppEnv),
	})

	return hook
}
