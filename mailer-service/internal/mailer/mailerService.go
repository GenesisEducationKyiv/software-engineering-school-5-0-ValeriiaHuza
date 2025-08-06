package mailer

import (
	"encoding/json"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/rabbitmq"

	"gopkg.in/gomail.v2"
)

type loggerInterface interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type dialer interface {
	DialAndSend(msg ...*gomail.Message) error
}

type rabbitMQConsumer interface {
	Consume(queue string, handler func(body []byte))
}

type weatherEmailBuilder interface {
	BuildWeatherUpdateEmail(sub SubscriptionDTO, weather WeatherDTO, time time.Time) string
	BuildConfirmationEmail(sub SubscriptionDTO) string
	BuildConfirmSuccessEmail(sub SubscriptionDTO) string
}

type MailService struct {
	mailEmail string
	dialer    dialer
	builder   weatherEmailBuilder
	logger    loggerInterface
}

func NewMailerService(mailEmail string, dialer dialer,
	builder weatherEmailBuilder, logger loggerInterface) *MailService {
	return &MailService{
		mailEmail: mailEmail,
		dialer:    dialer,
		builder:   builder,
		logger:    logger,
	}
}

func (ms *MailService) StartEmailWorker(consumer rabbitMQConsumer) {
	consumer.Consume(rabbitmq.SendEmail, func(body []byte) {
		var job EmailJob
		if err := json.Unmarshal(body, &job); err != nil {
			ms.logger.Error("Failed to unmarshal EmailJob", "error", err)
			return
		}

		ms.logger.Info("Processing EmailJob", "jobType", job.EmailType, "jobEmail", job.To)

		switch job.EmailType {
		case EmailTypeCreateSubscription:
			ms.SendConfirmationEmail(job.Subscription)
		case EmailTypeConfirmSuccess:
			ms.SendConfirmSuccessEmail(job.Subscription)
		default:
			ms.logger.Error("Unknown email type", "emailType", job.EmailType)
		}
	})

	consumer.Consume(rabbitmq.WeatherUpdate, func(body []byte) {
		var job WeatherUpdateJob
		if err := json.Unmarshal(body, &job); err != nil {
			ms.logger.Error("Failed to unmarshal WeatherUpdateJob", "error", err)
			return
		}
		ms.logger.Info("Processing WeatherUpdateJob", "jobType", job.EmailType, "jobEmail", job.To, "jobWeather", job.Weather)

		ms.SendWeatherUpdateEmail(job.Subscription, job.Weather)
	})
}

func (ms *MailService) SendConfirmationEmail(sub SubscriptionDTO) {
	body := ms.builder.BuildConfirmationEmail(sub)
	ms.send(sub.Email, "Weather updates confirmation link", body)
}

func (ms *MailService) SendConfirmSuccessEmail(sub SubscriptionDTO) {
	body := ms.builder.BuildConfirmSuccessEmail(sub)
	ms.send(sub.Email, "Weather updates subscription", body)
}

func (ms *MailService) SendWeatherUpdateEmail(sub SubscriptionDTO, weather WeatherDTO) {
	body := ms.builder.BuildWeatherUpdateEmail(sub, weather, time.Now())
	ms.send(sub.Email, "Weather Update", body)
}
func (ms *MailService) send(to, subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", ms.mailEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := ms.dialer.DialAndSend(m); err != nil {
		ms.logger.Error("Failed to send email", "to", to, "error", err)
		return
	}
	ms.logger.Info("Email sent successfully", "to", to, "subject", subject)

}
