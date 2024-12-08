package controller

import (
	"context"
	"fmt"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/pubsubx"
)

type pubsubController struct {
}

func (p *pubsubController) Serve(ctx context.Context, message *pubsubx.Message) {
	logger.InfoWithContext(ctx, fmt.Sprintf("Received Data %s", string(message.Data)))
}

func NewPubsubController() contract.PubSubMessageController {
	return &pubsubController{}
}
