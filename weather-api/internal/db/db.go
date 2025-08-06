package db

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type loggerInterface interface {
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

func ConnectToDatabase(config config.Config, logger loggerInterface) (*gorm.DB, error) {

	dsn := config.GetDSNString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return nil, err
	}

	logger.Info("Connected to database", "dsn", dsn)

	if err := AutomatedMigration(db); err != nil {
		logger.Error("Failed to run database migrations", "error", err)
		return nil, err
	}

	return db, nil
}
