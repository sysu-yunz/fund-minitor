package command

import (
	"fund/cryptoc"
	"fund/data"
	"fund/log"
	"fund/watch"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Command(bot *tgbotapi.BotAPI, update tgbotapi.Update)  {
	var reply string
	switch update.Message.Text {
	case "/fund_watch":
		reply = watch.FundWatch()
	case "/fund":
		reply = data.RealTimeFundReply()
	case "/bitcoin":
		reply = cryptoc.GetBtcUSDReply()
	case "/index":
		reply = data.GlobalIndexReply()
	case "/bond":
		reply = data.BondReply()

	default:
		reply = "暂时无法理解： " + update.Message.Text
	}

	log.Debug("Reply to update: %+v %+v", update.Message.Text, update)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML
	//msg.ParseMode = tgbotapi.ModeMarkdown
	log.Debug("Bot Token: %+v", bot.Token)
	bot.Send(msg)


	msg2 := tgbotapi.NewMessage(update.Message.Chat.ID+1, "this is test msg")
	bot.Send(msg2)

	log.Debug("[%s] %s", update.Message.From.UserName, update.Message.Text)
}
