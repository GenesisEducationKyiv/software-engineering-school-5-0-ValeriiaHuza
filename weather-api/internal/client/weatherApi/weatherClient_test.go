//go:build unit
// +build unit

package weatherapi

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	packageClient "github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/client"
	"github.com/stretchr/testify/assert"
)

type MockRoundTripper struct {
	resp *http.Response
	err  error
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func newMockClient(respBody string, statusCode int, err error) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{
			resp: &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(respBody)),
				Header:     make(http.Header),
			},
			err: err,
		},
	}
}

func TestFetchWeather_Success(t *testing.T) {
	mockBody := `{
		"current": {
			"temp_c": 21.5,
			"humidity": 60,
			"condition": {
				"text": "Sunny"
			}
		}
	}`
	client := newMockClient(mockBody, 200, nil)
	apiClient := NewWeatherAPIClient("dummy-key", "api-url", client)

	result, err := apiClient.FetchWeather("London")

	assert.NoError(t, err)
	assert.Equal(t, result.Temperature, 21.5)
	assert.Equal(t, result.Humidity, 60.0)
	assert.Equal(t, result.Description, "Sunny")
}

func TestFetchWeather_HTTPError(t *testing.T) {
	client := newMockClient("", 0, errors.New("network error"))
	apiClient := NewWeatherAPIClient("dummy-key", "api-url", client)

	_, err := apiClient.FetchWeather("London")

	assert.Error(t, err)
}

func TestFetchWeather_BadJSON(t *testing.T) {
	mockBody := `not a json`
	client := newMockClient(mockBody, 200, nil)
	apiClient := NewWeatherAPIClient("dummy-key", "api-url", client)

	_, err := apiClient.FetchWeather("London")
	assert.Error(t, err)
}

func TestFetchWeather_APIError_CityNotFound(t *testing.T) {
	mockBody := `{
		"error": {
			"code": 1006,
			"message": "No matching location found."
		}
	}`
	client := newMockClient(mockBody, 200, nil)
	apiClient := NewWeatherAPIClient("dummy-key", "api-url", client)

	_, err := apiClient.FetchWeather("UnknownCity")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, packageClient.ErrCityNotFound), "expected ErrCityNotFound, got %v", err)
}

func TestFetchWeather_APIError_InvalidRequest(t *testing.T) {
	mockBody := `{
		"error": {
			"code": 2006,
			"message": "Invalid API key."
		}
	}`
	client := newMockClient(mockBody, 200, nil)
	apiClient := NewWeatherAPIClient("dummy-key", "api-url", client)

	_, err := apiClient.FetchWeather("London")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, packageClient.ErrInvalidRequest), "expected ErrInvalidRequest, got %v", err)
}
