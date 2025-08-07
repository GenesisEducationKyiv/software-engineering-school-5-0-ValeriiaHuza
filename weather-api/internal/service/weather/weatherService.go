package weather

import (
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/redis"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
)

type weatherChain interface {
	GetWeather(city string) (*client.WeatherDTO, error)
}

type redisProvider interface {
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	Get(key string, dest interface{}) error
}

type WeatherService struct {
	weatherChain  weatherChain
	redisProvider redisProvider
	logger        logger.Logger
}

func NewWeatherAPIService(weatherChain weatherChain, redisProvider redisProvider,
	logger logger.Logger) *WeatherService {
	return &WeatherService{
		weatherChain:  weatherChain,
		redisProvider: redisProvider,
		logger:        logger,
	}
}

func (ws *WeatherService) GetWeather(city string) (*client.WeatherDTO, error) {

	var weatherFromRedis client.WeatherDTO

	err := ws.redisProvider.Get(redis.WeatherKey+city, &weatherFromRedis)

	if err == nil {
		ws.logger.Info("Weather data retrieved from Redis", "city", city)
		return &weatherFromRedis, nil
	}

	// Log Redis errors (not cache misses)
	if err.Error() != "redis: nil" { // or use redis.Nil constant if available
		ws.logger.Error("Failed to get weather from Redis", "city", city, "error", err)

	}

	weatherDto, err := ws.weatherChain.GetWeather(city)
	if err != nil {
		ws.logger.Error("Failed to get weather from chain", "city", city, "error", err)

		return nil, err
	}

	err = ws.redisProvider.SetWithTTL(redis.WeatherKey+city, weatherDto, redis.WeatherTTL)

	if err != nil {
		ws.logger.Error("Failed to save weather in Redis", "city", city, "error", err)
	}

	return weatherDto, nil
}
