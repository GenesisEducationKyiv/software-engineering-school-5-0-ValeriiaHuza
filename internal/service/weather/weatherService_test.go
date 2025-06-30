//go:build unit
// +build unit

package weather

import (
	"errors"
	"testing"
	"time"

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

type mockRedisProvider struct {
	mock.Mock
}

func (m *mockRedisProvider) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *mockRedisProvider) Get(key string, dest interface{}) error {
	args := m.Called(key, dest)

	if dto, ok := args.Get(1).(*client.WeatherDTO); ok {
		if out, ok := dest.(*client.WeatherDTO); ok {
			*out = *dto
		}
	}

	return args.Error(0)
}

// --- Tests ---

func TestGetWeather_Success(t *testing.T) {
	expected := &client.WeatherDTO{
		Temperature: 17,
		Humidity:    20.5,
		Description: "Sunny",
	}

	mockRedis := new(mockRedisProvider)
	mockClient := new(mockWeatherAPIClient)

	mockRedis.On("Get", "weather:Kyiv", mock.Anything).
		Return(nil, expected)

	service := NewWeatherAPIService(mockClient, mockRedis)

	result, err := service.GetWeather("Kyiv")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRedis.AssertCalled(t, "Get", mock.Anything, mock.Anything)
	mockClient.AssertExpectations(t)
}

func TestGetWeather_CacheMiss_Success(t *testing.T) {
	mockRedis := new(mockRedisProvider)
	mockClient := new(mockWeatherAPIClient)
	service := NewWeatherAPIService(mockClient, mockRedis)

	city := "Lviv"
	expected := &client.WeatherDTO{
		Temperature: 15.0,
		Humidity:    60,
		Description: "Cloudy",
	}

	mockRedis.On("Get", "weather:"+city, mock.Anything).Return(errors.New("not found"), nil)

	mockClient.On("FetchWeather", city).Return(expected, nil)

	mockRedis.On("SetWithTTL", "weather:"+city, expected, mock.Anything).Return(nil)

	result, err := service.GetWeather(city)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRedis.AssertCalled(t, "Get", "weather:"+city, mock.Anything)
	mockClient.AssertCalled(t, "FetchWeather", city)
	mockRedis.AssertCalled(t, "SetWithTTL", "weather:"+city, expected, mock.Anything)
}

func TestGetWeather_CacheMiss_APIError(t *testing.T) {
	mockRedis := new(mockRedisProvider)
	mockClient := new(mockWeatherAPIClient)
	service := NewWeatherAPIService(mockClient, mockRedis)

	city := "Odesa"
	mockRedis.On("Get", mock.Anything, mock.Anything).Return(errors.New("not found"), nil)
	mockClient.On("FetchWeather", city).Return(nil, errors.New("api error"))

	result, err := service.GetWeather(city)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockRedis.AssertCalled(t, "Get", mock.Anything, mock.Anything)
	mockClient.AssertCalled(t, "FetchWeather", city)
}

func TestGetWeather_CacheMiss_APISuccess_RedisSetError(t *testing.T) {
	mockRedis := new(mockRedisProvider)
	mockClient := new(mockWeatherAPIClient)
	service := NewWeatherAPIService(mockClient, mockRedis)

	city := "Dnipro"
	expected := &client.WeatherDTO{Temperature: 10.0, Humidity: 70, Description: "Rainy"}

	mockRedis.On("Get", mock.Anything, mock.Anything).Return(errors.New("not found"), nil)
	mockClient.On("FetchWeather", city).Return(expected, nil)
	mockRedis.On("SetWithTTL", mock.Anything, expected, mock.Anything).Return(errors.New("redis error"))

	result, err := service.GetWeather(city)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRedis.AssertCalled(t, "SetWithTTL", mock.Anything, expected, mock.Anything)
}
