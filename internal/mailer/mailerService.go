package mailer

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"gopkg.in/gomail.v2"
)

type weatherEmailBuilder interface {
	BuildWeatherUpdateEmail(sub subscription.Subscription, weather client.WeatherDTO) string
	BuildConfirmationEmail(sub subscription.Subscription) string
	BuildConfirmSuccessEmail(sub subscription.Subscription) string
}

type MailService struct {
	mailEmal string
	dialer   *gomail.Dialer
	builder  weatherEmailBuilder
}

func NewMailerService(mailEmal string, dialer *gomail.Dialer,
	builder weatherEmailBuilder) *MailService {
	return &MailService{
		mailEmal: mailEmal,
		dialer:   dialer,
		builder:  builder}
}

func (ms *MailService) SendConfirmationEmail(sub subscription.Subscription) {
	body := ms.builder.BuildConfirmationEmail(sub)
	ms.send(sub.Email, "Weather updates confirmation link", body)
}

func (ms *MailService) SendConfirmSuccessEmail(sub subscription.Subscription) {
	body := ms.builder.BuildConfirmSuccessEmail(sub)
	ms.send(sub.Email, "Weather updates subscription", body)
}

func (ms *MailService) SendWeatherUpdateEmail(sub subscription.Subscription, weather client.WeatherDTO) {
	body := ms.builder.BuildWeatherUpdateEmail(sub, weather)
	ms.send(sub.Email, "Weather Update", body)
}

func (ms *MailService) send(to, subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", ms.mailEmal)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := ms.dialer.DialAndSend(m); err != nil {
		log.Println("Email send failed:", err)
	} else {
		log.Println("Email sent to", to)
	}
}
