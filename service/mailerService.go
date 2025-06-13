package service

import (
	"log"
	"os"

	"github.com/ValeriiaHuza/weather_api/dto"
	"github.com/ValeriiaHuza/weather_api/models"
	"github.com/ValeriiaHuza/weather_api/utils"
	"gopkg.in/gomail.v2"
)

type MailerService interface {
	SendConfirmationEmail(sub models.Subscription)
	SendConfirmSuccessEmail(sub models.Subscription)
	SendWeatherUpdateEmail(sub models.Subscription, weather dto.WeatherDTO)
}

type MailerServiceImpl struct {
	Builder utils.WeatherEmailBuilder
}

func NewMailerService(builder utils.WeatherEmailBuilder) *MailerServiceImpl {
	return &MailerServiceImpl{Builder: builder}
}

func (ms *MailerServiceImpl) SendConfirmationEmail(sub models.Subscription) {
	body := ms.Builder.BuildConfirmationEmail(sub)
	ms.send(sub.Email, "Weather updates confirmation link", body)
}

func (ms *MailerServiceImpl) SendConfirmSuccessEmail(sub models.Subscription) {
	body := ms.Builder.BuildConfirmSuccessEmail(sub)
	ms.send(sub.Email, "Weather updates subscription", body)
}

func (ms *MailerServiceImpl) SendWeatherUpdateEmail(sub models.Subscription, weather dto.WeatherDTO) {
	body := ms.Builder.BuildWeatherUpdateEmail(sub, weather)
	ms.send(sub.Email, "Weather Update", body)
}

func (ms *MailerServiceImpl) send(to, subject, body string) {
	if os.Getenv("MAIL_EMAIL") == "" || os.Getenv("MAIL_PASSWORD") == "" {
		log.Fatalln("MAIL_EMAIL and MAIL_PASSWORD environment variables are required")
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
