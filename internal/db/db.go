package db

import (
	"fmt"
	"log"

	"github.com/ValeriiaHuza/weather_api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {

	host := config.AppConfig.DBHost
	port := config.AppConfig.DBPort
	user := config.AppConfig.DBUsername
	password := config.AppConfig.DBPassword
	dbName := config.AppConfig.DBName

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database")

	AutomatedMigration(db)

	return db
}
