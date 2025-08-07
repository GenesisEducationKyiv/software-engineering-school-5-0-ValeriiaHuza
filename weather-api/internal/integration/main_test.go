//go:build integration
// +build integration

package integration

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	weatherapi "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client/weatherApi"
	dbPackage "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/rabbitmq"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/redis"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/repository"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/routes"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/service/weather"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"

	"github.com/gin-gonic/gin"
)

var (
	testRouter      *gin.Engine
	testRepo        *repository.SubscriptionRepository
	terminateDB     func()
	terminateRedis  func()
	terminateRabbit func()
)

func setupRouter() (*gin.Engine, *repository.SubscriptionRepository, func()) {
	ctx := context.Background()

	logger, err := logger.NewTestLogger()

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Setup Postgres container
	db, terminateDB, err := SetupPostgresContainer()
	if err != nil {
		log.Fatalf("Failed to setup test DB: %v", err)
	}

	dbPackage.AutomatedMigration(db)

	// Setup Redis container
	redisTest, terminateRedis, err := SetupRedisContainer()
	if err != nil {
		log.Fatalf("Failed to setup test redis: %v", err)
	}

	// Setup RabbitMQ container
	rabbitMQTest, terminateRabbit, err := SetupRabbitMQContainer()
	if err != nil {
		log.Fatalf("Failed to setup test rabbitmq: %v", err)
	}

	// Setup fake weather API
	fakeWeatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("q")
		switch city {
		case "Nowhere":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"code":1006,"message":"No matching location found."}}`))
		default:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
                "current": {
                    "temp_c": 21.5,
                    "humidity": 55,
                    "condition": { "text": "Sunny" }
                }
            }`))
		}
	}))

	redisProvider := redis.NewRedisProvider(redisTest, ctx, logger)
	fakeWeatherClient := weatherapi.NewWeatherAPIClient("fake-key", fakeWeatherServer.URL, http.DefaultClient, logger)
	weatherChain := client.NewWeatherChain(fakeWeatherClient, logger)
	weatherService := weather.NewWeatherAPIService(weatherChain, &redisProvider, logger)
	weatherController := weather.NewWeatherController(weatherService)

	repo := repository.NewSubscriptionRepository(db)
	emailPublisher := rabbitmq.NewRabbitMQPublisher(rabbitMQTest.Channel)
	subscribeService := subscription.NewSubscribeService(weatherService, repo, emailPublisher, logger)
	subscribeController := subscription.NewSubscribeController(subscribeService)

	r := gin.Default()
	api := r.Group("/api")
	routes.WeatherRoute(api, weatherController)
	routes.SubscribeRoute(api, subscribeController)

	// Single cleanup function in reverse order of initialization
	return r, repo, func() {
		if terminateRabbit != nil {
			terminateRabbit()
		}
		if terminateRedis != nil {
			terminateRedis()
		}
		if terminateDB != nil {
			terminateDB()
		}
	}
}

func TestMain(m *testing.M) {
	var cleanup func()

	testRouter, testRepo, cleanup = setupRouter()
	code := m.Run()

	if cleanup != nil {
		cleanup()
	}

	os.Exit(code)

}
