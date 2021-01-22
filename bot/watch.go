package bot

import (
	"fund/global"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func FundWatch(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	arguments := update.Message.CommandArguments()
	if f, ok := global.MgoDB.ValidFundCode(arguments); ok {
		if !global.MgoDB.FundWatched(chatID, arguments) {
			global.MgoDB.InsertWatch(update.Message.Chat.ID, arguments)
			TextReply(update, f.FundName+"\n"+f.FundType)
		} else {
			TextReply(update, "Fund already watched !")
		}
	} else {
		TextReply(update, "Invalid fundCode !")
	}
}

func FundUnwatch(update tgbotapi.Update)  {
	chatID := update.Message.Chat.ID
	arguments := update.Message.CommandArguments()
	if f, ok := global.MgoDB.ValidFundCode(arguments); ok {
		if global.MgoDB.FundWatched(chatID, arguments) {
			global.MgoDB.DeleteWatch(update.Message.Chat.ID, arguments)
			TextReply(update, f.FundName+"\n"+f.FundType)
		} else {
			TextReply(update, "Fund not on watch !")
		}
	} else {
		TextReply(update, "Invalid fundCode !")
	}
}
