package db

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(config config.Config) *gorm.DB {

	dsn := config.GetDSNString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database")

	AutomatedMigration(db)

	return db
}
