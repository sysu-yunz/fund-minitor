package job

import (
	"fmt"
	"fund/notifier"

	"github.com/robfig/cron/v3"
)

// cron job every second
func Update() {
	go func() {
		c := cron.New()
		c.AddFunc("@every 1d", func() {
			fmt.Println("定时发邮件任务")
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
