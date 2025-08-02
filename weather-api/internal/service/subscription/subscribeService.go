package subscription

import (
	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

	logger.GetLogger().Info("Subscribing for weather updates",
		zap.String("email", email),
		zap.String("city", city),
		zap.String("frequency", string(frequency)))

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
		logger.GetLogger().Error("Failed to create subscription",
			zap.String("email", email),
			zap.Error(err))
		return ErrFailedToSaveSubscription
	}

	logger.GetLogger().Info("Subscription created",
		zap.String("email", email))

	job := EmailJob{
		To:           newSubscription.Email,
		EmailType:    EmailTypeCreateSubscription,
		Subscription: newSubscription,
	}

	err := ss.mailPublisher.Publish(rabbitmq.SendEmail, job)

	if err != nil {
		logger.GetLogger().Error("Failed to publish email job",
			zap.String("email", newSubscription.Email),
			zap.Error(err))
		return errors.New("failed to publish email job")
	}

	logger.GetLogger().Info("Confirmation email sent",
		zap.String("email", email))

	return nil
}

func (ss *SubscribeService) ConfirmSubscription(token string) error {

	sub, err := ss.subscriptionRepository.FindByToken(token)

	if err != nil {
		return ErrTokenNotFound
	}

	sub.Confirmed = true

	if err := ss.subscriptionRepository.Update(*sub); err != nil {
		logger.GetLogger().Error("Failed to update subscription",
			zap.String("token", token),
			zap.Error(err))

		return ErrFailedToSaveSubscription
	}

	job := EmailJob{
		To:           sub.Email,
		EmailType:    EmailTypeConfirmSuccess,
		Subscription: *sub,
	}

	err = ss.mailPublisher.Publish(rabbitmq.SendEmail, job)

	if err != nil {
		logger.GetLogger().Error("Failed to publish email job",
			zap.String("email", sub.Email),
			zap.Error(err))
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
		logger.GetLogger().Error("Failed to delete subscription",
			zap.String("token", token),
			zap.Error(err))
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
	logger.GetLogger().Info("Sending subscription emails",
		zap.String("frequency", string(freq)),
		zap.Int("count", len(subs)))

	for _, sub := range subs {
		weather, err := ss.weatherService.GetWeather(sub.City)
		if err != nil {
			logger.GetLogger().Error("Failed to fetch weather data",
				zap.String("city", sub.City),
				zap.Error(err))
			continue
		}

		job := WeatherUpdateJob{
			To:           sub.Email,
			Subscription: sub,
			Weather:      *weather}

		if err := ss.mailPublisher.Publish(rabbitmq.WeatherUpdate, job); err != nil {
			logger.GetLogger().Error("Failed to publish weather update",
				zap.String("email", sub.Email),
				zap.Error(err))
		}
	}
}

func (ss *SubscribeService) GetConfirmedSubscriptionsByFrequency(freq Frequency) []Subscription {
	subs, err := ss.subscriptionRepository.FindByFrequencyAndConfirmation(freq)

	if err != nil {
		logger.GetLogger().Error("Failed to fetch confirmed subscriptions",
			zap.String("frequency", string(freq)),
			zap.Error(err))
		return []Subscription{}
	}

	return subs
}
