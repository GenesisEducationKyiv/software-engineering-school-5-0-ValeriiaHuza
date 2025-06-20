//go:build unit
// +build unit

package weather

import (
	"errors"
	"testing"

	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/stretchr/testify/assert"
)

type mockWeatherAPIClient struct {
	fetchWeatherFunc func(city string) (*client.WeatherDTO, error)
}

func (m *mockWeatherAPIClient) FetchWeather(city string) (*client.WeatherDTO, error) {
	return m.fetchWeatherFunc(city)
}

func TestGetWeather_Success(t *testing.T) {
	expected := &client.WeatherDTO{
		Temperature: 17,
		Humidity:    20.5,
		Description: "Sunny",
	}
	mockClient := &mockWeatherAPIClient{
		fetchWeatherFunc: func(city string) (*client.WeatherDTO, error) {
			assert.Equal(t, "Kyiv", city)
			return expected, nil
		},
	}
	service := NewWeatherAPIService(mockClient)

	result, err := service.GetWeather("Kyiv")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetWeather_Error(t *testing.T) {
	mockErr := errors.New("network error")
	mockClient := &mockWeatherAPIClient{
		fetchWeatherFunc: func(city string) (*client.WeatherDTO, error) {
			return nil, mockErr
		},
	}
	service := NewWeatherAPIService(mockClient)

	result, err := service.GetWeather("London")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, mockErr, err)
}
