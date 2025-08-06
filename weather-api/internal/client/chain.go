package client

type loggerInterface interface {
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

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
	logger   loggerInterface
}

func (h *WeatherChain) SetNext(provider weatherChainProvider) {
	h.next = provider
}

func NewWeatherChain(provider weatherProvider, logger loggerInterface) *WeatherChain {
	return &WeatherChain{
		provider: provider,
		logger:   logger,
	}
}

func (c *WeatherChain) GetWeather(city string) (*WeatherDTO, error) {
	weather, err := c.provider.FetchWeather(city)
	if err == nil {
		return weather, nil
	}

	c.logger.Error("Weather provider error. Trying next provider... ", "city", city, "error", err)

	if c.next != nil {
		return c.next.GetWeather(city)
	}

	c.logger.Error("All weather providers failed", "city", city, "error", err)

	return nil, err
}
