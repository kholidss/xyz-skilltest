package broker

import (
	"fmt"
	"github.com/kholidss/xyz-skilltest/internal/consts"
	"github.com/kholidss/xyz-skilltest/internal/controller/contract"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type subscriber struct {
	conn           *amqp.Connection
	queueName      string
	queueDlxName   string
	exchangeName   string
	exchangeDlName string
	controller     contract.MessageController
}

func NewSubscriber(conn *amqp.Connection, name string, controller contract.MessageController) Subscriber {
	qName := fmt.Sprintf("queue_%s", name)
	xName := fmt.Sprintf("exchange_%s", name)

	return &subscriber{
		conn:           conn,
		queueName:      qName,
		queueDlxName:   fmt.Sprintf("%v_dlq", qName),
		exchangeName:   xName,
		exchangeDlName: fmt.Sprintf("%v_dlx", xName),
		controller:     controller,
	}
}

func (s *subscriber) Listen(topics []string) error {
	channel, err := s.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return s.listenQueue(channel, topics)
}

func (s *subscriber) listenQueue(channel *amqp.Channel, topics []string) error {
	if err := declareExchange(channel, s.exchangeName); err != nil {
		return err
	}

	if err := declareExchange(channel, s.exchangeDlName); err != nil {
		return err
	}

	queue, err := declareQueue(channel, s.queueName, s.exchangeDlName, 10000)
	if err != nil {
		return err
	}

	queueDlx, err := declareQueue(channel, s.queueDlxName, s.exchangeName, 15000)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if err = channel.QueueBind(
			queue.Name,
			topic,
			s.exchangeName,
			false,
			nil,
		); err != nil {
			return err
		}

		if err = channel.QueueBind(
			queueDlx.Name,
			topic,
			s.exchangeDlName,
			false,
			nil,
		); err != nil {
			return err
		}

		if err = channel.ExchangeBind(
			s.exchangeName,
			topic,
			"exchange_event",
			false,
			nil,
		); err != nil {
			return err
		}
	}

	msg1, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,
	)

	if err != nil {
		return err
	}

	s.consumeMessage(msg1)

	return nil
}

func (s *subscriber) consumeMessage(msgs <-chan amqp.Delivery) {
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			logger.Info(fmt.Sprintf("Processing %v", d.MessageId))

			deaths, exists := d.Headers["x-death"].([]any)
			if exists && deaths[0].(amqp.Table)["count"].(int64) > consts.MaxRetryAttemps {
				logger.Info(fmt.Sprintf("Processing Greater then threshold %v", deaths[0]))

				// Do anything with unprocessable message
				// Insert to failed_jobs table or anything

				_ = d.Ack(false)
				continue
			}

			if err := s.controller.Serve(d); err != nil {
				logger.Error(err)
				_ = d.Nack(false, true)
				continue
			}

			if err := d.Ack(false); err != nil {
				logger.Error(fmt.Sprintf("Failed ack message cause: %+v", err))
				_ = d.Nack(false, true)
			}
		}
	}()
	fmt.Printf("Waiting for message [Exchange, Queue] [%v, %v]\n", s.exchangeName, s.queueName)
	<-forever
}
