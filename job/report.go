package job

import (
	"fund/notifier"
	"net/http"

	"github.com/gin-gonic/gin"
)

// find schedulers in https://cronhub.io/schedulers

func DailyReport(c *gin.Context) {
	e := &notifier.Email{
		To:      "dukeyunz@hotmail.com",
		Subject: "Fund notification",
	}
	e.Send("Test email from heroku every 5min.")
	c.String(http.StatusOK, "ok")
}
