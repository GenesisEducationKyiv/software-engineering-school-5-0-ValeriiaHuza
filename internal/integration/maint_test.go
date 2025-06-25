package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ValeriiaHuza/weather_api/internal/client"
	"github.com/ValeriiaHuza/weather_api/internal/service/weather"
	"github.com/gin-gonic/gin"
)

var testRouter *gin.Engine

func setupRouter() *gin.Engine {
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

	client := client.NewWeatherAPIClient("fake-key", fakeWeatherServer.URL, http.DefaultClient)
	service := weather.NewWeatherAPIService(client)
	controller := weather.NewWeatherController(service)

	r := gin.Default()
	r.GET("/weather", controller.GetWeather)
	return r
}

func TestMain(m *testing.M) {
	testRouter = setupRouter()
	m.Run()
}
