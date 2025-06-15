package service

import (
	"encoding/json"
	"log"

	"github.com/ValeriiaHuza/weather_api/dto"
	"github.com/ValeriiaHuza/weather_api/error"
	"github.com/ValeriiaHuza/weather_api/utils"
)

type WeatherService interface {
	GetWeather(city string) (*dto.WeatherDTO, *error.AppError)
}

type WeatherAPIService struct {
	WeatherClient utils.WeatherAPIClient
}

func NewWeatherAPIService(weatherClient utils.WeatherAPIClient) *WeatherAPIService {
	return &WeatherAPIService{
		WeatherClient: weatherClient,
	}
}

func (ws *WeatherAPIService) GetWeather(city string) (*dto.WeatherDTO, *error.AppError) {
	if city == "" {
		return nil, error.ErrInvalidRequest
	}

	body, err := ws.WeatherClient.FetchWeather(city)
	if err != nil {
		log.Println("HTTP error:", err)
		return nil, err
	}

	var weather dto.WeatherResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		log.Println("Failed to parse JSON:", err)
		return nil, error.ErrInvalidRequest
	}

	weatherDTO := dto.WeatherDTO{
		Temperature: weather.Current.TempC,
		Humidity:    weather.Current.Humidity,
		Description: weather.Current.Condition.Text,
	}

	return &weatherDTO, nil
}
