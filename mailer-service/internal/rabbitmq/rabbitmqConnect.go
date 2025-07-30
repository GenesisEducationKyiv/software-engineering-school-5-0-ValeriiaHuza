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
		conn.Close() // Clean up the connection
		return nil, err
	}
	rabbitMQ := &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}

	log.Println("Connected to RabbitMQ successfully")

	return rabbitMQ, nil
}

func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
			return err
		}
	}
	return nil
}
