package emailBuilder

import (
	"fmt"
	"html"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/mailer"
)

type loggerInterface interface {
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

type WeatherEmailBuilder struct {
	appUrl string
	logger loggerInterface
}

func NewWeatherEmailBuilder(appUrl string, logger loggerInterface) *WeatherEmailBuilder {
	return &WeatherEmailBuilder{
		appUrl: appUrl,
		logger: logger,
	}
}

func (w *WeatherEmailBuilder) BuildWeatherUpdateEmail(
	sub mailer.SubscriptionDTO,
	weather mailer.WeatherDTO,
	time time.Time) string {

	unsubscribeLink := w.buildURL("/api/unsubscribe/") + sub.Token

	escapedCity := html.EscapeString(sub.City)
	escapedDescription := html.EscapeString(weather.Description)

	return fmt.Sprintf(`
		<p><strong>Weather update for %s</strong></p>
		<p><strong>Date:</strong> %s<br>
		<strong>Time:</strong> %s</p>
		<p><strong>Temperature:</strong> %.1f°C<br>
		<strong>Humidity:</strong> %.0f%%<br>
		<strong>Description:</strong> %s</p>
		<p><a href="%s">Unsubscribe here</a></p>`,
		escapedCity,
		time.Format("January 2, 2006"),
		time.Format("15:04"),
		weather.Temperature,
		weather.Humidity,
		escapedDescription,
		unsubscribeLink,
	)
}

func (w *WeatherEmailBuilder) BuildConfirmationEmail(sub mailer.SubscriptionDTO) string {
	confirmationLink := w.buildURL("/api/confirm/") + sub.Token

	w.logger.Info("Building confirmation email", "confirmationLink", confirmationLink)

	escapedCity := html.EscapeString(sub.City)

	return fmt.Sprintf(`
		<p>Hello from Weather Updates!</p>
		<p>You subscribed for <strong>%s</strong> updates for <strong>%s</strong> weather.</p>
		<p>Please confirm your subscription by clicking the link below:</p>
		<p><a href="%s">Your link</a></p>`,
		string(sub.Frequency), escapedCity, confirmationLink)
}

func (w *WeatherEmailBuilder) BuildConfirmSuccessEmail(sub mailer.SubscriptionDTO) string {
	unsubscribeLink := w.buildURL("/api/unsubscribe/") + sub.Token
	return fmt.Sprintf(`
		<p>Hello from Weather Updates!</p>
		<p>You have successfully confirmed your subscription!</p>
		<p>If you want to unsubscribe, click the link below:</p>
		<p><a href="%s">Your link</a></p>`,
		unsubscribeLink)
}

func (w *WeatherEmailBuilder) buildURL(path string) string {
	w.logger.Info("Building URL", "appUrl", w.appUrl, "path", path)
	return w.appUrl + path
}
