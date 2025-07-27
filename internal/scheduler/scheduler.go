package scheduler

import (
	"log"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-ValeriiaHuza/internal/service/subscription"
	"github.com/robfig/cron/v3"
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
		log.Println("Failed to schedule daily job:", err)
	}

	// Every hour
	if _, err := c.AddFunc("0 * * * *", func() {
		ss.subscribeService.SendSubscriptionEmails(subscription.FrequencyHourly)
	}); err != nil {
		log.Println("Failed to schedule hourly job:", err)
	}

	c.Start()
}
