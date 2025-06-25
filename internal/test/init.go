package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupPostgresContainer(t *testing.T) (db *gorm.DB, terminate func()) {
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

	assert.NoError(t, err)

	port, err := postgresContainer.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	host, err := postgresContainer.Host(ctx)
	assert.NoError(t, err)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port.Port(), dbUser, dbPassword, dbName)
	// Wait for DB to be ready
	var gormDB *gorm.DB
	for i := 0; i < 10; i++ {
		gormDB, err = gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	assert.NoError(t, err)

	terminate = func() {
		_ = postgresContainer.Terminate(ctx)
	}

	return gormDB, terminate
}
