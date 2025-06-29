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

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client"
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

	log.Printf("Sending request to: %s", geocodingURL)
	resp, err := c.client.Get(geocodingURL)

	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read geocoding response body:", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Geogoding error : %s", body)
		return nil, errors.New("could not get city coordinates")
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
