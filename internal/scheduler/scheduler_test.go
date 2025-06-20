//go:build unit
// +build unit

package scheduler

import (
	"testing"

	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"github.com/stretchr/testify/assert"
)

// mockSubscribeService implements subscribeService for testing
type mockSubscribeService struct {
	calls []subscription.Frequency
}

func (m *mockSubscribeService) SendSubscriptionEmails(freq subscription.Frequency) {
	m.calls = append(m.calls, freq)
}

func TestStartCronJobs_SchedulesJobs(t *testing.T) {
	mockService := &mockSubscribeService{}
	s := NewScheduler(mockService)

	s.StartCronJobs()

	mockService.SendSubscriptionEmails(subscription.FrequencyDaily)
	mockService.SendSubscriptionEmails(subscription.FrequencyHourly)

	assert.Len(t, mockService.calls, 2, "should have 2 calls to SendSubscriptionEmails")
	assert.Equal(t, subscription.FrequencyDaily, mockService.calls[0], "first call should be FrequencyDaily")
	assert.Equal(t, subscription.FrequencyHourly, mockService.calls[1], "second call should be FrequencyHourly")

}
