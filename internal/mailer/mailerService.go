package mailer

import (
	"log"
	"os"

	"github.com/ValeriiaHuza/weather_api/internal/emailBuilder"
	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"github.com/ValeriiaHuza/weather_api/internal/service/weather"
	"gopkg.in/gomail.v2"
)

type MailService struct {
	Builder emailBuilder.WeatherEmailBuilder
}

func NewMailerService(builder emailBuilder.WeatherEmailBuilder) *MailService {
	return &MailService{Builder: builder}
}

func (ms *MailService) SendConfirmationEmail(sub subscription.Subscription) {
	body := ms.Builder.BuildConfirmationEmail(sub)
	ms.send(sub.Email, "Weather updates confirmation link", body)
}

func (ms *MailService) SendConfirmSuccessEmail(sub subscription.Subscription) {
	body := ms.Builder.BuildConfirmSuccessEmail(sub)
	ms.send(sub.Email, "Weather updates subscription", body)
}

func (ms *MailService) SendWeatherUpdateEmail(sub subscription.Subscription, weather weather.WeatherDTO) {
	body := ms.Builder.BuildWeatherUpdateEmail(sub, weather)
	ms.send(sub.Email, "Weather Update", body)
}

func (ms *MailService) send(to, subject, body string) {
	if os.Getenv("MAIL_EMAIL") == "" || os.Getenv("MAIL_PASSWORD") == "" {
		log.Println("MAIL_EMAIL and MAIL_PASSWORD environment variables are required")
		return
	}

	from := os.Getenv("MAIL_EMAIL")
	password := os.Getenv("MAIL_PASSWORD")

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
