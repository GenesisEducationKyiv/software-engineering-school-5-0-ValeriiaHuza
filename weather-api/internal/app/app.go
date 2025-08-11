package app

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/config"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	openweather "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client/openWeather"
	weatherapi "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client/weatherApi"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/httpclient"
	metricP "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/metrics"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/repository"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/routes"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/scheduler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/service/weather"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	redisProvider "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/redis"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Run() error {
	var ctx = context.Background()

	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Error syncing logger: %v\n", err)
		}
	}()

	logger.Info("Starting Weather Api Service...")

	config, err := config.LoadEnvVariables()

	if err != nil {
		return err
	}

	db, err := db.ConnectToDatabase(*config, *logger)

	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get sql.DB from gorm.DB", "error", err)
		return err
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database connection", "error", err)

		}
	}()

	redis, err := redisProvider.ConnectToRedis(ctx, *config, *logger)

	if err != nil {
		return err
	}

	defer func() {
		if err := redis.Close(); err != nil {
			logger.Error("Failed to close Redis connection", "error", err)
		}
	}()

	rabbit, err := rabbitmq.ConnectToRabbitMQ(config.RabbitMQUrl, *logger)
	if err != nil {
		return err
	}
	defer rabbit.Conn.Close()
	defer rabbit.Channel.Close()

	if err := declareQueues(rabbit); err != nil {
		return err
	}

	emailPublisher := rabbitmq.NewRabbitMQPublisher(rabbit.Channel)

	router := setupRouter()

	redisPrv := redisProvider.NewRedisProvider(redis, ctx, *logger)

	services := initServices(*config, db, redisPrv, emailPublisher, *logger)

	initRoutes(router, services)
	startBackgroundJobs(*services.subscribeService, *logger)

	return startServer(*config, router)
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(metricP.MetricsMiddleware())

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}

func initRoutes(router *gin.Engine, services *Services) {
	api := router.Group("/api")

	weatherController := weather.NewWeatherController(services.weatherService)
	routes.WeatherRoute(api, weatherController)

	subscribeController := subscription.NewSubscribeController(services.subscribeService)
	routes.SubscribeRoute(api, subscribeController)
}

func startBackgroundJobs(subscribeService subscription.SubscribeService, logger logger.Logger) {
	schedulerService := scheduler.NewScheduler(&subscribeService, logger)
	schedulerService.StartCronJobs()
}

func startServer(config config.Config, router *gin.Engine) error {
	port := strconv.Itoa(config.AppPort)

	return router.Run(":" + port)
}

func initServices(config config.Config, database *gorm.DB,
	redisPrv redisProvider.RedisProvider, emailPublisher *rabbitmq.RabbitMQPublisher,
	logger logger.Logger) *Services {

	weatherApiChain := buildWeatherResponsibilityChain(config, logger)

	weatherService := weather.NewWeatherAPIService(weatherApiChain, &redisPrv, logger)

	subscribeRepo := repository.NewSubscriptionRepository(database)

	subscribeService := subscription.NewSubscribeService(weatherService, subscribeRepo, emailPublisher, logger)

	return &Services{
		weatherService:   weatherService,
		subscribeService: subscribeService,
	}
}

func buildWeatherResponsibilityChain(config config.Config, logger logger.Logger) *client.WeatherChain {
	http := httpclient.InitHttpClient()

	geocoding := openweather.NewGeocodingClient(config.OpenWeatherKey, config.OpenWeatherUrl, &http, logger)

	weatherApiClient := weatherapi.NewWeatherAPIClient(config.WeatherAPIKey,
		config.WeatherAPIUrl, &http, logger)
	openWeatherClient := openweather.NewWeatherAPIClient(config.OpenWeatherKey,
		config.OpenWeatherUrl, geocoding, &http, logger)

	weatherApiChain := client.NewWeatherChain(weatherApiClient, logger)
	openWeatherChain := client.NewWeatherChain(openWeatherClient, logger)

	weatherApiChain.SetNext(openWeatherChain)

	return weatherApiChain
}

type Services struct {
	weatherService   *weather.WeatherService
	subscribeService *subscription.SubscribeService
}

func declareQueues(r *rabbitmq.RabbitMQ) error {
	queues := []string{
		rabbitmq.SendEmail,
		rabbitmq.WeatherUpdate,
	}

	for _, q := range queues {
		_, err := r.Channel.QueueDeclare(q, true, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", q, err)
		}
	}
	return nil
}
