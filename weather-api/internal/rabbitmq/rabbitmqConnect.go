package rabbitmq

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func ConnectToRabbitMQ(connectionUrl string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(connectionUrl)
	if err != nil {
		logger.GetLogger().Error("Failed to connect to RabbitMQ", zap.Error(err))
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.GetLogger().Error("Failed to open a RabbitMQ channel", zap.Error(err))
		conn.Close() // Clean up the connection
		return nil, err
	}

	rabbitMQ := &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}

	logger.GetLogger().Info("Connected to RabbitMQ", zap.String("url", connectionUrl))

	return rabbitMQ, nil
}

func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			logger.GetLogger().Error("Failed to close RabbitMQ channel", zap.Error(err))
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			logger.GetLogger().Error("Failed to close RabbitMQ connection", zap.Error(err))
			return err
		}
	}
	return nil
}
