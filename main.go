package main

import (
	"encoding/json"
	"fmt"
	"fund/bot"
	"fund/config"
	"fund/data"
	"fund/db"
	"fund/global"
	"fund/job"
	"fund/log"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	err error
)

func init() {
	botToken := config.EnvVariable("BOT_TOKEN")
	global.Bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Init Bot %+v", err)
	}
	global.Bot.Debug = true
	log.Info("Authorized on account %s", global.Bot.Self.UserName)

	mgoPwd := config.EnvVariable("MGO_PWD")
	global.MgoDB = db.NewDB(mgoPwd)

	data.UpdateCookie()
}

func main() {
	// local
	if len(os.Args) > 1 {
		// go job.Update()
		bot.Run()
	} else {
		port := config.EnvVariable("PORT")
		fmt.Println("Server started on port:", port)
		global.Bot.RemoveWebhook()
		_, err = global.Bot.SetWebhook(tgbotapi.NewWebhook("https://thawing-scrubland-62700.herokuapp.com/" + global.Bot.Token))
		if err != nil {
			log.Fatal("%v", err)
		}

		router := gin.New()
		router.Use(gin.Logger())
		router.Use(gin.Recovery())

		router.POST("/"+global.Bot.Token, webhookHandler)
		router.GET("/update", job.DailyReport)

		err := router.Run(":" + port)
		if err != nil {
			log.Error("%v", err)
		}
	}
}

func ChartsReply() {
	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "家庭资产配置",
	}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "",
			TriggerOn: "",
			Formatter: "{b}-{d}%",
		}))

	pie.AddSeries("xxxxxxxx", generatePieItems())
	p, err := os.Create("job/pie.html")
	if err != nil {
		fmt.Println(err)
	}

	pie.Render(p)
}

func generatePieItems() []opts.PieData {
	items := make([]opts.PieData, 0)

	items = append(items, opts.PieData{Name: "主动债券基金", Value: 10000})
	items = append(items, opts.PieData{Name: "指数基金", Value: 10000})
	items = append(items, opts.PieData{Name: "主动行业基金", Value: 10000})

	return items
}

func webhookHandler(c *gin.Context) {
	defer c.Request.Body.Close()

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("%+v", err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Error("%+v", err)
		return
	}

	// to monitor changes run: heroku logs --tail
	log.Info("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)
	bot.Handle(update)
}
