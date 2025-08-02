package openweather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"go.uber.org/zap"
)

type geocodingClient interface {
	GetCityCoordinates(city string) (*Coordinates, error)
}

type WeatherAPIClient struct {
	apiKey    string
	apiUrl    string
	geocoding geocodingClient
	client    *http.Client
}

func NewWeatherAPIClient(apiKey string, apiUrl string, geocoding geocodingClient, http *http.Client) *WeatherAPIClient {
	return &WeatherAPIClient{
		apiKey:    apiKey,
		geocoding: geocoding,
		apiUrl:    strings.TrimRight(apiUrl, "/"),
		client:    http,
	}
}

func (c *WeatherAPIClient) FetchWeather(city string) (*client.WeatherDTO, error) {

	coord, err := c.geocoding.GetCityCoordinates(city)

	if err != nil {
		return nil, err
	}

	openWeatherUrl := fmt.Sprintf("%s/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric",
		c.apiUrl, coord.Lat, coord.Lon, c.apiKey)
	sanitizedUrl := strings.Replace(openWeatherUrl, c.apiKey, "[REDACTED]", 1)
	logger.GetLogger().Info("Sending request to OpenWeather API", zap.String("url", sanitizedUrl))

	resp, err := c.client.Get(openWeatherUrl)

	if err != nil {
		logger.GetLogger().Error("HTTP request to OpenWeather failed", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.GetLogger().Error("Failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Error("Failed to read OpenWeather response body", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.GetLogger().Error("OpenWeather API returned non-200 status code",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("responseBody", string(body)),
		)
		return nil, fmt.Errorf("OpenWeather API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	logger.GetLogger().Info("OpenWeather response received", zap.String("responseBody", string(body)))

	var weather OpenWeatherResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		logger.GetLogger().Error("Failed to parse JSON response from OpenWeather API", zap.Error(err))
		return nil, err
	}

	if len(weather.Weather) == 0 {
		return nil, errors.New("weather data not found")
	}

	weatherDTO := client.WeatherDTO{
		Temperature: weather.Main.Temp,
		Humidity:    weather.Main.Humidity,
		Description: weather.Weather[0].Description,
	}

	return &weatherDTO, nil
}
