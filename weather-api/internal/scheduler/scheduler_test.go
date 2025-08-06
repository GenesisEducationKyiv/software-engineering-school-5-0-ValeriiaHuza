//go:build unit
// +build unit

package scheduler

import (
	"testing"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"

	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type mockSubscribeService struct {
	mock.Mock
}

func (m *mockSubscribeService) SendSubscriptionEmails(freq subscription.Frequency) {
	m.Called(freq)
}

// --- Tests ---

func TestStartCronJobs_SchedulesJobs(t *testing.T) {
	mockService := new(mockSubscribeService)

	// Set expectations
	mockService.On("SendSubscriptionEmails", subscription.FrequencyDaily).Return()
	mockService.On("SendSubscriptionEmails", subscription.FrequencyHourly).Return()

	mockLog, _ := logger.NewLogger()

	scheduler := NewScheduler(mockService, mockLog) // Assuming constructor exists
	scheduler.StartCronJobs()

	mockService.SendSubscriptionEmails(subscription.FrequencyDaily)
	mockService.SendSubscriptionEmails(subscription.FrequencyHourly)

	// Assert expectations
	mockService.AssertCalled(t, "SendSubscriptionEmails", subscription.FrequencyDaily)
	mockService.AssertCalled(t, "SendSubscriptionEmails", subscription.FrequencyHourly)
	mockService.AssertExpectations(t)
}
