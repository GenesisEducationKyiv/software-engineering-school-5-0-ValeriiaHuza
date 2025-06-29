package main

import (
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/app"
)

func main() {
	log.Println("Starting Weather API...")
	if err := app.Run(); err != nil {
		log.Fatalf("App terminated with error: %v", err)
	}
}
