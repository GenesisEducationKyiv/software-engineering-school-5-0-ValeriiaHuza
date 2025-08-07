package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
)

type WeatherAPIClient struct {
	apiKey string
	apiUrl string
	client *http.Client
	logger logger.Logger
}

func NewWeatherAPIClient(apiKey string, apiUrl string, http *http.Client, logger logger.Logger) *WeatherAPIClient {
	return &WeatherAPIClient{
		apiKey: apiKey,
		apiUrl: strings.TrimRight(apiUrl, "/"),
		client: http,
		logger: logger,
	}
}

func (c *WeatherAPIClient) FetchWeather(city string) (*client.WeatherDTO, error) {
	city = url.QueryEscape(city)

	weatherURL := fmt.Sprintf("%s/current.json?key=%s&q=%s", c.apiUrl, c.apiKey, city)

	c.logger.Info("Sending request to Weather API", "city", city, "url", weatherURL)

	resp, err := c.client.Get(weatherURL)

	if err != nil {
		c.logger.Error("HTTP request to Weather API failed", "error", err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("Failed to close response body", "error", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body", "error", err)

		return nil, err
	}

	if apiErr := c.parseAPIError(body); apiErr != nil {
		return nil, apiErr
	}

	var weather WeatherAPIResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		c.logger.Error("Failed to parse JSON response from Weather API", "error", err)
		return nil, err
	}

	weatherDTO := client.WeatherDTO{
		Temperature: weather.Current.TempC,
		Humidity:    weather.Current.Humidity,
		Description: weather.Current.Condition.Text,
	}

	return &weatherDTO, nil
}

func (c *WeatherAPIClient) parseAPIError(body []byte) error {
	var apiErr WeatherAPIErrorResponse
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return nil
	}

	if apiErr.Error.Message != "" {
		c.logger.Error("Weather API error", "message", apiErr.Error.Message, "code", apiErr.Error.Code)

		if apiErr.Error.Code == 1006 {
			return client.ErrCityNotFound
		}
		return client.ErrInvalidRequest
	}
	return nil
}
