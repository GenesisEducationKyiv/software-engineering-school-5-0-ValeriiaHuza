package db

import (
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/subscription"
	"gorm.io/gorm"
)

func AutomatedMigration(db *gorm.DB) {
	if err := db.AutoMigrate(&subscription.Subscription{}); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}
