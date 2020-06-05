package data

import (
	"fund/reply"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GlobalIndexReply(update tgbotapi.Update) {
	reply.TextReply(update,  "Global Index Not Implemented!")
}
