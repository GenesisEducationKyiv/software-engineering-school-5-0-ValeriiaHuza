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
	"gorm.io/gorm"
)

func main() {
	config.LoadEnvVariables()

	db := initDatabase()
	router := setupRouter()

	services := initServices(db)
	initRoutes(router, services)
	startBackgroundJobs(services.SubscribeService)

	startServer(router)
}

func initDatabase() *gorm.DB {
	db.ConnectToDatabase()
	return db.DB
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

	weatherController := controller.NewWeatherController(services.WeatherService)
	routes.WeatherRoute(api, weatherController)

	subscribeController := controller.NewSubscribeController(services.SubscribeService)
	routes.SubscribeRoute(api, subscribeController)
}

func startBackgroundJobs(subscribeService service.SubscribeService) {
	schedulerService := scheduler.NewScheduler(subscribeService)
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
	weatherClient := &utils.WeatherAPIClientImpl{}
	weatherService := service.NewWeatherAPIService(weatherClient)

	subscribeRepo := repository.NewSubscriptionRepository(database)
	emailBuilder := utils.NewWeatherEmailBuilder()
	mailerService := service.NewMailerService(*emailBuilder)
	subscribeService := service.NewSubscribeService(weatherService, mailerService, subscribeRepo)

	return &Services{
		WeatherService:   weatherService,
		SubscribeService: subscribeService,
	}
}

type Services struct {
	WeatherService   service.WeatherService
	SubscribeService service.SubscribeService
}
