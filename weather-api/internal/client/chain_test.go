//go:build unit
// +build unit

package client

import (
	"errors"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWeatherProvider struct {
	mock.Mock
}

func (m *mockWeatherProvider) FetchWeather(city string) (*WeatherDTO, error) {
	args := m.Called(city)
	dto, _ := args.Get(0).(*WeatherDTO)
	return dto, args.Error(1)
}

func TestWeatherChain_SuccessFirstProvider_Mock(t *testing.T) {

	want := &WeatherDTO{Temperature: 25}
	provider := new(mockWeatherProvider)
	provider.On("FetchWeather", "Kyiv").Return(want, nil)

	mockLog, _ := logger.NewTestLogger()
	chain := NewWeatherChain(provider, *mockLog)

	got, err := chain.GetWeather("Kyiv")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
	provider.AssertExpectations(t)
}

func TestWeatherChain_SecondProviderSuccess_Mock(t *testing.T) {

	provider1 := new(mockWeatherProvider)
	provider2 := new(mockWeatherProvider)
	want := &WeatherDTO{Temperature: 18}

	provider1.On("FetchWeather", "Lviv").Return(nil, errors.New("fail1"))
	provider2.On("FetchWeather", "Lviv").Return(want, nil)

	mockLog, _ := logger.NewTestLogger()
	chain := NewWeatherChain(provider1, *mockLog)
	chain.SetNext(NewWeatherChain(provider2, *mockLog))

	got, err := chain.GetWeather("Lviv")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
	provider1.AssertExpectations(t)
	provider2.AssertExpectations(t)
}

func TestWeatherChain_AllProvidersFail_Mock(t *testing.T) {

	provider1 := new(mockWeatherProvider)
	provider2 := new(mockWeatherProvider)

	provider1.On("FetchWeather", "Odesa").Return(nil, errors.New("fail1"))
	provider2.On("FetchWeather", "Odesa").Return(nil, errors.New("fail2"))

	mockLog, _ := logger.NewTestLogger()
	chain := NewWeatherChain(provider1, *mockLog)
	chain.SetNext(NewWeatherChain(provider2, *mockLog))

	got, err := chain.GetWeather("Odesa")
	assert.Nil(t, got)
	assert.Error(t, err)

	assert.Equal(t, errors.New("fail2").Error(), err.Error())
	provider1.AssertExpectations(t)
	provider2.AssertExpectations(t)
}
