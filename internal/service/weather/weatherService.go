package weather

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/internal/client"
)

type weatherAPIClient interface {
	FetchWeather(city string) (*client.WeatherDTO, error)
}

type WeatherService struct {
	weatherClient weatherAPIClient
}

func NewWeatherAPIService(weatherClient weatherAPIClient) *WeatherService {
	return &WeatherService{
		weatherClient: weatherClient,
	}
}

func (ws *WeatherService) GetWeather(city string) (*client.WeatherDTO, error) {

	weatherDto, err := ws.weatherClient.FetchWeather(city)
	if err != nil {
		log.Println("HTTP error in GetWeather :", err)
		return nil, err
	}

	return weatherDto, nil
}
