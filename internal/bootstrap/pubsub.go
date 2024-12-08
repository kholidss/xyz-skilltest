package bootstrap

import (
	"context"
	"fmt"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/pubsubx"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func RegistryPubSubConsumer(cfg *config.Config) pubsubx.Subscriberer {
	credOpt := option.WithCredentialsFile(cfg.GCPConfig.PubsubAccountPath)
	cl, err := pubsub.NewClient(context.Background(), cfg.GCPConfig.ProjectID, credOpt)
	if err != nil {
		logger.Fatal(fmt.Sprintf("google pusbsub conusmer error:%v", err))
	}

	return pubsubx.NewGSubscriber(cl)
}

func RegistryPubSubPublisher(cfg *config.Config) pubsubx.Publisher {
	credOpt := option.WithCredentialsFile(cfg.GCPConfig.PubsubAccountPath)
	cl, err := pubsub.NewClient(context.Background(), cfg.GCPConfig.ProjectID, credOpt)
	if err != nil {
		logger.Fatal(fmt.Sprintf("google pusbsub publisher error:%v", err))
	}

	return pubsubx.NewGPublisher(cl)
}
