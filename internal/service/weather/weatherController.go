package weather

import (
	"net/http"

	appErr "github.com/ValeriiaHuza/weather_api/error"
	"github.com/gin-gonic/gin"
)

type weatherService interface {
	GetWeather(city string) (*WeatherDTO, *appErr.AppError)
}

type WeatherController struct {
	service weatherService
}

func NewWeatherController(service weatherService) *WeatherController {
	return &WeatherController{service: service}
}

func (wc *WeatherController) GetWeather(c *gin.Context) {
	city := c.Query("city")

	weather, err := wc.service.GetWeather(city)

	if err != nil {
		c.String(err.StatusCode, err.Message)
		return
	}
	c.JSON(http.StatusOK, weather)
}
