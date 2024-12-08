package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareQueue(channel *amqp.Channel, queueName, dlExchangeName string, expiredTime int64) (amqp.Queue, error) {
	return channel.QueueDeclare(
		queueName, // queueName
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-queue-type":           "quorum",
			"x-message-ttl":          expiredTime,
			"x-dead-letter-exchange": dlExchangeName,
		},
	)
}

func declareExchange(channel *amqp.Channel, exchangeName string) error {
	return channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		true,
		false,
		false,
		nil,
	)
}
