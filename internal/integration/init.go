//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"time"

	redisClient "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupPostgresContainer() (*gorm.DB, func(), error) {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to map port: %w", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get host: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port.Port(), dbUser, dbPassword, dbName)

	// Wait for DB to be ready
	var gormDB *gorm.DB
	for i := 0; i < 10; i++ {
		gormDB, err = gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	terminate := func() {
		_ = postgresContainer.Terminate(ctx)
	}

	return gormDB, terminate, nil
}

func SetupRedisContainer() (*redisClient.Client, func(), error) {
	ctx := context.Background()

	// Start Redis container
	redisContainer, err := redis.Run(ctx,
		"redis:7-alpine",
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	// Get host and port
	host, err := redisContainer.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get redis host: %w", err)
	}

	port, err := redisContainer.MappedPort(ctx, "6379/tcp")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get redis port: %w", err)
	}

	// Create redis client
	rdb := redisClient.NewClient(&redisClient.Options{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
		DB:   0,
	})

	// Ping until available
	for i := 0; i < 10; i++ {
		_, err = rdb.Ping(ctx).Result()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	terminate := func() {
		_ = redisContainer.Terminate(ctx)
	}

	return rdb, terminate, nil
}
