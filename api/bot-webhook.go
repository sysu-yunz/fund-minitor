package handler

import (
	"encoding/json"
	"fund/bot"
	"fund/data"
	"fund/db"
	"fund/global"
	"fund/log"
	"io/ioutil"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	err error
)

func init() {
	botToken := os.Getenv("BOT_TOKEN")
	global.Bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Init Bot %+v", err)
	}
	global.Bot.Debug = true
	log.Info("Authorized on account %s", global.Bot.Self.UserName)

	mgoPwd := os.Getenv("MGO_PWD")
	global.MgoDB = db.NewDB(mgoPwd)
	go data.UpdateCookie()
}

func Handler(w http.ResponseWriter, req *http.Request) {
	bytes, _ := ioutil.ReadAll(req.Body)

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Error("Unmarshal update: %v", err)
	}
	log.Debug("Update: %+v %+v", update.Message.Text, update)

	bot.Handle(update)
}
