package job

import (
	"fmt"
	"fund/log"
	"fund/notifier"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// find schedulers in https://www.easycron.com/

func DailyReport(c *gin.Context) {
	username := os.Getenv("username")
	password := os.Getenv("password")
	u, p, ok := c.Request.BasicAuth()
	if !ok {
		fmt.Println("Error parsing basic auth")
		c.String(http.StatusUnauthorized, "Error parsing basic auth")
		return
	}
	if u != username {
		log.Info("Username provided is correct: %s\n", u)
		c.String(http.StatusUnauthorized, "Error parsing basic auth")
		return
	}
	if p != password {
		log.Info("Password provided is correct: %s\n", u)
		c.String(http.StatusUnauthorized, "Error parsing basic auth")
		return
	}
	log.Info("Username: %s\n", u)
	log.Info("Password: %s\n", p)
	e := &notifier.Email{
		To:      "dukeyunz@hotmail.com",
		Subject: "Fund notification",
	}
	e.Send("Hello from heroku !")
	c.String(http.StatusOK, "ok")
}
