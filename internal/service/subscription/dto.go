package subscription

import "errors"

var (
	ErrInvalidRequest           = errors.New("invalid request")
	ErrInvalidInput             = errors.New("invalid input")
	ErrEmailAlreadySubscribed   = errors.New("email already subscribed")
	ErrInvalidToken             = errors.New("invalid token")
	ErrTokenNotFound            = errors.New("token not found")
	ErrFailedToSaveSubscription = errors.New("failed to save subscription")
)
