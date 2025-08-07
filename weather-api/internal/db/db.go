package db

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(config config.Config, logger logger.Logger) (*gorm.DB, error) {

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
