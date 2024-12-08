package bootstrap

import (
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/util"
)

func RegistryLogger(cfg *config.Config) {
	loggerConfig := logger.Config{
		Environment: util.EnvironmentTransform(cfg.AppEnv),
		Debug:       cfg.AppDebug,
		Level:       cfg.LogLevel,
		ServiceName: cfg.AppName,
	}

	logger.Setup(loggerConfig)

	switch cfg.LogDriver {
	case logger.LogDriverLoki:
		logger.AddHook(cfg.LoggerConfig.WithLokiHook(cfg))
		break
	case logger.LogDriverGraylog:
		logger.AddHook(cfg.LoggerConfig.WithGraylogHook(cfg))
		break
	default:
		break
	}
}
