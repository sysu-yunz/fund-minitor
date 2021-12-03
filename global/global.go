package global

import (
	"fund/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	Bot    *tgbotapi.BotAPI
	MgoDB  *db.MgoC
	Cookie string
)
