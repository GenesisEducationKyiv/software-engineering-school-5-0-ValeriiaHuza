package routes

import (
	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"github.com/ValeriiaHuza/weather_api/internal/service/weather"
	"github.com/gin-gonic/gin"
)

func WeatherRoute(router *gin.RouterGroup, weatherController *weather.WeatherController) {

	router.GET("/weather", weatherController.GetWeather)

}

func SubscribeRoute(router *gin.RouterGroup, subscribeController *subscription.SubscribeController) {

	router.POST("/subscribe", subscribeController.SubscribeForWeatherUpdates)
	router.GET("/confirm/:token", subscribeController.ConfirmSubscription)
	router.GET("/unsubscribe/:token", subscribeController.Unsubscribe)

}
