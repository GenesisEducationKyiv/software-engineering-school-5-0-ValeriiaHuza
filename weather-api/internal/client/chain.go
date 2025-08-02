package client

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"go.uber.org/zap"
)

type weatherProvider interface {
	FetchWeather(city string) (*WeatherDTO, error)
}

type weatherChainProvider interface {
	GetWeather(city string) (*WeatherDTO, error)
	SetNext(next weatherChainProvider)
}

type WeatherChain struct {
	next     weatherChainProvider
	provider weatherProvider
}

func (h *WeatherChain) SetNext(provider weatherChainProvider) {
	h.next = provider
}

func NewWeatherChain(provider weatherProvider) *WeatherChain {
	return &WeatherChain{
		provider: provider,
	}
}

func (c *WeatherChain) GetWeather(city string) (*WeatherDTO, error) {
	weather, err := c.provider.FetchWeather(city)
	if err == nil {
		return weather, nil
	}

	logger.GetLogger().Error("Weather provider error. Trying next provider...", zap.String("city", city), zap.Error(err))

	if c.next != nil {
		return c.next.GetWeather(city)
	}

	logger.GetLogger().Error("All weather providers failed", zap.String("city", city), zap.Error(err))

	return nil, err
}
