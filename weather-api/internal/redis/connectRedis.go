package redis

import (
	"context"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ConnectToRedis(ctx context.Context, config config.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})

	// Ping to check Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		logger.GetLogger().Error("Failed to connect to Redis", zap.Error(err))
		redisClient.Close()
		return nil, err
	}

	logger.GetLogger().Info("Connected to Redis", zap.String("host", config.RedisHost), zap.Int("port", config.RedisPort))

	return redisClient, nil
}
