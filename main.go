package main

import (
	"fund/bot"
	"fund/config"
	"fund/db"
	"fund/global"
	"fund/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	err error
)

func init() {
	botToken := config.ViperEnvVariable("BOT_TOKEN")
	global.Bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Init Bot %+v", err)
	}
	global.Bot.Debug = true
	log.Info("Authorized on account %s", global.Bot.Self.UserName)

	mgoPwd := config.ViperEnvVariable("MGO_PWD")
	global.MgoDB = db.NewDB(mgoPwd)
}

func main() {
	// TODO: Email
	//e := &email.Email{
	//	To: "dukeyunz@hotmail.com",
	//	Subject: "Fund notification",
	//}

	// TODO: Bot
	bot.Run()

	// TODO: Monthly history
	//fundCode := "161716"
	//monthProfitRate := profit.MonthProfitRate(fundCode)
	//if monthProfitRate > 0.01 {
	//	e.Send(fmt.Sprintf("hello, enough profit, you need to sell, last 20 workday profit rate: %.2f%% !", monthProfitRate*100.0))
	//} else if monthProfitRate < 0.03 {
	//	e.Send(fmt.Sprintf("hello, it's good time to buy, last 20 workday profit rate: %.2f%% !", monthProfitRate*100.0))
	//}

	// TODO: Real time fund data
	//watchFunds := config.GetWatches()
	//for _, f := range watchFunds {
	//	data.GetRealTime(f)
	//}

	// TODO: BTC Price
	//cryptoc.GetBtcUSD()

	// TODO: Prepare money for upcoming repayment
}
