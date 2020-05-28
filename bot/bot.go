package bot

import (
	"fmt"
	"fund/config"
	"fund/cryptoc"
	"fund/data"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func Bot()  {
	fmt.Println("Hello, bot!")
	run()
}

func run()  {
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotAPIToken)
	if err != nil {
		log.Panic()
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		var reply string
		switch update.Message.Text {
		case "fund":
			reply = data.RealTimeFundReply()
		case "bitcoin":
			reply = cryptoc.GetBtcUSDReply()
		default:
			reply = "暂时无法理解： "+update.Message.Text
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ParseMode = tgbotapi.ModeHTML
		//msg.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(msg)

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}
}
