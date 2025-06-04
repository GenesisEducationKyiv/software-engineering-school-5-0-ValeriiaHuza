package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/ValeriiaHuza/weather_api/models"
	"gopkg.in/gomail.v2"
)

func SendConfirmationEmail(subscription models.Subscription) {

	confirmationLink := BuildURL("/api/confirm/") + subscription.Token

	body := fmt.Sprintf(`
		<p>Hello from Weather Updates!</p>
		<p>You subscribed for <strong>%s</strong> updates for <strong>%s</strong> weather.</p>
		<p>Please confirm your subscription by clicking the link below:</p>
		<p><a href="%s">Your link</a></p>`,
		string(subscription.Frequency), subscription.City, confirmationLink)

	SendEmail(subscription.Email, "Weather updates confirmation link", body)

}

func SendConfirmSuccessMail(subscription models.Subscription) {

	unsubscribeLink := BuildURL("/api/unsubscribe/") + subscription.Token

	body := fmt.Sprintf(`
		<p>Hello from Weather Updates!</p>
		<p>You have successfully confirmed your subscription!</p>
		<p>If you want to unsubscribe, click the link below:</p>
		<p><a href="%s">Your link</a></p>`,
		unsubscribeLink)

	SendEmail(subscription.Email,
		"Weather updates subscription",
		body)

}

func SendEmail(to, subject, body string) {
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

func BuildURL(path string) string {
	host := os.Getenv("APP_URL")
	return host + path
}
