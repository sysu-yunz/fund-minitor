package main

import (
	"fmt"
	"fund/bot"
	"fund/config"
	"fund/cron"
	"fund/db"
	"fund/global"
	"fund/log"
	"fund/notifier"
	"net/http"
	"os"

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

		updates := global.Bot.ListenForWebhook("/" + global.Bot.Token)
		go http.ListenAndServe("0.0.0.0"+":"+port, SendEmail{})
		go cron.Update()

		for update := range updates {
			bot.Handle(update)
		}
	}
}

type SendEmail struct{}

func (SendEmail) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("%v", r)
	fmt.Println("定时发邮件任务")
	e := &notifier.Email{
		To:      "dukeyunz@hotmail.com",
		Subject: "Fund notification",
	}
	e.Send("Test email from heroku every 5min.")
	w.Write([]byte("ok"))
}
