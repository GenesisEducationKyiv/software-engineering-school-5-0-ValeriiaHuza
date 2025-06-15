package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ValeriiaHuza/weather_api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, user, password, dbName, port)

	if err := ConnectToDBWithRetry(dsn); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(&models.Subscription{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
}

func ConnectToDBWithRetry(dsn string) error {
	const maxRetries = 10
	for i := 1; i <= maxRetries; i++ {
		fmt.Printf("Connecting to DB (attempt %d/%d)...\n", i, maxRetries)
		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Connected to database")
			return nil
		}
		fmt.Printf("Database not ready: %v\n", err)
		time.Sleep(3 * time.Second)

		if i == maxRetries {
			return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
		}
	}
	return fmt.Errorf("failed to connect to database")
}
