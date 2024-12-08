package http

import (
	"fmt"
	"github.com/kholidss/xyz-skilltest/internal/bootstrap"
	"github.com/kholidss/xyz-skilltest/internal/controller"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"strings"
)

func StartNSListener(cfg *config.Config) {
	name := cfg.NoSleepConfig.FlagName
	topics := strings.Split(cfg.NoSleepConfig.Topics, "|")
	mController := controller.NewLogController()

	subs := bootstrap.RegistryRabbitMQSubscriber(name, cfg, mController)

	err := subs.Listen(topics)

	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to Listen Topic Cause: %v", err))
	}
}
