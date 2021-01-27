package bot

import (
	"fund/global"
	"fund/log"
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
			case "index":
				GlobalIndexReply(update)
			case "chart":
				ChartsReply(update)
			default:
				TextReply(update, "暂时无法理解： "+update.Message.Text)
			}
		}
	}
}

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

