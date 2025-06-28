package weather

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/internal/client"
)

type weatherChain interface {
	GetWeather(city string) (*client.WeatherDTO, error)
}

type WeatherService struct {
	weatherChain weatherChain
}

func NewWeatherAPIService(weatherChain weatherChain) *WeatherService {
	return &WeatherService{
		weatherChain: weatherChain,
	}
}

func (ws *WeatherService) GetWeather(city string) (*client.WeatherDTO, error) {

	weatherDto, err := ws.weatherChain.GetWeather(city)
	if err != nil {
		log.Println("HTTP error in GetWeather :", err)
		return nil, err
	}

	return weatherDto, nil
}
