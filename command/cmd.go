package command

import (
	"fund/cryptoc"
	"fund/data"
	"fund/global"
	"fund/log"
	"fund/watch"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Handle(update tgbotapi.Update)  {
	if update.Message.IsCommand() {
		var reply string
		switch update.Message.Command() {
		case "fund_watch":
			reply = watch.FundWatch(update)
		case "fund":
			reply = data.RealTimeFundReply(update.Message.Chat.ID)
		case "bitcoin":
			reply = cryptoc.GetBtcUSDReply()
		case "index":
			reply = data.GlobalIndexReply()
		case "bond":
			reply = data.BondReply()

		default:
			reply = "暂时无法理解： " + update.Message.Text
		}

		log.Debug("Reply to update: %+v %+v", update.Message.Text, update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = tgbotapi.ModeHTML
		//msg.ParseMode = tgbotapi.ModeMarkdown
		global.Bot.Send(msg)
		log.Debug("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}
