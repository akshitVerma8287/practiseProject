package cronjobs

import (
	"project/services"

	"github.com/robfig/cron/v3"
)

func StartCronJobs() {

	c := cron.New()

	// Runs every 10 seconds
	c.AddFunc("@every 10s", func() {
		services.LogStudentsAbove18()
	})

	c.Start()
}