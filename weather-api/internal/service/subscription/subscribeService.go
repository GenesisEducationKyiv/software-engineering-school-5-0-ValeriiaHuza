package subscription

import (
	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
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
	logger                 logger.Logger
}

func NewSubscribeService(weatherService weatherService,
	repository subscriptionRepository,
	mailPublisher mailPublisher, logger logger.Logger) *SubscribeService {
	return &SubscribeService{
		weatherService:         weatherService,
		subscriptionRepository: repository,
		mailPublisher:          mailPublisher,
		logger:                 logger,
	}
}

func (ss *SubscribeService) SubscribeForWeatherUpdates(email string,
	city string, frequency Frequency) error {

	if _, err := ss.weatherService.GetWeather(city); err != nil {
		return err
	}

	ss.logger.Info("Validating subscription input",
		"email", email,
		"city", city,
		"frequency", frequency)

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
		ss.logger.Error("Failed to create subscription",
			"email", email,
			"error", err)
		return ErrFailedToSaveSubscription
	}

	ss.logger.Info("Subscription created", "email", email)

	job := EmailJob{
		To:           newSubscription.Email,
		EmailType:    EmailTypeCreateSubscription,
		Subscription: newSubscription,
	}

	err := ss.mailPublisher.Publish(rabbitmq.SendEmail, job)

	if err != nil {
		ss.logger.Error("Failed to publish email job",
			"email", newSubscription.Email,
			"error", err)
		return errors.New("failed to publish email job")
	}

	ss.logger.Info("Subscription email job published", "email", email)

	return nil
}

func (ss *SubscribeService) ConfirmSubscription(token string) error {

	sub, err := ss.subscriptionRepository.FindByToken(token)

	if err != nil {
		return ErrTokenNotFound
	}

	sub.Confirmed = true

	if err := ss.subscriptionRepository.Update(*sub); err != nil {
		ss.logger.Error("Failed to update subscription",
			"token", token,
			"error", err)

		return ErrFailedToSaveSubscription
	}

	job := EmailJob{
		To:           sub.Email,
		EmailType:    EmailTypeConfirmSuccess,
		Subscription: *sub,
	}

	err = ss.mailPublisher.Publish(rabbitmq.SendEmail, job)

	if err != nil {
		ss.logger.Error("Failed to publish confirmation email job",
			"email", sub.Email,
			"error", err)
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
		ss.logger.Error("Failed to delete subscription",
			"token", token,
			"error", err)

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
	ss.logger.Info("Sending subscription emails",
		"frequency", string(freq),
		"count", len(subs))

	for _, sub := range subs {
		weather, err := ss.weatherService.GetWeather(sub.City)
		if err != nil {
			ss.logger.Error("Failed to fetch weather data",
				"city", sub.City,
				"error", err)
			continue
		}

		job := WeatherUpdateJob{
			To:           sub.Email,
			Subscription: sub,
			Weather:      *weather}

		if err := ss.mailPublisher.Publish(rabbitmq.WeatherUpdate, job); err != nil {
			ss.logger.Error("Failed to publish weather update",
				"email", sub.Email,
				"error", err)
		}
	}
}

func (ss *SubscribeService) GetConfirmedSubscriptionsByFrequency(freq Frequency) []Subscription {
	subs, err := ss.subscriptionRepository.FindByFrequencyAndConfirmation(freq)

	if err != nil {
		ss.logger.Error("Failed to fetch confirmed subscriptions",
			"frequency", string(freq),
			"error", err)
		return []Subscription{}
	}

	return subs
}
