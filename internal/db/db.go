package db

import (
	"fmt"
	"log"

	"github.com/ValeriiaHuza/weather_api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase(config config.Config) *gorm.DB {

	host := config.DBHost
	port := config.DBPort
	user := config.DBUsername
	password := config.DBPassword
	dbName := config.DBName

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database")

	AutomatedMigration(db)

	return db
}
