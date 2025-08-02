package scheduler

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/service/subscription"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/logger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type subscribeService interface {
	SendSubscriptionEmails(freq subscription.Frequency)
}

type Scheduler struct {
	subscribeService subscribeService
}

func NewScheduler(subscribeService subscribeService) *Scheduler {
	return &Scheduler{
		subscribeService: subscribeService,
	}
}

func (ss *Scheduler) StartCronJobs() {
	c := cron.New()

	// at 9 oclock
	if _, err := c.AddFunc("0 9 * * *", func() {
		ss.subscribeService.SendSubscriptionEmails(subscription.FrequencyDaily)
	}); err != nil {
		logger.GetLogger().Error("Failed to schedule daily job", zap.Error(err))
	}

	// Every hour
	if _, err := c.AddFunc("0 * * * *", func() {
		ss.subscribeService.SendSubscriptionEmails(subscription.FrequencyHourly)
	}); err != nil {
		logger.GetLogger().Error("Failed to schedule hourly job", zap.Error(err))
	}

	c.Start()
}
