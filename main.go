package main

import (
	"fund/bot"
	"fund/config"
	"fund/db"
	"fund/global"
	"fund/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	err error
)

func init() {
	botToken := config.ViperEnvVariable("BOT_TOKEN")
	global.Bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Init Bot %+v", err)
	}
	global.Bot.Debug = true
	log.Info("Authorized on account %s", global.Bot.Self.UserName)

	mgoPwd := config.ViperEnvVariable("MGO_PWD")
	global.MgoDB = db.NewDB(mgoPwd)
}

func main() {
	bot.Run()

	// TODO: Email
	//e := &email.Email{
	//	To: "dukeyunz@hotmail.com",
	//	Subject: "Fund notification",
	//}
}

