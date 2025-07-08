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

	"github.com/ValeriiaHuza/weather_api/internal/client"
	dbPackage "github.com/ValeriiaHuza/weather_api/internal/db"
	"github.com/ValeriiaHuza/weather_api/internal/redis"
	"github.com/ValeriiaHuza/weather_api/internal/repository"
	"github.com/ValeriiaHuza/weather_api/internal/routes"
	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"github.com/ValeriiaHuza/weather_api/internal/service/weather"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client"
	weatherapi "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client/weatherApi"
	dbPackage "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/repository"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/routes"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/weather"
	"github.com/gin-gonic/gin"
)

var (
	testRouter      *gin.Engine
	testRepo        *repository.SubscriptionRepository
	testMailService *FakeMailService
	terminateDB     func()
	terminateRedis  func()
)

func setupRouter() (*gin.Engine, *repository.SubscriptionRepository, *FakeMailService, func(), func()) {
	ctx := context.Background()

	// Setup Postgres container
	db, terminateDB, err := SetupPostgresContainer()
	if err != nil {
		log.Fatalf("Failed to setup test DB: %v", err)
	}

	dbPackage.AutomatedMigration(db)

	redisTest, terminateRedis, err := SetupRedisContainer()
	if err != nil {
		log.Fatalf("Failed to setup test redis: %v", err)
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

<<<<<<< HEAD:internal/integration/maint_test.go
	redisProvider := redis.NewRedisProvider(redisTest, ctx)
	client := client.NewWeatherAPIClient("fake-key", fakeWeatherServer.URL, http.DefaultClient)
	weatherService := weather.NewWeatherAPIService(client, &redisProvider)
=======
	fakeWeatherClient := weatherapi.NewWeatherAPIClient("fake-key", fakeWeatherServer.URL, http.DefaultClient)

	weatherChain := client.NewWeatherChain(fakeWeatherClient)

	weatherService := weather.NewWeatherAPIService(weatherChain)
>>>>>>> main:internal/integration/main_test.go

	weatherController := weather.NewWeatherController(weatherService)

	fakeMailService := NewFakeMailService()
	repo := repository.NewSubscriptionRepository(db)
	subscribeService := subscription.NewSubscribeService(weatherService, fakeMailService, repo)

	subscribeController := subscription.NewSubscribeController(subscribeService)

	r := gin.Default()
	api := r.Group("/api")
	routes.WeatherRoute(api, weatherController)
	routes.SubscribeRoute(api, subscribeController)
	return r, repo, fakeMailService, terminateDB, terminateRedis
}

func TestMain(m *testing.M) {
	testRouter, testRepo, testMailService, terminateDB, terminateRedis = setupRouter()
	code := m.Run()

	if terminateDB != nil {
		terminateDB()
	}

	if terminateRedis != nil {
		terminateRedis()
	}

	os.Exit(code)

}
