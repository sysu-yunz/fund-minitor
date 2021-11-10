package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

// cron job every second
func Update() {
	go func() {
		c := cron.New()
		c.AddFunc("@every 1m", func() {
			fmt.Println("Every second")
		})
		c.Start()
		// select {}
	}()
}
