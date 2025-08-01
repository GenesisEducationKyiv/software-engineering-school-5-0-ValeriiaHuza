//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/client"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/weather"
	"github.com/stretchr/testify/assert"
)

func TestWeatherEndpoint_Scenarios(t *testing.T) {
	tests := []struct {
		name           string
		cityQuery      string
		expectedStatus int
		expectBody     string
	}{
		{"valid city", "Kyiv", http.StatusOK, `{"temperature":21.5,"humidity":55,"description":"Sunny"}`},
		{"missing city", "", http.StatusBadRequest, weather.ErrInvalidCityInput.Error()},
		{"city not found", "Nowhere", http.StatusNotFound, client.ErrCityNotFound.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/weather?city=%s", tt.cityQuery), nil)
			resp := httptest.NewRecorder()

			testRouter.ServeHTTP(resp, req)

			t.Logf("Response: %d - %s", resp.Code, resp.Body.String())

			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Equal(t, tt.expectBody, resp.Body.String())
		})
	}
}
