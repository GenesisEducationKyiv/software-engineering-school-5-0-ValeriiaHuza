package openweather

import (
	"net/http"

	"github.com/ValeriiaHuza/weather_api/internal/client"
)

type WeatherAPIClient struct {
	apiKey string
	apiUrl string
	client *http.Client
}

func NewWeatherAPIClient(apiKey string, apiUrl string, http *http.Client) *WeatherAPIClient {
	return &WeatherAPIClient{
		apiKey: apiKey,
		apiUrl: apiUrl,
		client: http,
	}
}

func (c *WeatherAPIClient) FetchWeather(city string) (*client.WeatherDTO, error) {

	return nil, nil
}
