//go:build integration
// +build integration

package integration

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/subscription"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func TestSubscriptionEndpoint_SuccessCreate(t *testing.T) {
// 	// create a new subscription
// 	email := "test@example.com"
// 	city := "Kyiv"
// 	resp := createTestSubscription(t, email, city)
// 	require.Equal(t, http.StatusOK, resp.Code)
// 	require.Contains(t, resp.Body.String(), "Subscription successful")

// 	//check subscription in the database
// 	sub := getSubscriptionByEmail(t, email)
// 	assert.Equal(t, city, sub.City)
// 	assert.False(t, sub.Confirmed)

// 	//check confirmation email was sent
// 	emailSent := testMailService.SentEmail
// 	assert.Equal(t, email, emailSent.To)
// 	assert.Equal(t, "Weather updates confirmation link", emailSent.Subject)
// 	assert.Contains(t, emailSent.Body, sub.Token)
// }

// func TestSubscriptionEndpoint_RepeatedSubscription(t *testing.T) {
// 	email := "repeat@example.com"
// 	city := "Kyiv"

// 	// create first subscription
// 	resp := createTestSubscription(t, email, city)
// 	require.Equal(t, http.StatusOK, resp.Code)
// 	require.Contains(t, resp.Body.String(), "Subscription successful")

// 	// create second subscription with the same email
// 	resp2 := createTestSubscription(t, email, city)
// 	assert.Equal(t, http.StatusConflict, resp2.Code)
// 	assert.Equal(t, resp2.Body.String(), subscription.ErrEmailAlreadySubscribed.Error())
// }

// func TestSubscriptionEndpoint_ConfirmSubscription_Success(t *testing.T) {
// 	email := "confirm@example.com"
// 	city := "Kyiv"

// 	// create subscription
// 	resp := createTestSubscription(t, email, city)
// 	require.Equal(t, http.StatusOK, resp.Code)
// 	require.Contains(t, resp.Body.String(), "Subscription successful")

// 	// get the subscription from the DB
// 	sub := getSubscriptionByEmail(t, email)

// 	// confirm subscription
// 	respConfirm := confirmSubscription(t, sub.Token)
// 	assert.Equal(t, http.StatusOK, respConfirm.Code)
// 	assert.Contains(t, respConfirm.Body.String(), "confirmed")

// 	// Check if the confirmation email was sent
// 	emailSent := testMailService.SentEmail
// 	assert.Equal(t, email, emailSent.To)
// 	assert.Equal(t, "Weather updates subscription", emailSent.Subject)

// 	// Check DB updated
// 	sub = getSubscriptionByEmail(t, email)
// 	assert.True(t, sub.Confirmed)
// }

// func TestSubscriptionEndpoint_ConfirmSubscription_Unsuccessful(t *testing.T) {
// 	confirmResp := confirmSubscription(t, "invalid-token")

// 	assert.Equal(t, http.StatusBadRequest, confirmResp.Code)
// 	assert.Equal(t, confirmResp.Body.String(), subscription.ErrTokenNotFound.Error())
// }

// func TestSubscriptionEndpoint_Unsubscribe_Success(t *testing.T) {
// 	// Create a subscription
// 	email := "unsub@example.com"
// 	city := "Kyiv"
// 	resp := createTestSubscription(t, email, city)
// 	require.Equal(t, http.StatusOK, resp.Code)
// 	require.Contains(t, resp.Body.String(), "Subscription successful")

// 	// Get subscription from DB
// 	sub := getSubscriptionByEmail(t, email)

// 	// Unsubscribe
// 	respUns := unsubscribe(t, sub.Token)
// 	assert.Equal(t, http.StatusOK, respUns.Code)
// 	assert.Contains(t, respUns.Body.String(), "unsubscribe")

// 	// Check DB deleted
// 	_, err := testRepo.FindByEmail(email)
// 	assert.Error(t, err)
// }

// func TestSubscriptionEndpoint_Unsubscribe_Unsuccessful(t *testing.T) {
// 	// Try to unsubscribe with an invalid token
// 	unsubResp := unsubscribe(t, "invalid-token")

// 	assert.Equal(t, http.StatusBadRequest, unsubResp.Code)
// 	assert.Equal(t, unsubResp.Body.String(), subscription.ErrTokenNotFound.Error())
// }

// func createTestSubscription(t *testing.T, email string, city string) *httptest.ResponseRecorder {
// 	t.Helper()
// 	body := fmt.Sprintf(`{
// 		"email": "%s",
// 		"city": "%s",
// 		"frequency": "daily"
// 	}`, email, city)

// 	req := httptest.NewRequest("POST", "/api/subscribe", strings.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp := httptest.NewRecorder()

// 	testRouter.ServeHTTP(resp, req)

// 	return resp
// }

// func getSubscriptionByEmail(t *testing.T, email string) *subscription.Subscription {
// 	t.Helper()
// 	sub, err := testRepo.FindByEmail(email)
// 	require.NoError(t, err)
// 	return sub
// }

// func confirmSubscription(t *testing.T, token string) *httptest.ResponseRecorder {
// 	t.Helper()
// 	req := httptest.NewRequest("GET", fmt.Sprintf("/api/confirm/%s", token), nil)
// 	resp := httptest.NewRecorder()
// 	testRouter.ServeHTTP(resp, req)
// 	return resp
// }

// func unsubscribe(t *testing.T, token string) *httptest.ResponseRecorder {
// 	t.Helper()
// 	req := httptest.NewRequest("GET", fmt.Sprintf("/api/unsubscribe/%s", token), nil)
// 	resp := httptest.NewRecorder()
// 	testRouter.ServeHTTP(resp, req)
// 	return resp
// }
