package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	Channel *amqp091.Channel
}

func NewRabbitMQPublisher(channel *amqp091.Channel) *RabbitMQPublisher {
	return &RabbitMQPublisher{Channel: channel}
}

func (p *RabbitMQPublisher) Publish(queue string, payload any) error {

	if p.Channel == nil {
		return fmt.Errorf("channel is nil")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message for queue %s: %w", queue, err)
	}

	err = p.Channel.Publish(
		"",    // exchange
		queue, // routing key (queue name)
		false, // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish to queue %s: %w", queue, err)
	}

	return nil
}
