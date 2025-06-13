package controller

import (
	"net/http"

	"github.com/ValeriiaHuza/weather_api/service"
	"github.com/gin-gonic/gin"
)

type WeatherController struct {
	Service *service.WeatherService
}

func NewWeatherController(service *service.WeatherService) *WeatherController {
	return &WeatherController{Service: service}
}

func (wc *WeatherController) GetWeather(c *gin.Context) {
	city := c.Query("city")

	weather, err := wc.Service.GetWeather(city)

	if err != nil {
		c.String(err.StatusCode, err.Message)
		return
	}
	c.JSON(http.StatusOK, weather)
}
