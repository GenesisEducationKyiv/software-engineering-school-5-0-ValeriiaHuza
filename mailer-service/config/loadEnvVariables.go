package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MailerPort int    `envconfig:"MAILER_PORT"`
	MailerURL  string `envconfig:"MAILER_URL" required:"true"`

	ApiURL string `envconfig:"API_URL" required:"true"`

	MailEmail    string `envconfig:"MAIL_EMAIL" required:"true"`
	MailPassword string `envconfig:"MAIL_PASSWORD" required:"true"`

	RabbitMQUrl string `envconfig:"RABBITMQ_URL" required:"true"`
	MQUsername  string `envconfig:"MQ_USERNAME" required:"true"`
	MQPassword  string `envconfig:"MQ_PASSWORD" required:"true"`
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

	if c.MailerPort == 0 {
		c.MailerPort = 8002 // Default port
	}

	if c.MailerURL == "" {
		errors = append(errors, "MAILER_URL is required")
	}

	if c.ApiURL == "" {
		errors = append(errors, "API_URL is required")
	}

	if c.MailEmail == "" {
		errors = append(errors, "MAIL_EMAIL is required")
	}
	if c.MailPassword == "" {
		errors = append(errors, "MAIL_PASSWORD is required")
	}

	if c.RabbitMQUrl == "" {
		errors = append(errors, "RABBITMQ_URL is required")
	}
	if c.MQUsername == "" {
		errors = append(errors, "MQ_USERNAME is required")
	}
	if c.MQPassword == "" {
		errors = append(errors, "MQ_PASSWORD is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("missing required environment variables: %v", errors)
	}

	return nil
}
