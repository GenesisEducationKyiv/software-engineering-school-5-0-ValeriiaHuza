package emailBuilder

import (
	"fmt"
	"html"
	"log"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/subscription"
)

type WeatherEmailBuilder struct {
	appUrl string
}

func NewWeatherEmailBuilder(appUrl string) *WeatherEmailBuilder {
	return &WeatherEmailBuilder{
		appUrl: appUrl,
	}
}

func (w *WeatherEmailBuilder) BuildWeatherUpdateEmail(
	sub subscription.Subscription,
	weather client.WeatherDTO,
	time time.Time) string {

	unsubscribeLink := w.buildURL("/api/unsubscribe/") + sub.Token

	sub.City = html.EscapeString(sub.City)
	weather.Description = html.EscapeString(weather.Description)

	return fmt.Sprintf(`
		<p><strong>Weather update for %s</strong></p>
		<p><strong>Date:</strong> %s<br>
		<strong>Time:</strong> %s</p>
		<p><strong>Temperature:</strong> %.1f°C<br>
		<strong>Humidity:</strong> %.0f%%<br>
		<strong>Description:</strong> %s</p>
		<p><a href="%s">Unsubscribe here</a></p>`,
		sub.City,
		time.Format("January 2, 2006"),
		time.Format("15:04"),
		weather.Temperature,
		weather.Humidity,
		weather.Description,
		unsubscribeLink,
	)
}

func (w *WeatherEmailBuilder) BuildConfirmationEmail(sub subscription.Subscription) string {
	confirmationLink := w.buildURL("/api/confirm/") + sub.Token

	log.Println(confirmationLink)

	return fmt.Sprintf(`
		<p>Hello from Weather Updates!</p>
		<p>You subscribed for <strong>%s</strong> updates for <strong>%s</strong> weather.</p>
		<p>Please confirm your subscription by clicking the link below:</p>
		<p><a href="%s">Your link</a></p>`,
		string(sub.Frequency), sub.City, confirmationLink)
}

func (w *WeatherEmailBuilder) BuildConfirmSuccessEmail(sub subscription.Subscription) string {
	unsubscribeLink := w.buildURL("/api/unsubscribe/") + sub.Token
	return fmt.Sprintf(`
		<p>Hello from Weather Updates!</p>
		<p>You have successfully confirmed your subscription!</p>
		<p>If you want to unsubscribe, click the link below:</p>
		<p><a href="%s">Your link</a></p>`,
		unsubscribeLink)
}

func (w *WeatherEmailBuilder) buildURL(path string) string {
	log.Println("Building URL with appUrl:", w.appUrl, "and path:", path)
	return w.appUrl + path
}
