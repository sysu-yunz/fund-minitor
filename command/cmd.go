package command

import (
	"fund/cryptoc"
	"fund/data"
	"fund/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Command(bot *tgbotapi.BotAPI, update tgbotapi.Update)  {
	var reply string
	switch update.Message.Text {
	case "fund":
		reply = data.RealTimeFundReply()
	case "bitcoin":
		reply = cryptoc.GetBtcUSDReply()
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

	log.Debug("[%s] %s", update.Message.From.UserName, update.Message.Text)
}
