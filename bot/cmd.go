package bot

import (
	"fund/data"
	"fund/global"
	"fund/log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tb "github.com/olekukonko/tablewriter"
)

func Handle(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		global.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
		global.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
	}

	if update.Message != nil {
		if update.Message.IsCommand() {
			log.Info("Command %s", update.Message.Text)
			switch update.Message.Command() {
			case "watch":
				FundWatch(update)
			case "unwatch":
				FundUnwatch(update)
			case "fund":
				RealTimeFundReply(update)
			case "hold":
				HoldReply(update)
			case "index":
				GlobalIndexReply(update)
			case "tv":
				TVReply(update)
			case "kpl":
				KPL(update)
			case "stock":
				StockReply(update)
			case "test":
				TestReply(update)
			default:
				StockReply(update)
			}
		} else {
			StockReply(update)
		}
	}
}

func TestReply(update tgbotapi.Update) {
	log.Info("***************%v", data.GetLagestCryptoID())
}

type RR interface {
	MsgReply(update tgbotapi.Update)
}

func TextReply(update tgbotapi.Update, s string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, s)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML
	m, err := global.Bot.Send(msg)
	if err != nil {
		log.Error("Text reply %+v ", err)
	}

	log.Debug("Replied update %+v with %+v", update.UpdateID, m)
}

func TableReply(update tgbotapi.Update, colSep string, cenSep string, reply [][]string) {
	tableString := &strings.Builder{}
	table := tb.NewWriter(tableString)
	table.SetColumnSeparator(colSep)
	table.SetCenterSeparator(cenSep)
	table.SetHeader(reply[0])

	for _, v := range reply[1:] {
		table.Append(v)
	}

	table.Render()
	TextReply(update, "<pre>"+tableString.String()+"</pre>")
	//return "```"+tableString.String()+"```"
}

func Keyboard(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("hello"),
		),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButtonContact("contact"),
			tgbotapi.NewKeyboardButtonLocation("location"),
		),
	)

	global.Bot.Send(msg)
}

func InLineKeyboard(update tgbotapi.Update) {
	d1 := tgbotapi.NewInlineKeyboardButtonData("chengqian", "hello, world")
	d2 := tgbotapi.NewInlineKeyboardButtonData("chengqian", "hello, world")
	d3 := tgbotapi.NewInlineKeyboardButtonData("chengqian", "hello, world")

	sw1 := tgbotapi.NewInlineKeyboardButtonSwitch("this is switch", "what is switch")
	sw2 := tgbotapi.NewInlineKeyboardButtonSwitch("this is switch", "what is switch")
	sw3 := tgbotapi.NewInlineKeyboardButtonSwitch("this is switch", "what is switch")
	sw4 := tgbotapi.NewInlineKeyboardButtonSwitch("this is switch", "what is switch")
	sw5 := tgbotapi.NewInlineKeyboardButtonSwitch("this is switch", "what is switch")
	sw6 := tgbotapi.NewInlineKeyboardButtonSwitch("this is switch", "what is switch")

	zhihuBtn := tgbotapi.NewInlineKeyboardButtonURL("zhihu", "https://www.zhihu.com")
	chengqianBtn := tgbotapi.NewInlineKeyboardButtonURL("chengqian", "https://duchengqian.com")
	baiduBtn := tgbotapi.NewInlineKeyboardButtonURL("baidu", "https://www.baidu.com")

	dRow := tgbotapi.NewInlineKeyboardRow(d1, d2, d3)
	swRow := tgbotapi.NewInlineKeyboardRow(sw1, sw2, sw3, sw4, sw5, sw6)
	urlRow := tgbotapi.NewInlineKeyboardRow(zhihuBtn, chengqianBtn, baiduBtn)

	a := tgbotapi.NewInlineKeyboardMarkup(dRow)
	b := tgbotapi.NewInlineKeyboardMarkup(swRow)
	c := tgbotapi.NewInlineKeyboardMarkup(urlRow)

	m := tgbotapi.NewEditMessageReplyMarkup(update.Message.Chat.ID, update.Message.MessageID, a)
	n := tgbotapi.NewEditMessageReplyMarkup(update.Message.Chat.ID, update.Message.MessageID, b)
	o := tgbotapi.NewEditMessageReplyMarkup(update.Message.Chat.ID, update.Message.MessageID, c)

	global.Bot.Send(m)
	global.Bot.Send(n)
	global.Bot.Send(o)

	//	tgbotapi.NewReplyKeyboard(urlRow)
}
