package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppURL        string
	AppPort       int
	DBHost        string
	DBHostPort    int
	DBPort        int
	DBUsername    string
	DBPassword    string
	DBName        string
	WeatherAPIKey string
	MailEmail     string
	MailPassword  string
}

var AppConfig *Config

func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatalf("Invalid APP_PORT: %v", err)
	}

	dbHostPort, err := strconv.Atoi(os.Getenv("DB_HOST_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_HOST_PORT: %v", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	AppConfig = &Config{
		AppURL:        os.Getenv("APP_URL"),
		AppPort:       appPort,
		DBHost:        os.Getenv("DB_HOST"),
		DBHostPort:    dbHostPort,
		DBPort:        dbPort,
		DBUsername:    os.Getenv("DB_USERNAME"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		MailEmail:     os.Getenv("MAIL_EMAIL"),
		MailPassword:  os.Getenv("MAIL_PASSWORD"),
	}

	AppConfig.validate()
}

func (c *Config) validate() {
	errors := []string{}

	if c.AppPort == 0 {
		c.AppPort = 8000 // Default port
	}

	if c.AppURL == "" {
		errors = append(errors, "APP_URL is required")
	}

	if c.DBHost == "" {
		errors = append(errors, "DB_HOST is required")
	}

	if c.DBHostPort == 0 {
		c.DBHostPort = 5432 // Default PostgreSQL port
	}

	if c.DBPort == 0 {
		c.DBPort = 5432 // Default PostgreSQL port
	}

	if c.DBUsername == "" {
		errors = append(errors, "DB_USERNAME is required")
	}
	if c.DBPassword == "" {
		errors = append(errors, "DB_PASSWORD is required")
	}
	if c.DBName == "" {
		errors = append(errors, "DB_NAME is required")
	}
	if c.WeatherAPIKey == "" {
		errors = append(errors, "WEATHER_API_KEY is required")
	}
	if c.MailEmail == "" {
		errors = append(errors, "MAIL_EMAIL is required")
	}
	if c.MailPassword == "" {
		errors = append(errors, "MAIL_PASSWORD is required")
	}

	if len(errors) > 0 {
		log.Fatalf("Missing required environment variables: %v", errors)
	}
}
