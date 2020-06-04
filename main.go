package main

import (
	"fund/config"
	"fund/db"
)

func main() {
	// TODO: Email
	//e := &email.Email{
	//	To: "dukeyunz@hotmail.com",
	//	Subject: "Fund notification",
	//}

	pwd := config.ViperEnvVariable("MGO_PWD")
	db.NewDB(pwd)

	// TODO: Bot
	//bot.Bot()

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
