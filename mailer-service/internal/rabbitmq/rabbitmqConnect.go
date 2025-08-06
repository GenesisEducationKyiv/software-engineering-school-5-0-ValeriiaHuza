package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
)

type loggerInterface interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
	logger  loggerInterface
}

func ConnectToRabbitMQ(connectionUrl string, logger loggerInterface) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(connectionUrl)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", "error", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open a RabbitMQ channel", "error", err)
		conn.Close() // Clean up the connection
		return nil, err
	}
	rabbitMQ := &RabbitMQ{
		Conn:    conn,
		Channel: ch,
		logger:  logger,
	}

	logger.Info("Connected to RabbitMQ successfully")

	return rabbitMQ, nil
}

func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			r.logger.Error("Failed to close RabbitMQ channel", "error", err)
			return err
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			r.logger.Error("Failed to close RabbitMQ connection", "error", err)
			return err
		}
	}
	return nil
}
