package bot

import (
	"fund/global"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Handle(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		global.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
		global.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
	}

	if update.Message != nil {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "watch":
				FundWatch(update)
			case "unwatch":
				FundUnwatch(update)
			case "fund":
				RealTimeFundReply(update)
			case "hold":
				HoldReply(update)
			case "bitcoin":
				GetBtcUSDReply(update)
			case "index":
				GlobalIndexReply(update)
			case "bond":
				BondReply(update)
			default:
				TextReply(update, "暂时无法理解： "+update.Message.Text)
			}
		}
	}
}
