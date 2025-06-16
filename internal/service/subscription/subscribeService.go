package subscription

import (
	"log"
	"regexp"

	appErr "github.com/ValeriiaHuza/weather_api/error"
	"github.com/ValeriiaHuza/weather_api/internal/service/weather"

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
	SendWeatherUpdateEmail(sub Subscription, weather weather.WeatherDTO)
}

type weatherService interface {
	GetWeather(city string) (*weather.WeatherDTO, *appErr.AppError)
}

type SubscribeService struct {
	WeatherService         weatherService
	MailService            mailService
	SubscriptionRepository subscriptionRepository
}

func NewSubscribeService(weatherService weatherService,
	mailService mailService,
	repository subscriptionRepository) *SubscribeService {
	return &SubscribeService{
		WeatherService:         weatherService,
		MailService:            mailService,
		SubscriptionRepository: repository,
	}
}

func (ss *SubscribeService) SubscribeForWeatherUpdates(email string,
	city string, frequencyStr string) *appErr.AppError {

	frequency, err := ss.validateSubscriptionInput(email, city, frequencyStr)
	if err != nil {
		return err
	}

	token := ss.generateToken()

	newSubscription := Subscription{Email: email,
		City:      city,
		Frequency: frequency,
		Token:     token,
		Confirmed: false,
	}

	if err := ss.SubscriptionRepository.Create(newSubscription); err != nil {
		log.Println("Db error : ", err.Error())
		return appErr.ErrFailedToSaveSubscription
	}

	ss.MailService.SendConfirmationEmail(newSubscription)

	return nil
}

func (ss *SubscribeService) validateSubscriptionInput(email string,
	city string, frequencyStr string) (Frequency, *appErr.AppError) {
	if email == "" || city == "" || frequencyStr == "" {
		return "", appErr.ErrInvalidInput
	}

	if !ss.isValidEmail(email) {
		return "", appErr.ErrInvalidInput
	}

	if _, err := ss.WeatherService.GetWeather(city); err != nil {
		return "", err
	}

	frequency, err := ParseFrequency(frequencyStr)
	if err != nil {
		return "", appErr.ErrInvalidInput
	}

	subscribed, err := ss.emailSubscribed(email)
	if subscribed {
		return "", appErr.ErrEmailSubscribed
	}
	if err != nil {
		return "", appErr.ErrInvalidInput
	}

	return frequency, nil
}

func (ss *SubscribeService) ConfirmSubscription(token string) *appErr.AppError {
	if token == "" {
		return appErr.ErrInvalidToken
	}

	sub, err := ss.SubscriptionRepository.FindByToken(token)

	if err != nil {
		return appErr.ErrTokenNotFound
	}

	sub.Confirmed = true

	if err := ss.SubscriptionRepository.Update(*sub); err != nil {
		log.Println("Failed to update subscription:", err)
		return appErr.ErrFailedToSaveSubscription
	}

	ss.MailService.SendConfirmSuccessEmail(*sub)

	return nil
}

func (ss *SubscribeService) Unsubscribe(token string) *appErr.AppError {
	if token == "" {
		return appErr.ErrInvalidToken
	}

	sub, err := ss.SubscriptionRepository.FindByToken(token)

	if err != nil {
		return appErr.ErrTokenNotFound
	}

	if err := ss.SubscriptionRepository.Delete(*sub); err != nil {
		log.Println("Failed to delete subscription:", err)
		return appErr.ErrInvalidInput
	}

	return nil
}

func (ss *SubscribeService) emailSubscribed(email string) (bool, error) {

	sub, err := ss.SubscriptionRepository.FindByEmail(email)

	return sub != nil, err
}

func (ss *SubscribeService) isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (ss *SubscribeService) generateToken() string {
	return uuid.New().String()
}

func (ss *SubscribeService) SendSubscriptionEmails(freq Frequency) {
	subs := ss.GetConfirmedSubscriptionsByFrequency(freq)
	log.Printf("Number of subscriptions found for frequency %s: %d", string(freq), len(subs))
	log.Printf("Found %d %s subscriptions.", len(subs), string(freq))

	for _, sub := range subs {
		weather, err := ss.WeatherService.GetWeather(sub.City)
		if err != nil {
			log.Println("Weather error for", sub.City, ":", err)
			continue
		}

		ss.MailService.SendWeatherUpdateEmail(sub, *weather)
	}
}

func (ss *SubscribeService) GetConfirmedSubscriptionsByFrequency(freq Frequency) []Subscription {
	subs, err := ss.SubscriptionRepository.FindByFrequencyAndConfirmation(freq)

	if err != nil {
		log.Println("Error fetching confirmed subscriptions:", err)
		return []Subscription{}
	}

	return subs
}
