//go:build integration
// +build integration

package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionEndpoint_Create(t *testing.T) {

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
	require.Len(t, testMailService.SentEmails, 1)
	email := testMailService.SentEmails[0]
	assert.Equal(t, "test@example.com", email.To)
	assert.Equal(t, "Weather updates confirmation link", email.Subject)
	assert.Contains(t, email.Body, sub.Token)
}
