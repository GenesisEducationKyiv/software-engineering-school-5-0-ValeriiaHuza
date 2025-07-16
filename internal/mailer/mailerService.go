package mailer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/subscription"
	"github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
)

type weatherEmailBuilder interface {
	BuildWeatherUpdateEmail(sub subscription.Subscription, weather client.WeatherDTO, time time.Time) string
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

func (ms *MailService) StartEmailWorker(channel *amqp091.Channel) {
	go ms.consume(channel, rabbitmq.SendEmail, func(body []byte) {
		var job subscription.EmailJob
		if err := json.Unmarshal(body, &job); err != nil {
			log.Println("Failed to unmarshal EmailJob:", err)
			return
		}
		log.Printf("Processing EmailJob: %+v", job)

		switch job.EmailType {
		case subscription.EmailTypeCreateSubscription:
			ms.SendConfirmationEmail(job.Subscription)
		case subscription.EmailTypeConfirmSuccess:
			ms.SendConfirmSuccessEmail(job.Subscription)
		default:
			log.Println("Unknown email type:", job.EmailType)
		}
	})

	go ms.consume(channel, rabbitmq.WeatherUpdate, func(body []byte) {
		var job subscription.WeatherUpdateJob
		if err := json.Unmarshal(body, &job); err != nil {
			log.Println("Failed to unmarshal WeatherUpdateJob:", err)
			return
		}
		log.Printf("Processing WeatherUpdateJob: %+v", job)

		ms.SendWeatherUpdateEmail(job.Subscription, job.Weather)
	})
}

func (ms *MailService) consume(channel *amqp091.Channel, queue string, handler func(body []byte)) {
	msgs, err := channel.Consume(
		queue,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer for queue %s: %v", queue, err)
	}

	for msg := range msgs {
		handler(msg.Body)
	}
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
	body := ms.builder.BuildWeatherUpdateEmail(sub, weather, time.Now())
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
