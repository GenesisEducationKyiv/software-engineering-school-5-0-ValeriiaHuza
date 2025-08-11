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
)

type GeocodingClient struct {
	apiKey string
	apiUrl string
	client *http.Client
	logger logger.Logger
}

func NewGeocodingClient(apiKey string, apiUrl string, http *http.Client, logger logger.Logger) *GeocodingClient {
	return &GeocodingClient{
		apiKey: apiKey,
		apiUrl: strings.TrimRight(apiUrl, "/"),
		client: http,
		logger: logger,
	}
}

func (c *GeocodingClient) GetCityCoordinates(city string) (*Coordinates, error) {

	city = url.QueryEscape(city)

	geocodingURL := fmt.Sprintf("%s/geo/1.0/direct?q=%s&limit=1&appid=%s", c.apiUrl, city, c.apiKey)

	c.logger.Info("Sending request to OpenWeather Geocoding API", "city", city, "url", geocodingURL)
	resp, err := c.client.Get(geocodingURL)

	if err != nil {
		c.logger.Error("HTTP request to OpenWeather Geocoding failed", "error", err)

		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.Error("Failed to close response body", "error", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read geocoding response body", "error", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Geocoding API returned non-200 status code",
			"statusCode", resp.StatusCode, "body", string(body))
		return nil, errors.New("could not get city coordinates")
	}

	var geocoding []Coordinates

	if err := json.Unmarshal(body, &geocoding); err != nil {
		c.logger.Error("Failed to parse JSON response from Geocoding API", "error", err)
		return nil, err
	}

	if len(geocoding) == 0 {
		return nil, client.ErrCityNotFound
	}

	return &geocoding[0], nil
}
