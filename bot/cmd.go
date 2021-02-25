package bot

import (
	"fmt"
	"fund/global"
	"fund/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/piquette/finance-go/quote"
)

func Handle(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		global.Bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
		global.Bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
	}

	if update.Message != nil {
		if update.Message.IsCommand() {
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
			case "chart":
				ChartsReply(update)
			case "test":
				Keyboard(update)
			case "appl":
				Yahoo(update)

			// TODO
			//case "buy":
			//	BuyReply(update)
			//case "sell":
			//	SellReply(update)
			//case "stock":
			//	StockReply(update)

			default:
				TextReply(update, "暂时无法理解： "+update.Message.Text)
			}
		}
	}
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

func Yahoo(update tgbotapi.Update)  {
	q, err := quote.Get("AAPL")
	if err != nil {
		// Uh-oh.
		panic(err)
	}

	// Success!
	fmt.Println(q)
	TextReply(update, fmt.Sprintf("%f", q.Ask))
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
