package scheduler

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/weather-api/internal/service/subscription"
	"github.com/robfig/cron/v3"
)

type loggerInterface interface {
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
}

type subscribeService interface {
	SendSubscriptionEmails(freq subscription.Frequency)
}

type Scheduler struct {
	subscribeService subscribeService
	logger           loggerInterface
}

func NewScheduler(subscribeService subscribeService, logger loggerInterface) *Scheduler {
	return &Scheduler{
		subscribeService: subscribeService,
		logger:           logger,
	}
}

func (ss *Scheduler) StartCronJobs() {
	c := cron.New()

	// at 9 oclock
	if _, err := c.AddFunc("0 9 * * *", func() {
		ss.subscribeService.SendSubscriptionEmails(subscription.FrequencyDaily)
	}); err != nil {
		ss.logger.Error("Failed to schedule daily job", "error", err)

	}

	// Every hour
	if _, err := c.AddFunc("0 * * * *", func() {
		ss.subscribeService.SendSubscriptionEmails(subscription.FrequencyHourly)
	}); err != nil {
		ss.logger.Error("Failed to schedule hourly job", "error", err)
	}

	c.Start()
}
