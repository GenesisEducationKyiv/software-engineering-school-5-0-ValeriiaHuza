package app

import (
	"log"
	"strconv"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client"
	openweather "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client/openWeather"
	weatherapi "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client/weatherApi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/emailBuilder"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/httpclient"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/mailer"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/repository"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/routes"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/scheduler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/weather"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/logger"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func Run() error {

	err := logger.InitLoggerFile("app.log")
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.CloseLogFile()

	config, err := config.LoadEnvVariables()

	if err != nil {
		return err
	}

	db, err := db.ConnectToDatabase(*config)

	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get sql.DB from gorm.DB: %v", err)
		return err
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	router := setupRouter()

	services := initServices(*config, db)

	initRoutes(router, services)
	startBackgroundJobs(services.subscribeService)

	return startServer(*config, router)
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

func startServer(config config.Config, router *gin.Engine) error {
	port := strconv.Itoa(config.AppPort)

	return router.Run(":" + port)
}

func initServices(config config.Config, database *gorm.DB) *Services {

	weatherApiChain := buildWeatherResponsibilityChain(config)

	weatherService := weather.NewWeatherAPIService(weatherApiChain)

	subscribeRepo := repository.NewSubscriptionRepository(database)
	emailBuilder := emailBuilder.NewWeatherEmailBuilder(config.AppURL)

	mailEmail := config.MailEmail
	dialer := gomail.NewDialer("smtp.gmail.com", 587, mailEmail, config.MailPassword)
	mailerService := mailer.NewMailerService(mailEmail, dialer, emailBuilder)

	subscribeService := subscription.NewSubscribeService(weatherService, mailerService, subscribeRepo)

	return &Services{
		weatherService:   *weatherService,
		subscribeService: *subscribeService,
	}
}

func buildWeatherResponsibilityChain(config config.Config) *client.WeatherChain {
	http := httpclient.InitHtttClient()

	geocoding := openweather.NewGeocodingClient(config.OpenWeatherKey, config.OpenWeatherUrl, &http)

	weatherApiClient := weatherapi.NewWeatherAPIClient(config.WeatherAPIKey, config.WeatherAPIUrl, &http)
	openWeatherClient := openweather.NewWeatherAPIClient(config.OpenWeatherKey, config.OpenWeatherUrl, geocoding, &http)

	weatherApiChain := client.NewWeatherChain(weatherApiClient)
	openWeatherChain := client.NewWeatherChain(openWeatherClient)

	weatherApiChain.SetNext(openWeatherChain)

	return weatherApiChain
}

type Services struct {
	weatherService   weather.WeatherService
	subscribeService subscription.SubscribeService
}
