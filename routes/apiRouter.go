package routes

import (
	"github.com/ValeriiaHuza/weather_api/controller"
	"github.com/gin-gonic/gin"
)

func WeatherRoute(router *gin.RouterGroup, weatherController *controller.WeatherController) {

	router.GET("/weather", weatherController.GetWeather)

}

func SubscribeRoute(router *gin.RouterGroup, subscribeController *controller.SubscribeController) {

	router.POST("/subscribe", subscribeController.SubscribeForWeatherUpdates)
	router.GET("/confirm/:token", subscribeController.ConfirmSubscription)
	router.GET("/unsubscribe/:token", subscribeController.Unsubscribe)

}
