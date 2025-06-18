package main

import (
	"context"
	"log"

	"github.com/ValeriiaHuza/weather_api/internal/app"
)

func main() {
	ctx := context.Background()

	if err := app.Run(ctx); err != nil {
		log.Fatalf("App terminated with error: %v", err)
	}
}
