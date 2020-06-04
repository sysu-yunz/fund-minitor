package global

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Bot *tgbotapi.BotAPI
	MgoDB *mongo.Client
)
