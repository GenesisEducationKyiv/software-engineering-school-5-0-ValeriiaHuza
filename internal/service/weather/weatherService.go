package weather

import (
	"encoding/json"
	"log"

	appErr "github.com/ValeriiaHuza/weather_api/error"
)

type weatherAPIClient interface {
	FetchWeather(city string) ([]byte, *appErr.AppError)
}

type WeatherService struct {
	WeatherClient weatherAPIClient
}

func NewWeatherAPIService(weatherClient weatherAPIClient) *WeatherService {
	return &WeatherService{
		WeatherClient: weatherClient,
	}
}

func (ws *WeatherService) GetWeather(city string) (*WeatherDTO, *appErr.AppError) {
	if city == "" {
		return nil, appErr.ErrInvalidRequest
	}

	body, err := ws.WeatherClient.FetchWeather(city)
	if err != nil {
		log.Println("HTTP error:", err)
		return nil, err
	}

	var weather WeatherResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		log.Println("Failed to parse JSON:", err)
		return nil, appErr.ErrInvalidRequest
	}

	weatherDTO := WeatherDTO{
		Temperature: weather.Current.TempC,
		Humidity:    weather.Current.Humidity,
		Description: weather.Current.Condition.Text,
	}

	return &weatherDTO, nil
}
