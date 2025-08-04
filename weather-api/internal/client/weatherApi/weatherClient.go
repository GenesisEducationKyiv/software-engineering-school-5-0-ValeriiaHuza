package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
)

type WeatherAPIClient struct {
	apiKey string
	apiUrl string
	client *http.Client
}

func NewWeatherAPIClient(apiKey string, apiUrl string, http *http.Client) *WeatherAPIClient {
	return &WeatherAPIClient{
		apiKey: apiKey,
		apiUrl: strings.TrimRight(apiUrl, "/"),
		client: http,
	}
}

func (c *WeatherAPIClient) FetchWeather(city string) (*client.WeatherDTO, error) {
	city = url.QueryEscape(city)

	weatherURL := fmt.Sprintf("%s/current.json?key=%s&q=%s", c.apiUrl, c.apiKey, city)

	log.Printf("Sending request to weather API for city: %s", city)

	resp, err := c.client.Get(weatherURL)

	if err != nil {
		log.Printf(" HTTP request failed: %v", err)
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

	weatherDTO := client.WeatherDTO{
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
		log.Printf("Weather API Error %d: %s\n", apiErr.Error.Code, apiErr.Error.Message)
		if apiErr.Error.Code == 1006 {
			return client.ErrCityNotFound
		}
		return client.ErrInvalidRequest
	}
	return nil
}
