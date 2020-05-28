package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func Bot()  {
	fmt.Println("Hello, bot!")
	run()
}

func run()  {
	bot, err := tgbotapi.NewBotAPI("1215007591:AAEIFtHy4V4WWuxeviR2Q1V4-M-LMZTXKUw")
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

		switch update.Message.Text {
		case "fund":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "实时基金涨跌")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		//
		//bot.Send(msg)
	}
}
