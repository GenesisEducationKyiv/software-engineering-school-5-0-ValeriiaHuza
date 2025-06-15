package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/ValeriiaHuza/weather_api/dto"
	"github.com/ValeriiaHuza/weather_api/models"
)

type WeatherEmailBuilder struct{}

func NewWeatherEmailBuilder() *WeatherEmailBuilder {
	return &WeatherEmailBuilder{}
}

func (w *WeatherEmailBuilder) BuildWeatherUpdateEmail(sub models.Subscription, weather dto.WeatherDTO) string {
	unsubscribeLink := w.BuildURL("/api/unsubscribe/") + sub.Token
	now := time.Now()
	return fmt.Sprintf(`
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
}

func (w *WeatherEmailBuilder) BuildConfirmationEmail(sub models.Subscription) string {
	confirmationLink := w.BuildURL("/api/confirm/") + sub.Token
	return fmt.Sprintf(`
		<p>Hello from Weather Updates!</p>
		<p>You subscribed for <strong>%s</strong> updates for <strong>%s</strong> weather.</p>
		<p>Please confirm your subscription by clicking the link below:</p>
		<p><a href="%s">Your link</a></p>`,
		string(sub.Frequency), sub.City, confirmationLink)
}

func (w *WeatherEmailBuilder) BuildConfirmSuccessEmail(sub models.Subscription) string {
	unsubscribeLink := w.BuildURL("/api/unsubscribe/") + sub.Token
	return fmt.Sprintf(`
		<p>Hello from Weather Updates!</p>
		<p>You have successfully confirmed your subscription!</p>
		<p>If you want to unsubscribe, click the link below:</p>
		<p><a href="%s">Your link</a></p>`,
		unsubscribeLink)
}

func (w *WeatherEmailBuilder) BuildURL(path string) string {
	host := os.Getenv("APP_URL")
	return host + path
}
