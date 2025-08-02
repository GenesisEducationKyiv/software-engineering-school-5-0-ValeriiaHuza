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
	"go.uber.org/zap"
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

	logger.GetLogger().Info("Sending request to Weather API", zap.String("city", city), zap.String("url", weatherURL))

	resp, err := c.client.Get(weatherURL)

	if err != nil {
		logger.GetLogger().Error("HTTP request to Weather API failed", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.GetLogger().Error("Failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Error("Failed to read response body", zap.Error(err))
		return nil, err
	}

	if apiErr := c.parseAPIError(body); apiErr != nil {
		return nil, apiErr
	}

	var weather WeatherAPIResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		logger.GetLogger().Error("Failed to parse JSON response from Weather API", zap.Error(err))
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
		logger.GetLogger().Error("Weather API error", zap.String("message", apiErr.Error.Message), zap.Int("code", apiErr.Error.Code))

		if apiErr.Error.Code == 1006 {
			return client.ErrCityNotFound
		}
		return client.ErrInvalidRequest
	}
	return nil
}
