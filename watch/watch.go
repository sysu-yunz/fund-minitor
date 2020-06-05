package watch

import (
	"fund/global"
	"fund/reply"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func FundWatch(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	arguments := update.Message.CommandArguments()
	if f, ok := global.MgoDB.ValidFundCode(arguments); ok {
		if !global.MgoDB.FundWatched(chatID, arguments) {
			global.MgoDB.InsertWatch(update.Message.Chat.ID, arguments)
			reply.TextReply(update, f.FundName+"\n"+f.FundType)
		} else {
			reply.TextReply(update, "Fund already watched !")
		}
	} else {
		reply.TextReply(update, "Invalid fundCode !")
	}
}

func FundUnwatch(update tgbotapi.Update)  {
	chatID := update.Message.Chat.ID
	arguments := update.Message.CommandArguments()
	if f, ok := global.MgoDB.ValidFundCode(arguments); ok {
		if global.MgoDB.FundWatched(chatID, arguments) {
			global.MgoDB.DeleteWatch(update.Message.Chat.ID, arguments)
			reply.TextReply(update, f.FundName+"\n"+f.FundType)
		} else {
			reply.TextReply(update, "Fund not on watch !")
		}
	} else {
		reply.TextReply(update, "Invalid fundCode !")
	}
}
