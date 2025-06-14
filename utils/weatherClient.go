package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ValeriiaHuza/weather_api/dto"
	appErr "github.com/ValeriiaHuza/weather_api/error"
)

type WeatherAPIClient interface {
	FetchWeather(city string) ([]byte, *appErr.AppError)
}

type WeatherAPIClientImpl struct {
}

func (c *WeatherAPIClientImpl) FetchWeather(city string) ([]byte, *appErr.AppError) {
	apiKey := os.Getenv("WEATHER_API_KEY")

	if apiKey == "" {
		log.Println("Missing WEATHER_API_KEY")
		return nil, appErr.ErrInvalidRequest
	}

	city = url.QueryEscape(city)

	weatherUrl := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%v&q=%v", apiKey, city)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(weatherUrl)
	if err != nil {
		log.Println("HTTP request failed:", err)
		return nil, appErr.ErrInvalidRequest
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return nil, appErr.ErrInvalidRequest
	}

	if apiErr := c.ParseAPIError(body); apiErr != nil {
		return nil, apiErr
	}

	return body, nil
}

func (ws *WeatherAPIClientImpl) ParseAPIError(body []byte) *appErr.AppError {
	var apiErr dto.APIErrorResponse
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return nil
	}

	if apiErr.Error.Message != "" {
		log.Printf("API Error %d: %s\n", apiErr.Error.Code, apiErr.Error.Message)
		if apiErr.Error.Code == 1006 {
			return appErr.ErrCityNotFound
		}
		return appErr.ErrInvalidRequest
	}
	return nil
}
