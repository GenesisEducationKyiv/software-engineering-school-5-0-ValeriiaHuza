package rabbitmq

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/logger"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	channel *amqp091.Channel
	logger  logger.Logger
}

func NewRabbitMQConsumer(channel *amqp091.Channel, logger logger.Logger) *RabbitMQConsumer {
	return &RabbitMQConsumer{channel: channel}
}

func (c *RabbitMQConsumer) Consume(queue string, handler func(body []byte)) {
	msgs, err := c.channel.Consume(
		queue,
		"",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		c.logger.Error("Failed to register a consumer", "error", err)
		return
	}

	go func() {
		for msg := range msgs {
			handler(msg.Body)
			if err := msg.Ack(false); err != nil {
				c.logger.Error("Failed to ack message", "error", err)
			}
		}
	}()
}
