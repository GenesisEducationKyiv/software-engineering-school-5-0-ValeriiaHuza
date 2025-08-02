package openweather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net/http"
	"net/url"
	"strings"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"go.uber.org/zap"
)

type GeocodingClient struct {
	apiKey string
	apiUrl string
	client *http.Client
}

func NewGeocodingClient(apiKey string, apiUrl string, http *http.Client) *GeocodingClient {
	return &GeocodingClient{
		apiKey: apiKey,
		apiUrl: strings.TrimRight(apiUrl, "/"),
		client: http,
	}
}

func (c *GeocodingClient) GetCityCoordinates(city string) (*Coordinates, error) {

	city = url.QueryEscape(city)

	geocodingURL := fmt.Sprintf("%s/geo/1.0/direct?q=%s&limit=1&appid=%s", c.apiUrl, city, c.apiKey)

	logger.GetLogger().Info("Sending request to OpenWeather Geocoding API", zap.String("city", city), zap.String("url", geocodingURL))
	resp, err := c.client.Get(geocodingURL)

	if err != nil {
		logger.GetLogger().Error("HTTP request to OpenWeather Geocoding failed", zap.Error(err))
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.GetLogger().Error("Failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.GetLogger().Error("Failed to read geocoding response body", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logger.GetLogger().Error("Geocoding API returned non-200 status code",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("responseBody", string(body)),
		)
		return nil, errors.New("could not get city coordinates")
	}

	var geocoding []Coordinates

	if err := json.Unmarshal(body, &geocoding); err != nil {
		logger.GetLogger().Error("Failed to parse JSON response from Geocoding API", zap.Error(err))
		return nil, err
	}

	if len(geocoding) == 0 {
		return nil, client.ErrCityNotFound
	}

	return &geocoding[0], nil
}
