package cron

import (
	"fmt"
	"fund/notifier"

	"github.com/robfig/cron/v3"
)

// cron job every second
func Update() {
	go func() {
		c := cron.New()
		c.AddFunc("@every 5m", func() {
			fmt.Println("Every second")
			e := &notifier.Email{
				To:      "dukeyunz@hotmail.com",
				Subject: "Fund notification",
			}
			e.Send("Test email from heroku every 5min.")
		})
		c.Start()
		// select {}
	}()
}
