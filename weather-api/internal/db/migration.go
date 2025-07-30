package db

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/service/subscription"
	"gorm.io/gorm"
)

func AutomatedMigration(db *gorm.DB) error {
	return db.AutoMigrate(&subscription.Subscription{})
}
