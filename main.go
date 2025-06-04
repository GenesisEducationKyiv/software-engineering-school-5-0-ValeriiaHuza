package main

import (
	"log"
	"os"

	"github.com/ValeriiaHuza/weather_api/config"
	"github.com/ValeriiaHuza/weather_api/controller"
	"github.com/ValeriiaHuza/weather_api/db"
	"github.com/ValeriiaHuza/weather_api/models"
	"github.com/ValeriiaHuza/weather_api/routes"
	"github.com/ValeriiaHuza/weather_api/service"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func init() {
	config.LoadEnvVariables()
}

func main() {
	router := gin.Default()

	// Connect to the database
	db.ConnectToDatabase()

	// Serve static files (e.g., subscribe.html)
	router.Static("/static", "./static")

	// Route for the HTML page
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// API routes group
	api := router.Group("/api")

	// Weather service and controller
	serviceWeather := service.NewWeatherService()
	controllerWeather := controller.NewWeatherController(serviceWeather)
	routes.WeatherRoute(api, controllerWeather)

	// Subscription service and controller
	serviceSubscription := service.NewSubscribeService(serviceWeather)
	controllerSubscription := controller.NewSubscribeController(serviceSubscription)
	routes.SubscribeRoute(api, controllerSubscription)

	// Start background jobs
	startCronJobs()

	// Start the server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000" // default fallback
	}
	router.Run(":" + port)
}

func startCronJobs() {
	c := cron.New()

	if err := c.AddFunc("0 0 9 * * *", func() {
		service.SendEmails(models.FrequencyDaily)
	}); err != nil {
		log.Println("Failed to schedule daily job:", err)
	}

	// Every hour
	if err := c.AddFunc("0 0 * * * *", func() {
		service.SendEmails(models.FrequencyHourly)
	}); err != nil {
		log.Println("Failed to schedule hourly job:", err)
	}

	c.Start()
}
