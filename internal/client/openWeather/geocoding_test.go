//go:build unit
// +build unit

package openweather

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/ValeriiaHuza/weather_api/internal/client"
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

// --- Tests ---

func TestGetCityCoordinates_Success(t *testing.T) {
	apiKey := "key"
	apiUrl := "open-weather"
	mockBody := `[{"lat":50.45,"lon":30.52}]`
	mockClient := newMockClient(mockBody, 200, nil)

	client := NewGeocodingClient(apiKey, apiUrl, mockClient)

	coords, err := client.GetCityCoordinates("Kyiv")
	assert.NoError(t, err)
	assert.NotNil(t, coords)
	assert.Equal(t, 50.45, coords.Lat)
	assert.Equal(t, 30.52, coords.Lon)
}

func TestGetCityCoordinates_NoCityFound(t *testing.T) {
	apiKey := "key"
	apiUrl := "open-weather"
	mockClient := newMockClient("[]", 200, nil)

	geoClient := NewGeocodingClient(apiKey, apiUrl, mockClient)

	coords, err := geoClient.GetCityCoordinates("Kyiv")
	assert.Error(t, err)
	assert.Equal(t, client.ErrCityNotFound.Error(), err.Error())
	assert.Nil(t, coords)
}

func TestGetCityCoordinates_Non200Status(t *testing.T) {
	apiKey := "key"
	apiUrl := "open-weather"
	mockClient := newMockClient("not found", 404, nil)

	client := NewGeocodingClient(apiKey, apiUrl, mockClient)

	coords, err := client.GetCityCoordinates("Kyiv")
	assert.Error(t, err)
	assert.Nil(t, coords)
}

func TestGetCityCoordinates_InvalidJSON(t *testing.T) {
	apiKey := "key"
	apiUrl := "open-weather"
	mockClient := newMockClient("{invalid json", 200, nil)

	client := NewGeocodingClient(apiKey, apiUrl, mockClient)

	coords, err := client.GetCityCoordinates("Kyiv")
	assert.Error(t, err)
	assert.Nil(t, coords)
}
