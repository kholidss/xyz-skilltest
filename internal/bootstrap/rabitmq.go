package bootstrap

import (
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/pkg/broker"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
)

func RegistryRabbitMQSubscriber(name string, cfg *config.Config, mController contract.MessageController) broker.Subscriber {
	conn, err := broker.ConnectRabbitMQ(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	return broker.NewSubscriber(conn, name, mController)
}

func RegistryRabbitMQPublisher(name string, cfg *config.Config) broker.Publisher {
	conn, err := broker.ConnectRabbitMQ(cfg)
	if err != nil {
		logger.Fatal("dial rabbit mq failed")
	}

	return broker.NewPublisher(cfg.AppId, conn)
}
