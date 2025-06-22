package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppURL        string `envconfig:"APP_URL" required:"true"`
	AppPort       int    `envconfig:"APP_PORT" required:"true"`
	DBHost        string `envconfig:"DB_HOST" required:"true"`
	DBHostPort    int    `envconfig:"DB_HOST_PORT" required:"true"`
	DBPort        int    `envconfig:"DB_PORT" required:"true"`
	DBUsername    string `envconfig:"DB_USERNAME" required:"true"`
	DBPassword    string `envconfig:"DB_PASSWORD" required:"true"`
	DBName        string `envconfig:"DB_NAME" required:"true"`
	WeatherAPIKey string `envconfig:"WEATHER_API_KEY" required:"true"`
	WeatherAPIUrl string `envconfig:"WEATHER_API_URL" required:"true"`
	MailEmail     string `envconfig:"MAIL_EMAIL" required:"true"`
	MailPassword  string `envconfig:"MAIL_PASSWORD" required:"true"`
}

func LoadEnvVariables() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("no .env file found or error loading it: %w", err)
	}

	var AppConfig Config

	err = envconfig.Process("", &AppConfig)
	if err != nil {
		return nil, fmt.Errorf("error processing environment variables: %w", err)
	}

	err = AppConfig.validate()
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return &AppConfig, nil
}

func (c *Config) validate() error {
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
	if c.WeatherAPIUrl == "" {
		errors = append(errors, "WEATHER_API_URL is required")
	}
	if c.MailEmail == "" {
		errors = append(errors, "MAIL_EMAIL is required")
	}
	if c.MailPassword == "" {
		errors = append(errors, "MAIL_PASSWORD is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("missing required environment variables: %v", errors)
	}

	return nil
}
