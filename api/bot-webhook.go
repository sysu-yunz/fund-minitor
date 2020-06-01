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

var bot *BotAPI

func init()  {
	bot, err := NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("Init Bot error: ", err)
	}

	bot.Debug = true
	log.Debug("Authorized on account %s", bot.Self.UserName)
}

func handler(w http.ResponseWriter, r *http.Request) UpdatesChannel {
	ch := make(chan Update, BotAPI{}.Buffer)
	bytes, _ := ioutil.ReadAll(r.Body)
	var update Update
	json.Unmarshal(bytes, &update)
	log.Debug("Update: ", update)
	ch <- update

	return ch
}

func Handler(w http.ResponseWriter, req *http.Request) {
	updates := handler(w, req)
	go func() {
		for update := range updates {
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			var reply string
			switch update.Message.Text {
			case "fund":
				reply = data.RealTimeFundReply()
			case "bitcoin":
				reply = cryptoc.GetBtcUSDReply()
			default:
				reply = "暂时无法理解： "+update.Message.Text
			}

			msg := NewMessage(update.Message.Chat.ID, reply)
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ParseMode = ModeHTML
			//msg.ParseMode = tgbotapi.ModeMarkdown
			bot.Send(msg)

			log.Debug("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
	}()
}
