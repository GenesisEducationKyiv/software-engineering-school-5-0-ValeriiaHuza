package subscription

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type subscribeService interface {
	SubscribeForWeatherUpdates(email string, city string, frequency Frequency) error
	ConfirmSubscription(token string) error
	Unsubscribe(token string) error
	GetConfirmedSubscriptionsByFrequency(freq Frequency) []Subscription
	SendSubscriptionEmails(freq Frequency)
}

type SubscribeController struct {
	service subscribeService
}

func NewSubscribeController(service subscribeService) *SubscribeController {
	return &SubscribeController{service: service}
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

	frequency, err := sc.validateSubscriptionInput(body.Email, body.City, body.Frequency)
	if err != nil {
		HandleError(c, err)
		return
	}

	errRes := sc.service.SubscribeForWeatherUpdates(body.Email, body.City, frequency)

	if errRes != nil {
		HandleError(c, err)
		return
	}

	c.String(http.StatusOK, "Subscription successful. Confirmation email sent.")
}

func (sc *SubscribeController) ConfirmSubscription(c *gin.Context) {
	token := c.Param("token")

	if token == "" {
		c.String(http.StatusBadRequest, ErrInvalidToken.Error())
		return
	}

	err := sc.service.ConfirmSubscription(token)

	if err != nil {
		HandleError(c, err)
		return
	}

	c.String(http.StatusOK, "You confirmed weather update.")
}

func (sc *SubscribeController) Unsubscribe(c *gin.Context) {
	token := c.Param("token")

	if token == "" {
		c.String(http.StatusBadRequest, ErrInvalidToken.Error())
		return
	}

	err := sc.service.Unsubscribe(token)

	if err != nil {
		HandleError(c, err)
		return
	}

	c.String(http.StatusOK, "You unsubscribe from weather update.")
}

func (sc *SubscribeController) validateSubscriptionInput(email string,
	city string, frequencyStr string) (Frequency, error) {
	if email == "" || city == "" || frequencyStr == "" {
		return "", ErrInvalidInput
	}

	if !sc.isValidEmail(email) {
		return "", ErrInvalidInput
	}

	frequency, err := ParseFrequency(frequencyStr)
	if err != nil {
		return "", ErrInvalidInput
	}

	return frequency, nil
}

func (sc *SubscribeController) isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
