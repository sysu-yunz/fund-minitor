package bot

import (
	"fmt"
	"fund/data"
	"fund/db"
	"fund/global"
	"fund/util"
	"sort"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cast"
)

func TVReply(update tgbotapi.Update) {

	tvs := map[string]string{
		"young-sheldon": "小谢尔顿",
		// "yellowstone-paramount-television": "黄石",
		"y-1883":      "黄石1833",
		"the-expanse": "浩瀚苍穹",
		// "ted-lasso":     "足球教练",
		// "billions":      "亿万",
	}

	res := ""
	for engName, chName := range tvs {
		res = res + chName + "\n" + data.GetTVData(engName) + "\n"
	}

	TextReply(update, res)
}

func Indices(update tgbotapi.Update) {
	type index struct {
		Symbol string
		Name   string
	}
	//indices := []string{"000001.SS","399001.SZ","^GSPC","^DJI","^IXIC","^RUT","^VIX","^HSI"}
	indices := []index{
		{"000001.SS", "上证"},
		{"399001.SZ", "深证"},
		{"000300.SS", "沪深300"},
		{"^HSI", "恒指"},
		{"^IXIC", "纳指"},
	}

	var reply [][]string
	ch := make(chan data.Meta, len(indices))

	for _, f := range indices {
		go data.IndexData(f.Symbol, ch)
	}

	for range indices {
		raw := <-ch
		price := fmt.Sprintf("%.1f", raw.RegularMarketPrice)
		symbol := raw.Symbol

		change := fmt.Sprintf("%.1f", raw.RegularMarketPrice-raw.PreviousClose)
		rate := fmt.Sprintf("%.2f", (raw.RegularMarketPrice-raw.PreviousClose)/raw.PreviousClose*100)

		name := symbol
		for _, d := range indices {
			if symbol == d.Symbol {
				name = d.Name
			}
		}

		reply = append(reply, []string{rate, price, change, name})
	}

	sort.Slice(reply, func(i, j int) bool {
		iF, _ := strconv.ParseFloat(reply[i][0], 64)
		jF, _ := strconv.ParseFloat(reply[j][0], 64)

		return iF > jF
	})

	// chinamoney.com.cn is not working, so remove bond index

	// bondRaw := data.GetChina10YearBondYield().Records
	// bond := cast.ToFloat64(bondRaw[0].TenRate)

	// bond10Str := fmt.Sprintf("%.1f", bond*100)
	// rate10 := cast.ToFloat64(bondRaw[0].TenRate)
	// rate10last := cast.ToFloat64(bondRaw[1].TenRate)
	// ror := fmt.Sprintf("%.2f", (rate10-rate10last)/rate10last*100)
	// percent := fmt.Sprintf("%.2f", (rate10-rate10last)*100)
	// reply = append(reply, []string{ror, bond10Str, percent, "国债10"})

	btc := data.GetCoinQuote("1").QData["1"]
	q := btc.Quote.USD
	btcRow := fmt.Sprintf("%.2f, %.1f, %.1f, 比特币", q.PercentChange24H, q.Price, q.Price-q.Price/(1+q.PercentChange24H/100))

	reply = append(reply, strings.Split(btcRow, ", "))

	// prepend header
	reply = append([][]string{{"%", "PRICE", "/", "NAME"}}, reply...)
	TableReply(update, " ", "+", reply)
}

func Subscription(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	watchFunds := global.MgoDB.GetWatchList(chatID)

	var reply [][]string

	ch := make(chan data.RealTimeRaw, len(watchFunds))

	for _, f := range watchFunds {
		go data.GetRealTime(f.Watch, ch)
	}

	for range watchFunds {
		raw := <-ch
		reply = append(reply, []string{raw.Fundcode, raw.Gszzl, util.ShortenFundName(raw.Name)})
	}

	sort.Slice(reply, func(i, j int) bool {
		iF, _ := strconv.ParseFloat(reply[i][1], 64)
		jF, _ := strconv.ParseFloat(reply[j][1], 64)
		return iF > jF
	})

	// prepend header
	reply = append([][]string{{"CODE", "RATE", "NAME"}}, reply...)
	TableReply(update, " ", "+", reply)
}

func Portfolio(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	holds := global.MgoDB.GetHolding(chatID)

	var reply [][]string
	reply = append(reply, []string{"EARN", "%", "COST", "NAME"})

	reply = append(reply, FundsHoldReply(holds.Shares)...)
	reply = append(reply, StocksHoldReply(holds.Stocks)...)
	reply = append(reply, CryptosHoldReply(holds.Cryptos)...)

	TableReply(update, " ", "+", reply)
}

func FundsHoldReply(shares []db.Share) [][]string {
	var hr []data.HoldReply

	ch := make(chan data.RealTimeRaw, len(shares))

	for _, f := range shares {
		go data.GetRealTime(f.Code, ch)
	}

	for range shares {
		raw := <-ch
		estimateValue := cast.ToFloat64(raw.Gsz)
		estimateRate := cast.ToFloat64(raw.Gszzl)

		hr = append(hr, data.HoldReply{
			Code: raw.Fundcode,
			Name: util.ShortenFundName(raw.Name),
			//Name: raw.Name,
			Rate:  estimateRate,
			Price: estimateValue,
		})
	}

	var reply [][]string

	for _, h := range hr {
		for _, f := range shares {
			if h.Code == f.Code {
				h.Cost = f.Cost
				h.Shares = f.Shares
				h.LastPrice = data.LastPrice(f.Code)
				h.Earn = h.Shares * (h.Price - h.Cost)
				h.TodayEarn = h.Shares * h.LastPrice * h.Rate / 100
				h.Cap = h.Price * h.Shares
			}
		}

		reply = append(reply, []string{
			// h.Code,
			fmt.Sprintf("%.1f", h.TodayEarn),
			fmt.Sprintf("%.2f", h.Rate),
			fmt.Sprintf("%.1f", h.Cost*h.Shares),
			h.Name,
		})
	}

	sort.Slice(reply, func(i, j int) bool {
		iF, _ := strconv.ParseFloat(reply[i][0], 64)
		jF, _ := strconv.ParseFloat(reply[j][0], 64)
		return iF > jF
	})

	return reply
}

func CryptosHoldReply(cryptos []db.CryptoHolding) [][]string {
	var reply [][]string
	for _, c := range cryptos {
		q := data.GetCoinQuote(c.ID).QData[c.ID].Quote
		// btcRow := fmt.Sprintf("%.1f, %.2f, %.1f, 比特币", q*holds.Bitcoin-bc, (q*holds.Bitcoin-bc)/bc*100, bc)
		row := fmt.Sprintf("%.1f, %.2f, %.1f, %s", q.USD.Price*c.Amount*6.4-c.Cost, (q.USD.Price*c.Amount*6.4-c.Cost)/c.Cost*100, c.Cost, c.Name)
		reply = append(reply, strings.Split(row, ", "))
	}

	return reply
}

func StocksHoldReply(stocks []db.StockHolding) [][]string {
	var reply [][]string
	for _, s := range stocks {
		q := data.GetStock(s.Symbol).Data.Quote
		row := fmt.Sprintf("%.1f, %.2f, %.1f, %s", q.Current*s.Amount-s.Cost, (q.Current*s.Amount-s.Cost)/s.Cost*100, s.Cost, s.Name)
		reply = append(reply, strings.Split(row, ", "))
	}

	return reply
}

func Sub(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	arguments := update.Message.CommandArguments()
	if f, ok := global.MgoDB.ValidFundCode(arguments); ok {
		if !global.MgoDB.FundWatched(chatID, arguments) {
			global.MgoDB.InsertWatch(update.Message.Chat.ID, arguments)
			TextReply(update, f.FundName+"\n"+f.FundType)
		} else {
			TextReply(update, "Fund already watched !")
		}
	} else {
		TextReply(update, "Invalid fundCode !")
	}
}

func Unsub(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	arguments := update.Message.CommandArguments()
	if f, ok := global.MgoDB.ValidFundCode(arguments); ok {
		if global.MgoDB.FundWatched(chatID, arguments) {
			global.MgoDB.DeleteWatch(update.Message.Chat.ID, arguments)
			TextReply(update, f.FundName+"\n"+f.FundType)
		} else {
			TextReply(update, "Fund not on watch !")
		}
	} else {
		TextReply(update, "Invalid fundCode !")
	}
}

func Quote(update tgbotapi.Update) {
	arguments := update.Message.CommandArguments()
	if arguments == "" {
		arguments = update.Message.Text
		// remove "/" in arguments
		arguments = strings.Replace(arguments, "/", "", -1)
		// change arguments to upper case
		arguments = strings.ToUpper(arguments)
	}

	symbol := data.GetSymbol(arguments)

	if symbol == "" {
		TextReply(update, "Invalid stock symbol !")
		return
	}

	stockData := data.GetStock(symbol)

	var reply [][]string
	reply = append(reply, []string{"Price", "%", "Code", "NAME"})

	q := stockData.Data.Quote

	stockReplyString := fmt.Sprintf("%.2f, %.2f, %v, %v", q.Current, q.Percent, q.Symbol, q.Name)
	reply = append(reply, strings.Split(stockReplyString, ", "))

	TableReply(update, " ", "+", reply)
}

func KPL(update tgbotapi.Update) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://yunz-rss.vercel.app/weibo/user/6074356560")
	rpl := "还未更新"
	for _, it := range feed.Items {
		if strings.Contains(it.Title, "首发名单") {
			fmt.Println(it.Link)
			rpl = it.Link
		}
	}

	TextReply(update, rpl)
}
