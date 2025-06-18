package db

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"gorm.io/gorm"
)

func AutomatedMigration(db *gorm.DB) {
	if err := db.AutoMigrate(&subscription.Subscription{}); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}
