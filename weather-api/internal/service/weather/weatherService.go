package weather

import (
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/redis"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"go.uber.org/zap"
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
}

func NewWeatherAPIService(weatherChain weatherChain, redisProvider redisProvider) *WeatherService {
	return &WeatherService{
		weatherChain:  weatherChain,
		redisProvider: redisProvider,
	}
}

func (ws *WeatherService) GetWeather(city string) (*client.WeatherDTO, error) {

	var weatherFromRedis client.WeatherDTO

	err := ws.redisProvider.Get(redis.WeatherKey+city, &weatherFromRedis)

	if err == nil {
		logger.GetLogger().Info("Weather data retrieved from Redis", zap.String("city", city))
		return &weatherFromRedis, nil
	}

	// Log Redis errors (not cache misses)
	if err.Error() != "redis: nil" { // or use redis.Nil constant if available
		logger.GetLogger().Error("Redis error while fetching weather", zap.String("city", city), zap.Error(err))

	}

	weatherDto, err := ws.weatherChain.GetWeather(city)
	if err != nil {
		logger.GetLogger().Error("Failed to get weather from chain", zap.String("city", city), zap.Error(err))

		return nil, err
	}

	err = ws.redisProvider.SetWithTTL(redis.WeatherKey+city, weatherDto, redis.WeatherTTL)

	if err != nil {
		logger.GetLogger().Error("Failed to save weather in Redis", zap.String("city", city), zap.Error(err))
	}

	return weatherDto, nil
}
