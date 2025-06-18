package app

import (
	"strconv"

	"github.com/ValeriiaHuza/weather_api/config"
	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/ValeriiaHuza/weather_api/internal/db"
	"github.com/ValeriiaHuza/weather_api/internal/emailBuilder"
	"github.com/ValeriiaHuza/weather_api/internal/httpclient"
	"github.com/ValeriiaHuza/weather_api/internal/mailer"
	"github.com/ValeriiaHuza/weather_api/internal/repository"
	"github.com/ValeriiaHuza/weather_api/internal/routes"
	"github.com/ValeriiaHuza/weather_api/internal/scheduler"
	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"github.com/ValeriiaHuza/weather_api/internal/service/weather"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func Run() error {
	config.LoadEnvVariables()

	db := initDatabase()
	router := setupRouter()

	services := initServices(db)

	initRoutes(router, services)
	startBackgroundJobs(services.subscribeService)

	return startServer(router)
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

func startServer(router *gin.Engine) error {
	port := strconv.Itoa(config.AppConfig.AppPort)

	return router.Run(":" + port)
}

func initServices(database *gorm.DB) *Services {

	http := httpclient.InitHtttClient()

	weatherClient := client.NewWeatherAPIClient(config.AppConfig.WeatherAPIKey, &http)
	weatherService := weather.NewWeatherAPIService(weatherClient)

	subscribeRepo := repository.NewSubscriptionRepository(database)
	emailBuilder := emailBuilder.NewWeatherEmailBuilder(config.AppConfig.AppURL)

	mailEmail := config.AppConfig.MailEmail
	dialer := gomail.NewDialer("smtp.gmail.com", 587, mailEmail, config.AppConfig.MailPassword)
	mailerService := mailer.NewMailerService(mailEmail, dialer, emailBuilder)

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
