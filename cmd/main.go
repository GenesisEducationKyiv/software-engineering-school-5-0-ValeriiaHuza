package main

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/internal/app"
)

func main() {
	log.Println("Starting Weather API...")
	if err := app.Run(); err != nil {
		log.Fatalf("App terminated with error: %v", err)
	}
}
