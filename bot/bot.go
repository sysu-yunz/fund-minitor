package bot

import (
	"fund/global"
	"fund/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Run()  {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := global.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Error("Get Updates Channel %+v", err)
	}

	for update := range updates {
		Handle(update)
	}
}

