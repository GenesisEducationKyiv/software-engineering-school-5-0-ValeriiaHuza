//go:build unit
// +build unit

package mailer

import (
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/logger"
	"github.com/stretchr/testify/mock"
	"gopkg.in/gomail.v2"
)

// --- Mocks ---

type mockEmailBuilder struct {
	mock.Mock
}

func (m *mockEmailBuilder) BuildWeatherUpdateEmail(sub SubscriptionDTO, weather WeatherDTO, t time.Time) string {
	args := m.Called(sub, weather, t)
	return args.String(0)
}
func (m *mockEmailBuilder) BuildConfirmationEmail(sub SubscriptionDTO) string {
	args := m.Called(sub)
	return args.String(0)
}
func (m *mockEmailBuilder) BuildConfirmSuccessEmail(sub SubscriptionDTO) string {
	args := m.Called(sub)
	return args.String(0)
}

type mockDialer struct {
	mock.Mock
}

func (m *mockDialer) DialAndSend(msg ...*gomail.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

// --- Tests ---

func setupMailerTest(t *testing.T) (*mockEmailBuilder, *mockDialer, *MailService) {
	builder := new(mockEmailBuilder)
	dialer := new(mockDialer)
	mockLog, _ := logger.NewTestLogger()
	ms := NewMailerService("test@example.com", dialer, builder, *mockLog)
	return builder, dialer, ms
}

func TestSendConfirmationEmail(t *testing.T) {
	builder, dialer, ms := setupMailerTest(t)

	sub := SubscriptionDTO{Email: "user@example.com"}
	expectedBody := "confirmation"

	builder.On("BuildConfirmationEmail", sub).Return(expectedBody)
	dialer.On("DialAndSend", mock.Anything).Return(nil)

	ms.SendConfirmationEmail(sub)

	builder.AssertCalled(t, "BuildConfirmationEmail", sub)
	dialer.AssertCalled(t, "DialAndSend", mock.Anything)
}

func TestSendConfirmSuccessEmail(t *testing.T) {
	builder, dialer, ms := setupMailerTest(t)

	sub := SubscriptionDTO{Email: "user@example.com"}
	expectedBody := "success"

	builder.On("BuildConfirmSuccessEmail", sub).Return(expectedBody)
	dialer.On("DialAndSend", mock.Anything).Return(nil)

	ms.SendConfirmSuccessEmail(sub)

	builder.AssertCalled(t, "BuildConfirmSuccessEmail", sub)
	dialer.AssertCalled(t, "DialAndSend", mock.Anything)
}

func TestSendWeatherUpdateEmail(t *testing.T) {
	builder, dialer, ms := setupMailerTest(t)

	sub := SubscriptionDTO{Email: "user@example.com"}
	weather := WeatherDTO{Temperature: 20}
	expectedBody := "weather update"

	builder.On("BuildWeatherUpdateEmail", sub, weather, mock.AnythingOfType("time.Time")).Return(expectedBody)
	dialer.On("DialAndSend", mock.Anything).Return(nil)

	ms.SendWeatherUpdateEmail(sub, weather)

	builder.AssertCalled(t, "BuildWeatherUpdateEmail", sub, weather, mock.AnythingOfType("time.Time"))
	dialer.AssertCalled(t, "DialAndSend", mock.Anything)
}
