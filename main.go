package main

import (
	"log"
	"os"

	"github.com/ValeriiaHuza/weather_api/config"
	"github.com/ValeriiaHuza/weather_api/controller"
	"github.com/ValeriiaHuza/weather_api/db"
	"github.com/ValeriiaHuza/weather_api/repository"
	"github.com/ValeriiaHuza/weather_api/routes"
	"github.com/ValeriiaHuza/weather_api/scheduler"
	"github.com/ValeriiaHuza/weather_api/service"
	"github.com/ValeriiaHuza/weather_api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnvVariables()

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

	serviceWeather := service.NewWeatherAPIService(&utils.WeatherAPIClientImpl{})
	controllerWeather := controller.NewWeatherController(serviceWeather)
	routes.WeatherRoute(api, controllerWeather)

	// Subscribe repository
	subscribeRepository := repository.NewSubscriptionRepository()

	//Email builder
	emailBuilder := utils.NewWeatherEmailBuilder()

	//MailerService
	mailerService := service.NewMailerService(*emailBuilder)

	// Subscription service and controller
	serviceSubscription := service.NewSubscribeService(serviceWeather, mailerService, subscribeRepository)
	controllerSubscription := controller.NewSubscribeController(serviceSubscription)
	routes.SubscribeRoute(api, controllerSubscription)

	// Start background jobs
	schedulerS := scheduler.NewScheduler(serviceSubscription)
	schedulerS.StartCronJobs()

	// Start the server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000" // default fallback
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
