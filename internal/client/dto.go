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

type WeatherAPIErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  float64 `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}
