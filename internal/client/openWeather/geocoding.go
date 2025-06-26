package openweather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ValeriiaHuza/weather_api/internal/client"
)

type GeocodingClient struct {
	apiKey string
	apiUrl string
	client *http.Client
}

func NewGeocodingClient(apiKey string, apiUrl string, http *http.Client) *GeocodingClient {
	return &GeocodingClient{
		apiKey: apiKey,
		apiUrl: apiUrl,
		client: http,
	}
}

func (c *GeocodingClient) GetCityCoordinates(city string) (*Coordinates, error) {

	city = url.QueryEscape(city)

	geocodingUrl := fmt.Sprintf("%s/geo/1.0/direct/?&q=%s&limit=1&appid=%s", strings.TrimRight(c.apiUrl, "/"), city, c.apiKey)

	resp, err := c.client.Get(geocodingUrl)

	if err != nil {
		log.Println("Geocodong HTTP request failed:", err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("could not het city coordinates")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read geocoding response body:", err)
		return nil, err
	}

	var geocoding []Coordinates

	if err := json.Unmarshal(body, &geocoding); err != nil {
		log.Println("Failed to parse JSON:", err)
		return nil, err
	}

	if len(geocoding) == 0 {
		return nil, client.ErrCityNotFound
	}

	return &geocoding[0], nil
}
