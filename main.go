package main

import (
	"fund/bot"
	"fund/config"
	"fund/db"
	"fund/global"
	"fund/log"
	"os"

	http "net/http"

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

// func main() {
// 	cron.Update()
// 	bot.Run()

// 	// TODO: Email
// 	//e := &email.Email{
// 	//	To: "dukeyunz@hotmail.com",
// 	//	Subject: "Fund notification",
// 	//}
// }

func main() {
	port := os.Getenv("PORT")
	_, err = global.Bot.SetWebhook(tgbotapi.NewWebhook("https://thawing-scrubland-62700.herokuapp.com/" + global.Bot.Token))
	if err != nil {
		log.Fatal("%v", err)
	}

	updates := global.Bot.ListenForWebhook("/" + global.Bot.Token)
	go http.ListenAndServe("0.0.0.0"+":"+port, nil)

	for update := range updates {
		bot.Handle(update)
	}
}
