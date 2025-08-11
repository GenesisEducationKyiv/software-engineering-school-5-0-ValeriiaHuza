package client

import "errors"

var (
	ErrCityNotFound   = errors.New("city not found")
	ErrInvalidRequest = errors.New("invalid request")
)

type WeatherDTO struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}
