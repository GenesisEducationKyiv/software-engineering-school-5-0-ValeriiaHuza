package weather

import (
	"errors"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client"
	"github.com/gin-gonic/gin"
)

type weatherService interface {
	GetWeather(city string) (*client.WeatherDTO, error)
}

type WeatherController struct {
	service weatherService
}

func NewWeatherController(service weatherService) *WeatherController {
	return &WeatherController{service: service}
}

func (wc *WeatherController) GetWeather(c *gin.Context) {
	city, err := validateCityQuery(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	response, err := wc.service.GetWeather(city)

	if err != nil {
		switch {
		case errors.Is(err, client.ErrCityNotFound):
			c.String(http.StatusNotFound, err.Error())
			return
		case errors.Is(err, client.ErrInvalidRequest):
			c.String(http.StatusBadRequest, err.Error())
			return
		default:
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func validateCityQuery(c *gin.Context) (string, error) {
	city := c.Query("city")
	if city == "" {
		return "", ErrInvalidCityInput
	}
	return city, nil
}
