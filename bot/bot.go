package bot

import (
	"fmt"
	"fund/command"
	"fund/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func Bot()  {
	fmt.Println("Hello, bot!")
	run()
}

func run()  {
	tgBotToken := config.ViperEnvVariable("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Panic()
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		command.Command(bot, update)
	}
}
