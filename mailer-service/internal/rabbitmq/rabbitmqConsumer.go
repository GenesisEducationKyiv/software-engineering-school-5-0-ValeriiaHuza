package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	channel *amqp091.Channel
}

func NewRabbitMQConsumer(channel *amqp091.Channel) *RabbitMQConsumer {
	return &RabbitMQConsumer{channel: channel}
}

func (c *RabbitMQConsumer) Consume(queue string, handler func(body []byte)) error {
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
		return err
	}
	go func() {
		for msg := range msgs {
			handler(msg.Body)
			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
		}
	}()
	return nil
}
