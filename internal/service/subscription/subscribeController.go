package subscription

import (
	"net/http"

	appErr "github.com/ValeriiaHuza/weather_api/error"
	"github.com/gin-gonic/gin"
)

type subscribeService interface {
	SubscribeForWeatherUpdates(email, city, frequency string) *appErr.AppError
	ConfirmSubscription(token string) *appErr.AppError
	Unsubscribe(token string) *appErr.AppError
	GetConfirmedSubscriptionsByFrequency(freq Frequency) []Subscription
	SendSubscriptionEmails(freq Frequency)
}

type SubscribeController struct {
	Service subscribeService
}

func NewSubscribeController(service subscribeService) *SubscribeController {
	return &SubscribeController{Service: service}
}

func (sc *SubscribeController) ConfirmSubscription(c *gin.Context) {
	token := c.Param("token")

	err := sc.Service.ConfirmSubscription(token)

	if err != nil {
		c.String(err.StatusCode, err.Message)
		return
	}

	c.String(http.StatusOK, "You confirmed weather update.")
}

func (sc *SubscribeController) Unsubscribe(c *gin.Context) {
	token := c.Param("token")

	err := sc.Service.Unsubscribe(token)

	if err != nil {
		c.String(err.StatusCode, err.Message)
		return
	}

	c.String(http.StatusOK, "You unsubscribe from weather update.")
}

func (sc *SubscribeController) SubscribeForWeatherUpdates(c *gin.Context) {

	var body struct {
		Email     string `json:"email"`
		City      string `json:"city"`
		Frequency string `json:"frequency"`
	}

	err := c.Bind(&body)

	if err != nil {
		c.String(http.StatusBadRequest, "Invalid input")
		return
	}

	errRes := sc.Service.SubscribeForWeatherUpdates(body.Email, body.City, body.Frequency)

	if errRes != nil {
		c.String(errRes.StatusCode, errRes.Message)
		return
	}

	c.String(http.StatusOK, "Subscription successful. Confirmation email sent.")
}
