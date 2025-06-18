package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type WeatherAPIClient struct {
	apiKey string
	client *http.Client
}

func NewWeatherAPIClient(apiKey string, http *http.Client) *WeatherAPIClient {
	return &WeatherAPIClient{
		apiKey: apiKey,
		client: http,
	}
}

func (c *WeatherAPIClient) FetchWeather(city string) (*WeatherDTO, error) {
	city = url.QueryEscape(city)

	weatherUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%v&q=%v", c.apiKey, city)

	resp, err := c.client.Get(weatherUrl)

	if err != nil {
		log.Println("HTTP request failed:", err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body:", err)
		return nil, err
	}

	if apiErr := c.parseAPIError(body); apiErr != nil {
		return nil, apiErr
	}

	var weather WeatherAPIResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		log.Println("Failed to parse JSON:", err)
		return nil, err
	}

	weatherDTO := WeatherDTO{
		Temperature: weather.Current.TempC,
		Humidity:    weather.Current.Humidity,
		Description: weather.Current.Condition.Text,
	}

	return &weatherDTO, nil
}

func (ws *WeatherAPIClient) parseAPIError(body []byte) error {
	var apiErr WeatherAPIErrorResponse
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return nil
	}

	if apiErr.Error.Message != "" {
		log.Printf("API Error %d: %s\n", apiErr.Error.Code, apiErr.Error.Message)
		if apiErr.Error.Code == 1006 {
			return ErrCityNotFound
		}
		return ErrInvalidRequest
	}
	return nil
}
