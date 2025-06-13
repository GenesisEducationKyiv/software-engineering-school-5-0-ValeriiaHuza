package error

import "net/http"

type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func New(status int, message string) *AppError {
	return &AppError{
		StatusCode: status,
		Message:    message,
	}
}

var (
	ErrCityNotFound   = New(http.StatusNotFound, "City not found")
	ErrInvalidRequest = New(http.StatusBadRequest, "Invalid request")
	ErrInvalidInput   = New(http.StatusBadRequest, "Invalid input")

	ErrEmailSubscribed = New(http.StatusConflict, "Email already subscribed")

	ErrInvalidToken  = New(http.StatusBadRequest, "Invalid token")
	ErrTokenNotFound = New(http.StatusNotFound, "Token not found")
)
