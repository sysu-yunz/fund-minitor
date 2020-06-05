package watch

import (
	"fund/global"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func FundWatch(update tgbotapi.Update) string  {
	chatID := update.Message.Chat.ID
	arguments := update.Message.CommandArguments()
	if f, ok := global.MgoDB.ValidFundCode(arguments); ok {
		if !global.MgoDB.FundWatched(chatID, arguments) {
			global.MgoDB.InsertWatch(update.Message.Chat.ID, arguments)
			return f.FundName+"\n"+f.FundType
		} else {
			return "Fund already on watch !"
		}
	}

	return "Invalid fundCode !"
}
