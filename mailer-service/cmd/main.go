package main

import (
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/mailer-service/internal/app"
)

func main() {
	log.Println("Starting Mailer Service...")
	if err := app.Run(); err != nil {
		log.Fatalf("Mailer Service terminated with error: %v", err)
	}
}
