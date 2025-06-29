package subscription

import (
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client"
	"github.com/google/uuid"
)

type subscriptionRepository interface {
	Create(sub Subscription) error
	Update(sub Subscription) error
	FindByToken(token string) (*Subscription, error)
	Delete(sub Subscription) error
	FindByEmail(email string) (*Subscription, error)
	FindByFrequencyAndConfirmation(freq Frequency) ([]Subscription, error)
}

type mailService interface {
	SendConfirmationEmail(sub Subscription)
	SendConfirmSuccessEmail(sub Subscription)
	SendWeatherUpdateEmail(sub Subscription, weather client.WeatherDTO)
}

type weatherService interface {
	GetWeather(city string) (*client.WeatherDTO, error)
}

type SubscribeService struct {
	weatherService         weatherService
	mailService            mailService
	subscriptionRepository subscriptionRepository
}

func NewSubscribeService(weatherService weatherService,
	mailService mailService,
	repository subscriptionRepository) *SubscribeService {
	return &SubscribeService{
		weatherService:         weatherService,
		mailService:            mailService,
		subscriptionRepository: repository,
	}
}

func (ss *SubscribeService) SubscribeForWeatherUpdates(email string,
	city string, frequency Frequency) error {

	if _, err := ss.weatherService.GetWeather(city); err != nil {
		return err
	}

	subscribed, err := ss.emailSubscribed(email)
	if subscribed {
		return ErrEmailAlreadySubscribed
	}
	if err != nil {
		return ErrInvalidInput
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

	ss.mailService.SendConfirmationEmail(newSubscription)

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

	ss.mailService.SendConfirmSuccessEmail(*sub)

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

func (ss *SubscribeService) emailSubscribed(email string) (bool, error) {

	sub, err := ss.subscriptionRepository.FindByEmail(email)

	return sub != nil, err
}

func (ss *SubscribeService) generateToken() string {
	return uuid.New().String()
}

func (ss *SubscribeService) SendSubscriptionEmails(freq Frequency) {
	subs := ss.GetConfirmedSubscriptionsByFrequency(freq)
	log.Printf("Number of subscriptions found for frequency %s: %d", string(freq), len(subs))
	log.Printf("Found %d %s subscriptions.", len(subs), string(freq))

	for _, sub := range subs {
		weather, err := ss.weatherService.GetWeather(sub.City)
		if err != nil {
			log.Println("Weather error for", sub.City, ":", err)
			continue
		}

		ss.mailService.SendWeatherUpdateEmail(sub, *weather)
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
