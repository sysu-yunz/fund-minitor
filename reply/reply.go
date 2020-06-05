package reply

import (
	"fund/global"
	"fund/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type RR interface {
	MsgReply(update tgbotapi.Update, )
}

func TextReply(update tgbotapi.Update, s string)  {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, s)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML
	m, err := global.Bot.Send(msg)
	if err != nil {
		log.Error("Text reply %+v ", err)
	}

	log.Debug("Replied update %+v with %+v", update.UpdateID, m)
}