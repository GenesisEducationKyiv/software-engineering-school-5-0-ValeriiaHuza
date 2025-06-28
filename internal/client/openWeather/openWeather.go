package openweather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ValeriiaHuza/weather_api/internal/client"
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

	openWeatherUrl := fmt.Sprintf("%s/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", c.apiUrl, coord.Lat, coord.Lon, c.apiKey)

	log.Printf("Sending request to OpenWeather API: %s", openWeatherUrl)

	resp, err := c.client.Get(openWeatherUrl)

	if err != nil {
		log.Printf("HTTP request to OpenWeather failed: %v", err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read open weather response body:", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Open weather error : %s", body)
		return nil, errors.New("could not get weather")
	}

	log.Printf("Open weather response : %s", string(body))

	var weather OpenWeatherResponse

	if err := json.Unmarshal(body, &weather); err != nil {
		log.Println("Failed to parse JSON:", err)
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
