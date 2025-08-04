//go:build integration
// +build integration

package integration

// import (
// 	"fmt"

// 	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
// 	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/service/subscription"
// )

// type MailService interface {
// 	SendConfirmationEmail(sub subscription.Subscription)
// 	SendConfirmSuccessEmail(sub subscription.Subscription)
// 	SendWeatherUpdateEmail(sub subscription.Subscription, weather client.WeatherDTO)
// }

// type FakeMailService struct {
// 	SentEmail SentEmail
// }

// func NewFakeMailService() *FakeMailService {
// 	return &FakeMailService{}
// }

// type SentEmail struct {
// 	To      string
// 	Subject string
// 	Body    string
// }

// func (f *FakeMailService) SendConfirmationEmail(sub subscription.Subscription) {
// 	// Just record the call instead of sending an email
// 	body := fmt.Sprintf("Fake email to %s for subscription token %s", sub.Email, sub.Token)

// 	f.SentEmail = SentEmail{
// 		To:      sub.Email,
// 		Subject: "Weather updates confirmation link",
// 		Body:    body,
// 	}
// }

// func (f *FakeMailService) SendConfirmSuccessEmail(sub subscription.Subscription) {
// 	body := fmt.Sprintf("Fake confirm success email to %s", sub.Email)
// 	f.SentEmail = SentEmail{
// 		To:      sub.Email,
// 		Subject: "Weather updates subscription",
// 		Body:    body,
// 	}
// }

// func (f *FakeMailService) SendWeatherUpdateEmail(sub subscription.Subscription, weather client.WeatherDTO) {
// 	f.SentEmail = SentEmail{
// 		To:      sub.Email,
// 		Subject: "Weather Update",
// 		Body:    fmt.Sprintf("Weather: %.1f°C, %s", weather.Temperature, weather.Description),
// 	}
// }
