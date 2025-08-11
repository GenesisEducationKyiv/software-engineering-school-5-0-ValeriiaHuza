package subscription

import (
	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
)

var (
	ErrInvalidRequest           = errors.New("invalid request")
	ErrInvalidInput             = errors.New("invalid input")
	ErrEmailAlreadySubscribed   = errors.New("email already subscribed")
	ErrInvalidToken             = errors.New("invalid token")
	ErrTokenNotFound            = errors.New("token not found")
	ErrFailedToSaveSubscription = errors.New("failed to save subscription")
)

type EmailType string

const (
	EmailTypeCreateSubscription EmailType = "CreateSubscription"
	EmailTypeConfirmSuccess     EmailType = "ConfirmSuccess"
)

type EmailJob struct {
	To           string
	EmailType    EmailType
	Subscription Subscription
}

type WeatherUpdateJob struct {
	To           string
	Weather      client.WeatherDTO
	Subscription Subscription
}
