package subscription

import (
	"errors"
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/rabbitmq"
	"github.com/google/uuid"
)

type mailPublisher interface {
	Publish(queue string, payload any) error
}

type subscriptionRepository interface {
	Create(sub Subscription) error
	Update(sub Subscription) error
	FindByToken(token string) (*Subscription, error)
	Delete(sub Subscription) error
	FindByEmail(email string) (*Subscription, error)
	FindByFrequencyAndConfirmation(freq Frequency) ([]Subscription, error)
}

type weatherService interface {
	GetWeather(city string) (*client.WeatherDTO, error)
}

type SubscribeService struct {
	weatherService         weatherService
	subscriptionRepository subscriptionRepository
	mailPublisher          mailPublisher
}

func NewSubscribeService(weatherService weatherService,
	repository subscriptionRepository,
	mailPublisher mailPublisher) *SubscribeService {
	return &SubscribeService{
		weatherService:         weatherService,
		subscriptionRepository: repository,
		mailPublisher:          mailPublisher,
	}
}

func (ss *SubscribeService) SubscribeForWeatherUpdates(email string,
	city string, frequency Frequency) error {

	if _, err := ss.weatherService.GetWeather(city); err != nil {
		return err
	}

	log.Printf("Subscribing %s for %s weather updates in %s", email, string(frequency), city)

	subscribed := ss.emailSubscribed(email)
	if subscribed {
		return ErrEmailAlreadySubscribed
	}

	token := ss.generateToken()

	newSubscription := Subscription{Email: email,
		City:      city,
		Frequency: frequency,
		Token:     token,
		Confirmed: false,
	}

	if err := ss.subscriptionRepository.Create(newSubscription); err != nil {
		log.Println("Db error : ", err.Error())
		return ErrFailedToSaveSubscription
	}

	log.Printf("Subscription created for email %s with token %s", email, token)

	job := EmailJob{
		To:           newSubscription.Email,
		EmailType:    EmailTypeCreateSubscription,
		Subscription: newSubscription,
	}

	err := ss.mailPublisher.Publish(rabbitmq.SendEmail, job)

	if err != nil {
		log.Printf("Failed to publish email job: %v", err)
		return errors.New("failed to publish email job")
	}

	log.Printf("Confirmation email sent to %s", email)

	return nil
}

func (ss *SubscribeService) ConfirmSubscription(token string) error {

	sub, err := ss.subscriptionRepository.FindByToken(token)

	if err != nil {
		return ErrTokenNotFound
	}

	sub.Confirmed = true

	if err := ss.subscriptionRepository.Update(*sub); err != nil {
		log.Println("Failed to update subscription:", err)
		return ErrFailedToSaveSubscription
	}

	job := EmailJob{
		To:           sub.Email,
		EmailType:    EmailTypeConfirmSuccess,
		Subscription: *sub,
	}

	err = ss.mailPublisher.Publish(rabbitmq.SendEmail, job)

	if err != nil {
		log.Printf("Failed to publish email job: %v", err)
		return errors.New("failed to publish email job")
	}

	return nil
}

func (ss *SubscribeService) Unsubscribe(token string) error {

	sub, err := ss.subscriptionRepository.FindByToken(token)

	if err != nil {
		return ErrTokenNotFound
	}

	if sub == nil {
		return nil
	}

	if err := ss.subscriptionRepository.Delete(*sub); err != nil {
		log.Println("Failed to delete subscription:", err)
		return ErrInvalidInput
	}

	return nil
}

func (ss *SubscribeService) emailSubscribed(email string) bool {
	_, err := ss.subscriptionRepository.FindByEmail(email)

	return err == nil
}

func (ss *SubscribeService) generateToken() string {
	return uuid.New().String()
}

func (ss *SubscribeService) SendSubscriptionEmails(freq Frequency) {
	subs := ss.GetConfirmedSubscriptionsByFrequency(freq)
	log.Printf("Found %d %s subscriptions", len(subs), string(freq))

	for _, sub := range subs {
		weather, err := ss.weatherService.GetWeather(sub.City)
		if err != nil {
			log.Println("Weather error for", sub.City, ":", err)
			continue
		}

		job := WeatherUpdateJob{
			To:           sub.Email,
			Subscription: sub,
			Weather:      *weather}

		if err := ss.mailPublisher.Publish(rabbitmq.WeatherUpdate, job); err != nil {
			log.Printf("Failed to publish weather update for %s: %v", sub.Email, err)
		}
	}
}

func (ss *SubscribeService) GetConfirmedSubscriptionsByFrequency(freq Frequency) []Subscription {
	subs, err := ss.subscriptionRepository.FindByFrequencyAndConfirmation(freq)

	if err != nil {
		log.Println("Error fetching confirmed subscriptions:", err)
		return []Subscription{}
	}

	return subs
}
