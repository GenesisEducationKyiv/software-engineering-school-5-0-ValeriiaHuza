package weather

import (
	"log"
	"time"

	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/ValeriiaHuza/weather_api/internal/redis"
)

type weatherAPIClient interface {
	FetchWeather(city string) (*client.WeatherDTO, error)
}

type redisProvider interface {
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	Get(key string, dest interface{}) error
}

type WeatherService struct {
	weatherClient weatherAPIClient
	redisProvider redisProvider
}

func NewWeatherAPIService(weatherClient weatherAPIClient, redisProvider redisProvider) *WeatherService {
	return &WeatherService{
		weatherClient: weatherClient,
		redisProvider: redisProvider,
	}
}

func (ws *WeatherService) GetWeather(city string) (*client.WeatherDTO, error) {

	var weatherFromRedis client.WeatherDTO

	err := ws.redisProvider.Get(redis.WeatherKey+city, &weatherFromRedis)

	if err == nil {
		log.Printf("Get weather for %s from Redis : %+v", city, weatherFromRedis)
		return &weatherFromRedis, nil
	}

	// Log Redis errors (not cache misses)
	if err.Error() != "redis: nil" { // or use redis.Nil constant if available
		log.Printf("Redis error while fetching weather for %s: %v", city, err)
	}

	weatherDto, err := ws.weatherClient.FetchWeather(city)
	if err != nil {
		log.Println("HTTP error in GetWeather :", err)
		return nil, err
	}

	err = ws.redisProvider.SetWithTTL(redis.WeatherKey+city, weatherDto, redis.WeatherTTL)

	if err != nil {
		log.Printf("Failed to save weather in redis %v", err.Error())
	}

	return weatherDto, nil
}
