package scheduler

import (
	"log"

	"github.com/ValeriiaHuza/weather_api/models"
	"github.com/ValeriiaHuza/weather_api/service"
	"github.com/robfig/cron"
)

type Scheduler struct {
	SubscribeService service.SubscribeService
}

func NewScheduler(subscribeService service.SubscribeService) *Scheduler {
	return &Scheduler{
		SubscribeService: subscribeService,
	}
}

func (ss *Scheduler) StartCronJobs() {
	c := cron.New()

	//at 9 oclock
	if err := c.AddFunc("0 9 * * *", func() {
		ss.SubscribeService.SendSubscriptionEmails(models.FrequencyDaily)
	}); err != nil {
		log.Println("Failed to schedule daily job:", err)
	}

	// Every hour
	if err := c.AddFunc("0 * * * *", func() {
		ss.SubscribeService.SendSubscriptionEmails(models.FrequencyHourly)
	}); err != nil {
		log.Println("Failed to schedule hourly job:", err)
	}

	c.Start()
}
