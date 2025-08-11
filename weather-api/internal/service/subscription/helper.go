package subscription

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidRequest),
		errors.Is(err, ErrInvalidInput),
		errors.Is(err, ErrInvalidToken),
		errors.Is(err, ErrTokenNotFound),
		errors.Is(err, ErrFailedToSaveSubscription):
		c.String(http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrEmailAlreadySubscribed):
		c.String(http.StatusConflict, err.Error())

	default:
		c.String(http.StatusBadRequest, err.Error())
	}
}
