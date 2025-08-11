package redis

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(ctx context.Context, config config.Config, logger logger.Logger) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})

	// Ping to check Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		redisClient.Close()
		return nil, err
	}

	logger.Info("Connected to Redis", "host", config.RedisHost, "port", config.RedisPort)

	return redisClient, nil
}
