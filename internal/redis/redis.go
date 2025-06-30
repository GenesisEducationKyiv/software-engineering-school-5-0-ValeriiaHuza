package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisProvider struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisProvider(redis *redis.Client, ctx context.Context) RedisProvider {
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

	log.Printf("Set to Redis : key - %s", key)

	return c.rdb.Set(c.ctx, key, data, 0).Err()
}

func (c *RedisProvider) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	log.Printf("Set to Redis with ttl : key - %s", key)

	return c.rdb.Set(c.ctx, key, data, ttl).Err()
}

func (c *RedisProvider) Get(key string, dest interface{}) error {
	data, err := c.rdb.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}

	log.Printf("Get from Redis : key - %s", key)

	return json.Unmarshal([]byte(data), dest)
}

func (c *RedisProvider) Delete(key string) error {
	log.Printf("Delete from Redis : %s", key)
	return c.rdb.Del(c.ctx, key).Err()
}
