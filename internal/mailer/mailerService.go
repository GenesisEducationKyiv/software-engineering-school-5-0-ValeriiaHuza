package mailer

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/config"
	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/ValeriiaHuza/weather_api/internal/emailBuilder"
	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"gopkg.in/gomail.v2"
)

type MailService struct {
	builder emailBuilder.WeatherEmailBuilder
}

func NewMailerService(builder emailBuilder.WeatherEmailBuilder) *MailService {
	return &MailService{builder: builder}
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
	from := config.AppConfig.MailEmail
	password := config.AppConfig.MailPassword

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Email send failed:", err)
	} else {
		log.Println("Email sent to", to)
	}
}
