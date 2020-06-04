package handler

import (
	"context"
	"encoding/json"
	"fund/command"
	"fund/db"
	"fund/log"
	"github.com/globalsign/mgo/bson"
	. "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	bot, err := NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Error("Init Bot %+v", err)
	}

	bot.Debug = true
	log.Debug("Authorized on account %s", bot.Self.UserName)

	dBase := db.NewDB(os.Getenv("MGO_PWD"))
	dBase.ListDatabaseNames(context.TODO(), bson.M{})

	bytes, _ := ioutil.ReadAll(req.Body)
	var update Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Error("Unmarshal update: ", err)
	}
	log.Debug("Update: %+v %+v", update.Message.Text, update)

	command.Command(bot, update)
}
