package broker

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type publisher struct {
	appid        string
	connection   *amqp.Connection
	exchangeName string
}

func NewPublisher(appid string, conn *amqp.Connection) Publisher {
	p := &publisher{
		appid:        appid,
		connection:   conn,
		exchangeName: "exchange_event", // when publishing we always pub to this exchange
	}

	if err := p.setup(); err != nil {
		logger.Fatal("publisher setup not initialized")
	}

	return p
}

func (p *publisher) setup() error {
	ch, err := p.connection.Channel()
	if err != nil {
		return err
	}

	return declareExchange(ch, p.exchangeName)
}

func (p *publisher) Publish(route string, payload MessagePayload) error {
	channel, err := p.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	dataPublish, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = channel.PublishWithContext(
		context.Background(),
		p.exchangeName,
		route,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         dataPublish,
			DeliveryMode: 2,
			Timestamp:    time.Now(),
			AppId:        p.appid,
			MessageId:    uuid.NewString(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
