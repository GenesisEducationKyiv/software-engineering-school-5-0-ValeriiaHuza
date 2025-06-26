//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionEndpoint_SuccessCreate(t *testing.T) {

	body := `{
		"email": "test@example.com",
		"city": "Kyiv",
		"frequency": "daily"
	}`

	req := httptest.NewRequest("POST", "/api/subscribe", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Subscription successful")

	// --- Assert DB ---
	sub, err := testRepo.FindByEmail("test@example.com")
	require.NoError(t, err)
	assert.Equal(t, "Kyiv", sub.City)
	assert.False(t, sub.Confirmed)
	assert.Equal(t, "daily", string(sub.Frequency))

	// --- Assert email sent ---
	email := testMailService.SentEmail
	assert.Equal(t, "test@example.com", email.To)
	assert.Equal(t, "Weather updates confirmation link", email.Subject)
	assert.Contains(t, email.Body, sub.Token)
}

func TestSubscriptionEndpoint_RepeatedSubscription(t *testing.T) {
	body := `{
        "email": "repeat@example.com",
        "city": "Kyiv",
        "frequency": "daily"
    }`

	// First subscription should succeed
	req1 := httptest.NewRequest("POST", "/api/subscribe", strings.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	resp1 := httptest.NewRecorder()
	testRouter.ServeHTTP(resp1, req1)
	assert.Equal(t, http.StatusOK, resp1.Code)

	// Second subscription with same email should fail (already subscribed)
	req2 := httptest.NewRequest("POST", "/api/subscribe", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()
	testRouter.ServeHTTP(resp2, req2)
	assert.Equal(t, http.StatusConflict, resp2.Code)
	assert.Contains(t, resp2.Body.String(), subscription.ErrEmailAlreadySubscribed.Error())
}

func TestSubscriptionEndpoint_ConfirmSubscription_Success(t *testing.T) {
	// Create a subscription
	body := `{
        "email": "confirm@example.com",
        "city": "Kyiv",
        "frequency": "daily"
    }`
	req := httptest.NewRequest("POST", "/api/subscribe", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Get the token from the DB
	sub, err := testRepo.FindByEmail("confirm@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, sub.Token)

	// Confirm the subscription
	confirmURL := fmt.Sprintf("/api/confirm/%s", sub.Token)
	confirmReq := httptest.NewRequest("GET", confirmURL, nil)
	confirmResp := httptest.NewRecorder()
	testRouter.ServeHTTP(confirmResp, confirmReq)
	assert.Equal(t, http.StatusOK, confirmResp.Code)
	assert.Contains(t, confirmResp.Body.String(), "confirmed")

	// Check if the confirmation email was sent
	email := testMailService.SentEmail
	assert.Equal(t, "confirm@example.com", email.To)
	assert.Equal(t, "Weather updates subscription", email.Subject)
	assert.Contains(t, email.Body, "confirm@example.com")

	// Check DB updated
	sub, err = testRepo.FindByEmail("confirm@example.com")
	require.NoError(t, err)
	assert.True(t, sub.Confirmed)
}
