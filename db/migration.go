package db

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/models"
)

func AutomatedMigration() {
	if err := DB.AutoMigrate(&models.Subscription{}); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}
