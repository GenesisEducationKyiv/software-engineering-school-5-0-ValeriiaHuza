package mailer

import (
	"encoding/json"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/logger"
	"go.uber.org/zap"

	"gopkg.in/gomail.v2"
)

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
}

func NewMailerService(mailEmail string, dialer dialer,
	builder weatherEmailBuilder) *MailService {
	return &MailService{
		mailEmail: mailEmail,
		dialer:    dialer,
		builder:   builder}
}

func (ms *MailService) StartEmailWorker(consumer rabbitMQConsumer) {
	consumer.Consume(rabbitmq.SendEmail, func(body []byte) {
		var job EmailJob
		if err := json.Unmarshal(body, &job); err != nil {
			logger.GetLogger().Error("Failed to unmarshal EmailJob", zap.Error(err))
			return
		}

		logger.GetLogger().Info("Processing EmailJob", zap.Any("jobType", job.EmailType), zap.Any("jobEmail", job.To))

		switch job.EmailType {
		case EmailTypeCreateSubscription:
			ms.SendConfirmationEmail(job.Subscription)
		case EmailTypeConfirmSuccess:
			ms.SendConfirmSuccessEmail(job.Subscription)
		default:
			logger.GetLogger().Error("Unknown email type", zap.String("emailType", string(job.EmailType)))
		}
	})

	consumer.Consume(rabbitmq.WeatherUpdate, func(body []byte) {
		var job WeatherUpdateJob
		if err := json.Unmarshal(body, &job); err != nil {
			logger.GetLogger().Error("Failed to unmarshal WeatherUpdateJob", zap.Error(err))
			return
		}
		logger.GetLogger().Info("Processing WeatherUpdateJob", zap.Any("jobType", job.EmailType), zap.Any("jobEmail", job.To), zap.Any("jobWeather", job.Weather))

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
		logger.GetLogger().Error("Failed to send email", zap.String("to", to), zap.Error(err))
		return
	}
	logger.GetLogger().Info("Email sent successfully", zap.String("to", to), zap.String("subject", subject))

}
