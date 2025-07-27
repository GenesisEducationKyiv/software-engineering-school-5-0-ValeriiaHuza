package mailer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/rabbitmq"

	"github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
)

type dialer interface {
	DialAndSend(msg ...*gomail.Message) error
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

func (ms *MailService) StartEmailWorker(channel *amqp091.Channel) {
	go ms.consume(channel, rabbitmq.SendEmail, func(body []byte) {
		var job EmailJob
		if err := json.Unmarshal(body, &job); err != nil {
			log.Println("Failed to unmarshal EmailJob:", err)
			return
		}
		log.Printf("Processing EmailJob: %+v", job)

		switch job.EmailType {
		case EmailTypeCreateSubscription:
			ms.SendConfirmationEmail(job.Subscription)
		case EmailTypeConfirmSuccess:
			ms.SendConfirmSuccessEmail(job.Subscription)
		default:
			log.Println("Unknown email type:", job.EmailType)
		}
	})

	go ms.consume(channel, rabbitmq.WeatherUpdate, func(body []byte) {
		var job WeatherUpdateJob
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
		false, // auto-ack
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
		if err := msg.Ack(false); err != nil {
			log.Printf("Failed to ack message - %s: %v", msg, err)
		}
	}
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
		log.Println("Email send failed:", err)
	} else {
		log.Println("Email sent to", to)
	}
}
