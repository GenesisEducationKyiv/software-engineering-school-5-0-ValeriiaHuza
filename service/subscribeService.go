package service

import (
	"log"
	"regexp"

	appErr "github.com/ValeriiaHuza/weather_api/error"
	"github.com/ValeriiaHuza/weather_api/models"
	"github.com/ValeriiaHuza/weather_api/repository"
	"github.com/google/uuid"
)

type SubscribeService interface {
	SubscribeForWeatherUpdates(email, city, frequency string) *appErr.AppError
	ConfirmSubscription(token string) *appErr.AppError
	Unsubscribe(token string) *appErr.AppError
	GetConfirmedSubscriptionsByFrequency(freq models.Frequency) []models.Subscription
	SendSubscriptionEmails(freq models.Frequency)
}

type SubscribeServiceImpl struct {
	WeatherService      WeatherService
	MailerService       MailerService
	SubscribeRepository repository.SubscriptionRepository
}

func NewSubscribeService(weatherService WeatherService, mailerService MailerService, repository repository.SubscriptionRepository) *SubscribeServiceImpl {
	return &SubscribeServiceImpl{
		WeatherService:      weatherService,
		MailerService:       mailerService,
		SubscribeRepository: repository,
	}
}

func (ss *SubscribeServiceImpl) SubscribeForWeatherUpdates(email string, city string, frequencyStr string) *appErr.AppError {

	frequency, err := ss.validateSubscriptionInput(email, city, frequencyStr)
	if err != nil {
		return err
	}

	token := ss.generateToken()

	newSubscription := models.Subscription{Email: email,
		City:      city,
		Frequency: frequency,
		Token:     token,
		Confirmed: false,
	}

	if err := ss.SubscribeRepository.Create(newSubscription); err != nil {
		log.Println("Db error : ", err.Error())
		return appErr.ErrInvalidInput
	}

	ss.MailerService.SendConfirmationEmail(newSubscription)

	return nil
}

func (ss *SubscribeServiceImpl) validateSubscriptionInput(email string, city string, frequencyStr string) (models.Frequency, *appErr.AppError) {
	if email == "" || city == "" || frequencyStr == "" {
		return "", appErr.ErrInvalidInput
	}

	if !ss.isValidEmail(email) {
		return "", appErr.ErrInvalidInput
	}

	if _, err := ss.WeatherService.GetWeather(city); err != nil {
		return "", appErr.ErrInvalidInput
	}

	frequency, err := models.ParseFrequency(frequencyStr)
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

func (ss *SubscribeServiceImpl) ConfirmSubscription(token string) *appErr.AppError {
	if token == "" {
		return appErr.ErrInvalidToken
	}

	sub, err := ss.SubscribeRepository.FindByToken(token)

	if err != nil {
		return appErr.ErrTokenNotFound
	}

	sub.Confirmed = true

	if err := ss.SubscribeRepository.Update(*sub); err != nil {
		return appErr.ErrInvalidToken
	}

	ss.MailerService.SendConfirmSuccessEmail(*sub)

	return nil
}

func (ss *SubscribeServiceImpl) Unsubscribe(token string) *appErr.AppError {
	if token == "" {
		return appErr.ErrInvalidToken
	}

	sub, err := ss.SubscribeRepository.FindByToken(token)

	if err != nil {
		return appErr.ErrTokenNotFound
	}

	if err := ss.SubscribeRepository.Delete(*sub); err != nil {
		return appErr.ErrInvalidToken
	}

	return nil
}

func (ss *SubscribeServiceImpl) emailSubscribed(email string) (bool, error) {

	sub, err := ss.SubscribeRepository.FindByEmail(email)

	return sub != nil, err
}

func (ss *SubscribeServiceImpl) isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (ss *SubscribeServiceImpl) generateToken() string {
	return uuid.New().String()
}

func (ss *SubscribeServiceImpl) SendSubscriptionEmails(freq models.Frequency) {
	subs := ss.GetConfirmedSubscriptionsByFrequency(freq)

	log.Printf("Found %d %s subscriptions.", len(subs), string(freq))

	for _, sub := range subs {
		weather, err := ss.WeatherService.GetWeather(sub.City)
		if err != nil {
			log.Println("Weather error for", sub.City, ":", err)
			continue
		}

		ss.MailerService.SendWeatherUpdateEmail(sub, *weather)
	}
}

func (ss *SubscribeServiceImpl) GetConfirmedSubscriptionsByFrequency(freq models.Frequency) []models.Subscription {
	subs, err := ss.SubscribeRepository.FindByFrequencyAndConfirmation(freq)

	if err != nil {
		log.Println("Error fetching confirmed subscriptions:", err)
		return []models.Subscription{}
	}

	return subs
}
