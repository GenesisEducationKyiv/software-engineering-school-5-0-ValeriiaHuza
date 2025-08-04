package db

import (
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(config config.Config) (*gorm.DB, error) {

	dsn := config.GetDSNString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Connected to database")

	if err := AutomatedMigration(db); err != nil {
		log.Printf("Failed to run database migrations: %v", err)
		return nil, err
	}

	return db, nil
}
