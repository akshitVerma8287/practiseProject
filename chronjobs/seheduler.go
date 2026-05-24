package cronjobs

import (
	"fmt"
	"project/services"
	"time"

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

func StartApiJob() {

	ticker := time.NewTicker(10 * time.Second)

	go func() {

		for {

			select {

			case <-ticker.C:

				fmt.Println("Running API Job...")

				services.CallMultipleApis()

			}
		}

	}()
}