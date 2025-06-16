package scheduler

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/internal/service/subscription"
	"github.com/robfig/cron/v3"
)

type subscribeService interface {
	SendSubscriptionEmails(freq subscription.Frequency)
}

type Scheduler struct {
	SubscribeService subscribeService
}

func NewScheduler(subscribeService subscribeService) *Scheduler {
	return &Scheduler{
		SubscribeService: subscribeService,
	}
}

func (ss *Scheduler) StartCronJobs() {
	c := cron.New()

	// at 9 oclock
	if _, err := c.AddFunc("0 9 * * *", func() {
		ss.SubscribeService.SendSubscriptionEmails(subscription.FrequencyDaily)
	}); err != nil {
		log.Println("Failed to schedule daily job:", err)
	}

	// Every hour
	if _, err := c.AddFunc("0 * * * *", func() {
		ss.SubscribeService.SendSubscriptionEmails(subscription.FrequencyHourly)
	}); err != nil {
		log.Println("Failed to schedule hourly job:", err)
	}

	c.Start()
}
