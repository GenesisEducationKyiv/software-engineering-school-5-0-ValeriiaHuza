package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func ConnectToRabbitMQ(connectionUrl string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(connectionUrl)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a RabbitMQ channel: %v", err)
		return nil, err
	}

	rabbitMQ := &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}

	return rabbitMQ, nil
}
