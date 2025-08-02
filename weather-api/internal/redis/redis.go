package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type RedisProvider struct {
	rdb redisClient
	ctx context.Context
}

func NewRedisProvider(redis redisClient, ctx context.Context) RedisProvider {
	return RedisProvider{
		rdb: redis,
		ctx: ctx,
	}
}

func (c *RedisProvider) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	logger.GetLogger().Info("Set to Redis", zap.String("key", key))

	return c.rdb.Set(c.ctx, key, data, 0).Err()
}

func (c *RedisProvider) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	logger.GetLogger().Info("Set to Redis with TTL", zap.String("key", key), zap.Duration("ttl", ttl))

	return c.rdb.Set(c.ctx, key, data, ttl).Err()
}

func (c *RedisProvider) Get(key string, dest interface{}) error {
	data, err := c.rdb.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}

	logger.GetLogger().Info("Get from Redis", zap.String("key", key))

	return json.Unmarshal([]byte(data), dest)
}

func (c *RedisProvider) Delete(key ...string) error {
	logger.GetLogger().Info("Delete from Redis", zap.Strings("keys", key))
	return c.rdb.Del(c.ctx, key...).Err()
}
