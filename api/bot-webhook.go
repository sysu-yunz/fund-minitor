package handler

import (
	"fund/cryptoc"
	"fund/data"
	"fund/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"os"
)

type ReqBody struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Type      string `json:"type"`
}

func Handler(w http.ResponseWriter, req *http.Request) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal("Init Bot error: ", err)
	}

	bot.Debug = true

	log.Debug("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/" + bot.Token)
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

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = tgbotapi.ModeHTML
		//msg.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(msg)

		log.Debug("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}
