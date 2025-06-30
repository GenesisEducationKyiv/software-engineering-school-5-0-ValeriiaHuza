package app

import (
	"context"
	"log"
	"strconv"

	"github.com/ValeriiaHuza/weather_api/config"
	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/ValeriiaHuza/weather_api/internal/db"
	"github.com/ValeriiaHuza/weather_api/internal/emailBuilder"
	"github.com/ValeriiaHuza/weather_api/internal/httpclient"
	"github.com/ValeriiaHuza/weather_api/internal/mailer"
	redisProvider "github.com/ValeriiaHuza/weather_api/internal/redis"
	"github.com/ValeriiaHuza/weather_api/internal/repository"
	"github.com/ValeriiaHuza/weather_api/internal/routes"
	"github.com/ValeriiaHuza/weather_api/internal/scheduler"
	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"github.com/ValeriiaHuza/weather_api/internal/service/weather"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

var ctx = context.Background()

func Run() error {
	config, err := config.LoadEnvVariables()

	if err != nil {
		return err
	}

	db := db.ConnectToDatabase(*config)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from gorm.DB: %v", err)
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	redis, err := redisProvider.ConnectToRedis(ctx, *config)

	if err != nil {
		return err
	}

	defer func() {
		if err := redis.Close(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
	}()

	router := setupRouter()

	redisPrv := redisProvider.NewRedisProvider(redis, ctx)

	services := initServices(*config, db, redisPrv)

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

func initServices(config config.Config, database *gorm.DB, redisPrv redisProvider.RedisProvider) *Services {

	http := httpclient.InitHtttClient()

	weatherClient := client.NewWeatherAPIClient(config.WeatherAPIKey, config.WeatherAPIUrl, &http)
	weatherService := weather.NewWeatherAPIService(weatherClient, &redisPrv)

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

type Services struct {
	weatherService   weather.WeatherService
	subscribeService subscription.SubscribeService
}
