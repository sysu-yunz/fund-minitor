package handler

import (
	"encoding/json"
	"fund/command"
	"fund/log"
	. "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
	"os"
)

var bot, _ = NewBotAPI(os.Getenv("BOT_TOKEN"))

func init() {
	bot.Debug = true
	log.Debug("Authorized on account %s", bot.Self.UserName)
}

func Handler(w http.ResponseWriter, req *http.Request) {
	bytes, _ := ioutil.ReadAll(req.Body)
	var update Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		log.Error("Unmarshal update: ", err)
	}
	log.Debug("Update: %+v %+v", update.Message.Text, update)

	go command.Command(bot, update)
}
