package main

import (
	"encoding/json"
	"fmt"
	"fund/bot"
	"fund/config"
	"fund/cron"
	"fund/db"
	"fund/global"
	"fund/log"
	"fund/notifier"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	err error
)

func init() {
	botToken := config.EnvVariable("BOT_TOKEN")
	global.Bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Init Bot %+v", err)
	}
	global.Bot.Debug = true
	log.Info("Authorized on account %s", global.Bot.Self.UserName)

	mgoPwd := config.EnvVariable("MGO_PWD")
	global.MgoDB = db.NewDB(mgoPwd)
}

func main() {
	// local
	if len(os.Args) > 1 {
		go cron.Update()
		bot.Run()
	} else {
		port := config.EnvVariable("PORT")
		fmt.Println("Server started on port:", port)
		global.Bot.RemoveWebhook()
		_, err = global.Bot.SetWebhook(tgbotapi.NewWebhook("https://thawing-scrubland-62700.herokuapp.com/" + global.Bot.Token))
		if err != nil {
			log.Fatal("%v", err)
		}

		router := gin.New()
		router.Use(gin.Logger())
		router.Use(gin.Recovery())

		// router.POST("/"+global.Bot.Token, webhookHandler)
		router.GET("/hello", hello)
		router.GET("/reminder", func(c *gin.Context) {
			e := &notifier.Email{
				To:      "dukeyunz@hotmail.com",
				Subject: "Fund notification",
			}
			e.Send("Test email from heroku every 5min.")
			c.String(http.StatusOK, "ok")
		})

		err := router.Run(":" + port)
		if err != nil {
			log.Error("%v", err)
		}
	}
}

func hello(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

func webhookHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("%+v", err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Error("%+v", err)
		return
	}

	// to monitor changes run: heroku logs --tail
	log.Info("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
	bot.Handle(update)
}
