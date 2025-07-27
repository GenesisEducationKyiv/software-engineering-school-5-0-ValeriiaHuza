package app

import (
	"fmt"
	"strconv"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/emailBuilder"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/mailer"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/logger"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

func Run() error {

	if err := logger.InitLoggerFile("app.log"); err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	defer logger.CloseLogFile()

	config, err := config.LoadEnvVariables()

	if err != nil {
		return err
	}

	rabbit, err := rabbitmq.ConnectToRabbitMQ(config.RabbitMQUrl)
	if err != nil {
		return err
	}
	defer rabbit.Conn.Close()
	defer rabbit.Channel.Close()

	if err := declareQueues(rabbit); err != nil {
		return err
	}

	emailPublisher := rabbitmq.NewRabbitMQPublisher(rabbit.Channel)

	initServices(*config, emailPublisher)

	router := gin.Default()

	port := strconv.Itoa(config.MailerPort)
	return router.Run(":" + port)
}

func initServices(config config.Config, emailPublisher *rabbitmq.RabbitMQPublisher) {
	emailBuilder := emailBuilder.NewWeatherEmailBuilder(config.ApiURL)

	mailEmail := config.MailEmail
	dialer := gomail.NewDialer("smtp.gmail.com", 587, mailEmail, config.MailPassword)
	mailerService := mailer.NewMailerService(mailEmail, dialer, emailBuilder)

	go mailerService.StartEmailWorker(emailPublisher.Channel)

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
