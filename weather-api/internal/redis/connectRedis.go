package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/config"
	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(ctx context.Context, config config.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})

	// Ping to check Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		redisClient.Close()
		return nil, err
	}

	log.Println("Connected to redis")

	return redisClient, nil
}
