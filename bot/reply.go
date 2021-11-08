package bot

import (
	"fmt"
	"fund/data"
	"fund/global"
	"fund/util"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tb "github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
)

func TVReply(update tgbotapi.Update) {

	tvs := map[string]string{
		"young-sheldon":                    "小谢尔顿",
		"yellowstone-paramount-television": "黄石",
		// "ted-lasso":     "足球教练",
		// "billions":      "亿万",
	}

	res := ""
	for engName, chName := range tvs {
		res = res + chName + "\n" + data.GetTVData(engName) + "\n"
	}

	TextReply(update, res)
}

func GlobalIndexReply(update tgbotapi.Update) {
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

	bondRaw := data.GetChina10YearBondYield().Records
	bond := cast.ToFloat64(bondRaw[0].TenRate)

	bond10Str := fmt.Sprintf("%.1f", bond*100)
	rate10 := cast.ToFloat64(bondRaw[0].TenRate)
	rate10last := cast.ToFloat64(bondRaw[1].TenRate)
	ror := fmt.Sprintf("%.2f", (rate10-rate10last)/rate10last*100)
	percent := fmt.Sprintf("%.2f", (rate10-rate10last)*100)
	reply = append(reply, []string{ror, bond10Str, percent, "国债10"})

	btc := data.GetBitcoin().CoinData[0]
	q := btc.Quote.USD
	btcRow := fmt.Sprintf("%.2f, %.1f, %.1f, 比特币", q.PercentChange24H, q.Price, q.Price-q.Price/(1+q.PercentChange24H/100))

	reply = append(reply, strings.Split(btcRow, ", "))

	tableString := &strings.Builder{}
	table := tb.NewWriter(tableString)
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("+")
	table.SetHeader([]string{"%", "PRICE", "/", "NAME"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	TextReply(update, "<pre>"+tableString.String()+"</pre>")
	//return "```"+tableString.String()+"```"
}

func ChartsReply(update tgbotapi.Update) {
	//chatID := update.Message.Chat.ID

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
	p, _ := os.Create("pie.html")
	pie.Render(p)
}

func generatePieItems() []opts.PieData {
	items := make([]opts.PieData, 0)

	items = append(items, opts.PieData{Name: "主动债券基金", Value: 10000})
	items = append(items, opts.PieData{Name: "指数基金", Value: 10000})
	items = append(items, opts.PieData{Name: "主动行业基金", Value: 10000})

	return items
}

func RealTimeFundReply(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	watchFunds := global.MgoDB.GetWatchList(chatID)

	var reply [][]string

	ch := make(chan data.RealTimeRaw, len(watchFunds))

	for _, f := range watchFunds {
		go data.GetRealTime(f.Watch, ch)
	}

	for range watchFunds {
		raw := <-ch
		reply = append(reply, []string{raw.Fundcode, raw.Gszzl, raw.Name})
	}

	sort.Slice(reply, func(i, j int) bool {
		iF, _ := strconv.ParseFloat(reply[i][1], 64)
		jF, _ := strconv.ParseFloat(reply[j][1], 64)
		return iF > jF
	})

	tableString := &strings.Builder{}
	table := tb.NewWriter(tableString)
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("+")
	table.SetHeader([]string{"CODE", "RATE", "NAME"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	TextReply(update, "<pre>"+tableString.String()+"</pre>")
	//return "```"+tableString.String()+"```"
}

func HoldReply(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	holdFunds := global.MgoDB.GetHolding(chatID)

	var hr []data.HoldReply
	ch := make(chan data.RealTimeRaw, len(holdFunds.Shares))

	for _, f := range holdFunds.Shares {
		go data.GetRealTime(f.Code, ch)
	}

	for range holdFunds.Shares {
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
		for _, f := range holdFunds.Shares {
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

	btc := data.GetBitcoin().CoinData[0]
	q := btc.Quote.USD.Price * 6.5
	bc := 3438.0
	btcRow := fmt.Sprintf("%.1f, %.2f, %.1f, 比特币", q*holdFunds.Bitcoin-bc, (q*holdFunds.Bitcoin-bc)/bc*100, bc)
	reply = append(reply, strings.Split(btcRow, ", "))

	tableString := &strings.Builder{}
	table := tb.NewWriter(tableString)
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("+")
	table.SetHeader([]string{"EARN", "%", "COST", "NAME"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	TextReply(update, "<pre>"+tableString.String()+"</pre>")
	//return "```"+tableString.String()+"```"
}

func FundWatch(update tgbotapi.Update) {
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

func FundUnwatch(update tgbotapi.Update) {
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
