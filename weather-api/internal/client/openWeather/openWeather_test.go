//go:build unit
// +build unit

package openweather

import (
	"errors"
	"net/http"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"github.com/stretchr/testify/assert"
)

// --- Mocks ---

type mockGeocodingClient struct {
	coord *Coordinates
	err   error
}

func (m *mockGeocodingClient) GetCityCoordinates(city string) (*Coordinates, error) {
	return m.coord, m.err
}

// --- Tests ---

func TestFetchWeather_Success(t *testing.T) {
	weatherJSON := `{
        "main": {"temp": 22.5, "humidity": 60},
        "weather": [{"description": "clear sky"}]
    }`
	geo := &mockGeocodingClient{coord: &Coordinates{Lat: 50.0, Lon: 30.0}}
	client := newMockClient(weatherJSON, 200, nil)
	mockLog, _ := logger.NewLogger()
	api := NewWeatherAPIClient("testkey", "http://api", geo, client, mockLog)

	weather, err := api.FetchWeather("Kyiv")

	assert.NoError(t, err)
	assert.NotNil(t, weather)
	assert.Equal(t, 22.5, weather.Temperature)
	assert.Equal(t, 60.0, weather.Humidity)
	assert.Equal(t, "clear sky", weather.Description)
}

func TestFetchWeather_GeocodingError(t *testing.T) {
	geo := &mockGeocodingClient{err: errors.New("geo error")}
	mockLog, _ := logger.NewLogger()
	api := NewWeatherAPIClient("testkey", "http://api", geo, http.DefaultClient, mockLog)

	weather, err := api.FetchWeather("Kyiv")

	assert.Nil(t, weather)
	assert.Error(t, err)
}

func TestFetchWeather_Non200Status(t *testing.T) {
	geo := &mockGeocodingClient{coord: &Coordinates{Lat: 50.0, Lon: 30.0}}
	client := newMockClient("could not get weather", 404, nil)
	mockLog, _ := logger.NewLogger()
	api := NewWeatherAPIClient("testkey", "http://api", geo, client, mockLog)

	result, err := api.FetchWeather("Kyiv")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "OpenWeather API request failed with status 404: could not get weather", err.Error())
}

func TestFetchWeather_BadJSON(t *testing.T) {
	geo := &mockGeocodingClient{coord: &Coordinates{Lat: 50.0, Lon: 30.0}}
	client := newMockClient("{bad json", 200, nil)
	mockLog, _ := logger.NewLogger()
	api := NewWeatherAPIClient("testkey", "http://api", geo, client, mockLog)

	result, err := api.FetchWeather("Kyiv")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFetchWeather_EmptyWeatherArray(t *testing.T) {
	weatherJSON := `{
		"main": {"temp": 22.5, "humidity": 60},
		"weather": []
	}`
	client := newMockClient(weatherJSON, 200, nil)
	geo := &mockGeocodingClient{coord: &Coordinates{Lat: 50.0, Lon: 30.0}}

	mockLog, _ := logger.NewLogger()
	api := NewWeatherAPIClient("testkey", "http://api", geo, client, mockLog)

	result, err := api.FetchWeather("Kyiv")

	assert.Error(t, err)
	assert.Nil(t, result)
}
