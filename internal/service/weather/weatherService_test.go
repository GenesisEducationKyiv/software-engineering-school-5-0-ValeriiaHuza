//go:build unit
// +build unit

package weather

import (
	"errors"
	"testing"

	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type mockWeatherAPIClient struct {
	mock.Mock
}

func (m *mockWeatherAPIClient) FetchWeather(city string) (*client.WeatherDTO, error) {
	args := m.Called(city)
	dto, _ := args.Get(0).(*client.WeatherDTO)
	return dto, args.Error(1)
}

// --- Tests ---

func TestGetWeather_Success(t *testing.T) {
	expected := &client.WeatherDTO{
		Temperature: 17,
		Humidity:    20.5,
		Description: "Sunny",
	}

	mockClient := new(mockWeatherAPIClient)
	mockClient.On("FetchWeather", "Kyiv").Return(expected, nil)

	service := NewWeatherAPIService(mockClient)

	result, err := service.GetWeather("Kyiv")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockClient.AssertExpectations(t)
}

func TestGetWeather_Error(t *testing.T) {
	mockErr := errors.New("network error")

	mockClient := new(mockWeatherAPIClient)
	mockClient.On("FetchWeather", "London").Return((*client.WeatherDTO)(nil), mockErr)

	service := NewWeatherAPIService(mockClient)

	result, err := service.GetWeather("London")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, mockErr, err)
	mockClient.AssertExpectations(t)
}
