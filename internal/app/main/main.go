package main

import (
	"log"
	"os"

	"github.com/ValeriiaHuza/weather_api/config"
	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/ValeriiaHuza/weather_api/internal/db"
	"github.com/ValeriiaHuza/weather_api/internal/emailBuilder"
	"github.com/ValeriiaHuza/weather_api/internal/mailer"
	"github.com/ValeriiaHuza/weather_api/internal/repository"
	"github.com/ValeriiaHuza/weather_api/internal/routes"
	"github.com/ValeriiaHuza/weather_api/internal/scheduler"
	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"github.com/ValeriiaHuza/weather_api/internal/service/weather"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnvVariables()

	db := initDatabase()
	router := setupRouter()

	services := initServices(db)

	initRoutes(router, services)
	startBackgroundJobs(services.subscribeService)

	startServer(router)
}
func initDatabase() *gorm.DB {
	return db.ConnectToDatabase()
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return router
}

func initRoutes(router *gin.Engine, services *Services) {
	api := router.Group("/api")

	weatherController := weather.NewWeatherController(&services.weatherService)
	routes.WeatherRoute(api, weatherController)

	subscribeController := subscription.NewSubscribeController(&services.subscribeService)
	routes.SubscribeRoute(api, subscribeController)
}

func startBackgroundJobs(subscribeService subscription.SubscribeService) {
	schedulerService := scheduler.NewScheduler(&subscribeService)
	schedulerService.StartCronJobs()
}

func startServer(router *gin.Engine) {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initServices(database *gorm.DB) *Services {
	weatherClient := &client.WeatherAPIClient{}
	weatherService := weather.NewWeatherAPIService(weatherClient)

	subscribeRepo := repository.NewSubscriptionRepository(database)
	emailBuilder := emailBuilder.NewWeatherEmailBuilder()
	mailerService := mailer.NewMailerService(*emailBuilder)
	subscribeService := subscription.NewSubscribeService(weatherService, mailerService, subscribeRepo)

	return &Services{
		weatherService:   *weatherService,
		subscribeService: *subscribeService,
	}
}

type Services struct {
	weatherService   weather.WeatherService
	subscribeService subscription.SubscribeService
}
