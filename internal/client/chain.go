package client

import (
	"log"
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

	log.Printf("Weather provider error for city '%s': %s. Trying next provider...", city, err)

	if c.next != nil {
		return c.next.GetWeather(city)
	}

	log.Printf("all providers failed for city '%s': last error: %s", city, err)

	return nil, err
}
