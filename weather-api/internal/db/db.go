package db

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(config config.Config) (*gorm.DB, error) {

	dsn := config.GetDSNString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.GetLogger().Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	logger.GetLogger().Info("Connected to database", zap.String("dsn", dsn))

	if err := AutomatedMigration(db); err != nil {
		logger.GetLogger().Error("Failed to run database migrations", zap.Error(err))
		return nil, err
	}

	return db, nil
}
