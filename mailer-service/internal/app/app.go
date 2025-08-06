package app

import (
	"fmt"
	"log"
	"strconv"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/emailBuilder"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/mailer"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/logger"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type loggerInterface interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Sync() error
}

func Run() error {

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}

	defer logger.Sync()

	logger.Info("Starting Mailer Service...")

	config, err := config.LoadEnvVariables()

	if err != nil {
		return err
	}

	rabbit, err := rabbitmq.ConnectToRabbitMQ(config.RabbitMQUrl, logger)
	if err != nil {
		return err
	}
	defer rabbit.Conn.Close()
	defer rabbit.Channel.Close()

	if err := declareQueues(rabbit); err != nil {
		return err
	}

	initServices(*config, *rabbit, logger)

	router := gin.Default()

	port := strconv.Itoa(config.MailerPort)
	return router.Run(":" + port)
}

func initServices(config config.Config, rabbit rabbitmq.RabbitMQ, logger loggerInterface) {
	emailBuilder := emailBuilder.NewWeatherEmailBuilder(config.ApiURL, logger)

	mailEmail := config.MailEmail
	dialer := gomail.NewDialer(config.MailDialerHost, config.MailDialerPort, mailEmail, config.MailPassword)
	mailerService := mailer.NewMailerService(mailEmail, dialer, emailBuilder, logger)

	rabbitmqConsumer := rabbitmq.NewRabbitMQConsumer(rabbit.Channel, logger)

	go mailerService.StartEmailWorker(rabbitmqConsumer)

}

func declareQueues(r *rabbitmq.RabbitMQ) error {
	queues := []string{
		rabbitmq.SendEmail,
		rabbitmq.WeatherUpdate,
	}

	for _, q := range queues {
		_, err := r.Channel.QueueDeclare(q, true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", q, err)
		}
	}
	return nil
}
