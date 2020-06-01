package handler

import (
	"encoding/json"
	"fund/cryptoc"
	"fund/data"
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

func handler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	var update Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		log.Error("Unmarshal update: ", err)
	}
	log.Debug("Update: %+v %+v", update.Message.Text, update)

	var reply string
	switch update.Message.Text {
	case "fund":
		reply = data.RealTimeFundReply()
	case "bitcoin":
		reply = cryptoc.GetBtcUSDReply()
	default:
		reply = "暂时无法理解： " + update.Message.Text
	}

	log.Debug("Reply to update: %+v %+v", update.Message.Text, update)
	msg := NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = ModeHTML
	//msg.ParseMode = tgbotapi.ModeMarkdown
	log.Debug("Bot Token: %+v", bot.Token)
	bot.Send(msg)

	log.Debug("[%s] %s", update.Message.From.UserName, update.Message.Text)
}

func Handler(w http.ResponseWriter, req *http.Request) {
	handler(w, req)
}
