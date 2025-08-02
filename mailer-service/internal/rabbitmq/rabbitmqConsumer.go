package rabbitmq

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/logger"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQConsumer struct {
	channel *amqp091.Channel
}

func NewRabbitMQConsumer(channel *amqp091.Channel) *RabbitMQConsumer {
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
		logger.GetLogger().Error("Failed to register a consumer", zap.Error(err))
		return
	}

	go func() {
		for msg := range msgs {
			handler(msg.Body)
			if err := msg.Ack(false); err != nil {
				logger.GetLogger().Error("Failed to ack message", zap.Error(err))
			}
		}
	}()
}
