package service

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/ValeriiaHuza/weather_api/db"
	"github.com/ValeriiaHuza/weather_api/error"
	"github.com/ValeriiaHuza/weather_api/models"
	"github.com/ValeriiaHuza/weather_api/utils"
	"github.com/google/uuid"
)

type SubscribeService struct {
	Service *WeatherService
}

func NewSubscribeService(service *WeatherService) *SubscribeService {
	return &SubscribeService{Service: service}
}

func (ss *SubscribeService) SubscribeForWeather(email string, city string, frequencyStr string) *error.AppError {

	if email == "" || city == "" || frequencyStr == "" {
		return error.ErrInvalidInput
	}

	if !ss.isValidEmail(email) {
		return error.ErrInvalidInput
	}

	if _, err := ss.Service.GetWeather(city); err != nil {
		return error.ErrInvalidInput
	}

	frequency, err := models.ParseFrequency(frequencyStr)
	if err != nil {
		return error.ErrInvalidInput
	}

	//check subscription in db
	if ss.emailSubscribed(email) {
		return error.ErrEmailSubscribed
	}

	token := ss.generateToken()

	newSubscription := models.Subscription{Email: email,
		City:      city,
		Frequency: frequency,
		Token:     token,
		Confirmed: false,
	}

	if err := db.DB.Create(&newSubscription).Error; err != nil {
		return error.ErrInvalidInput
	}

	utils.SendConfirmationEmail(newSubscription)

	return nil
}

func (ss *SubscribeService) ConfirmSubscription(token string) *error.AppError {
	if token == "" {
		return error.ErrInvalidToken
	}

	var sub models.Subscription
	//find token in db
	result := db.DB.Where("token = ?", token).First(&sub)

	if result.Error != nil {
		return error.ErrTokenNotFound
	}

	sub.Confirmed = true

	if err := db.DB.Save(&sub).Error; err != nil {
		return error.ErrInvalidToken
	}

	utils.SendConfirmSuccessMail(sub)

	return nil
}

func (ss *SubscribeService) Unsubscribe(token string) *error.AppError {
	if token == "" {
		return error.ErrInvalidToken
	}

	var sub models.Subscription
	//find token in db
	result := db.DB.Where("token = ?", token).First(&sub)

	if result.Error != nil {
		return error.ErrTokenNotFound
	}

	if err := db.DB.Unscoped().Delete(&sub).Error; err != nil {
		return error.ErrInvalidToken
	}

	return nil
}

func (*SubscribeService) emailSubscribed(email string) bool {
	var sub models.Subscription

	result := db.DB.Where("email = ?", email).First(&sub)

	return result.Error == nil
}

func (*SubscribeService) isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (*SubscribeService) generateToken() string {
	return uuid.New().String()
}

func SendEmails(freq models.Frequency) {
	subs := GetSubscriptionsByFrequency(freq)
	serviceW := NewWeatherService()

	log.Printf("Found %d %s subscriptions.", len(subs), string(freq))

	for _, sub := range subs {
		weather, err := serviceW.GetWeather(sub.City)
		if err != nil {
			log.Println("Weather error for", sub.City, ":", err)
			continue
		}

		unsubscribeLink := utils.BuildURL("/api/unsubscribe/") + sub.Token

		now := time.Now()
		message := fmt.Sprintf(`
		<p><strong>Weather update for %s</strong></p>
		<p><strong>Date:</strong> %s<br>
		<strong>Time:</strong> %s</p>
		<p><strong>Temperature:</strong> %.1f°C<br>
		<strong>Humidity:</strong> %.0f%%<br>
		<strong>Description:</strong> %s</p>
		<p><a href="%s">Unsubscribe here</a></p>`,
			sub.City,
			now.Format("January 2, 2006"),
			now.Format("15:04"),
			weather.Temperature,
			weather.Humidity,
			weather.Description,
			unsubscribeLink,
		)

		utils.SendEmail(sub.Email, "Weather Update", message)
	}
}

func GetSubscriptionsByFrequency(freq models.Frequency) []models.Subscription {
	var subs []models.Subscription
	result := db.DB.Where("frequency = ? AND confirmed = true", freq).Find(&subs)
	if result.Error != nil {
		log.Println("DB error:", result.Error)
	}
	return subs
}
