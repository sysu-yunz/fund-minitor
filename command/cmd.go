package command

import (
	"fund/cryptoc"
	"fund/data"
	"fund/global"
	"fund/reply"
	"fund/watch"
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
				watch.FundWatch(update)
			case "unwatch":
				watch.FundUnwatch(update)
			case "fund":
				data.RealTimeFundReply(update)
			case "bitcoin":
				cryptoc.GetBtcUSDReply()
			case "index":
				data.GlobalIndexReply()
			case "bond":
				data.BondReply(update)
			default:
				reply.TextReply(update, "暂时无法理解： "+update.Message.Text)
			}
		}
	}
}
