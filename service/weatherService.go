package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ValeriiaHuza/weather_api/dto"
	"github.com/ValeriiaHuza/weather_api/error"
)

type WeatherService struct {
}

func NewWeatherService() *WeatherService {
	return &WeatherService{}
}

func (ws *WeatherService) GetWeather(city string) (*dto.WeatherDTO, *error.AppError) {
	if city == "" {
		return nil, error.ErrInvalidRequest
	}

	body, err := ws.fetchWeatherData(city)
	if err != nil {
		return nil, err
	}

	if apiErr := ws.parseAPIError(body); apiErr != nil {
		return nil, apiErr
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

func (ws *WeatherService) fetchWeatherData(city string) ([]byte, *error.AppError) {
	apiKey := os.Getenv("WEATHER_API_KEY")

	if apiKey == "" {
		log.Println("Missing WEATHER_API_KEY")
		return nil, error.ErrInvalidRequest
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%v&q=%v", apiKey, city)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("HTTP request failed:", err)
		return nil, error.ErrInvalidRequest
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return nil, error.ErrInvalidRequest
	}

	return body, nil
}

func (ws *WeatherService) parseAPIError(body []byte) *error.AppError {
	var apiErr dto.APIErrorResponse
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return nil
	}

	if apiErr.Error.Message != "" {
		log.Printf("API Error %d: %s\n", apiErr.Error.Code, apiErr.Error.Message)
		if apiErr.Error.Code == 1006 {
			return error.ErrCityNotFound
		}
		return error.ErrInvalidRequest
	}
	return nil
}
