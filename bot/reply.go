package bot

import (
	"fmt"
	"fund/data"
	"fund/global"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tb "github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"sort"
	"strconv"
	"strings"
)

func GlobalIndexReply(update tgbotapi.Update) {
	//indices := []string{"000001.SS","399001.SZ","^GSPC","^DJI","^IXIC","^RUT","^VIX","^HSI"}
	indices := []string{"000001.SS","399001.SZ", "^IXIC","^HSI"}

	var reply [][]string
	ch := make(chan data.Meta, len(indices))

	for _, f := range indices {
		go data.IndexData(f, ch)
	}

	for range indices {
		raw := <-ch
		price := fmt.Sprintf("%.1f", raw.RegularMarketPrice)
		symbol := raw.Symbol
		change := fmt.Sprintf("%.1f", raw.RegularMarketPrice - raw.PreviousClose)
		rate := fmt.Sprintf("%.2f", (raw.RegularMarketPrice - raw.PreviousClose)/raw.PreviousClose*100)
		reply = append(reply, []string{symbol, rate, price, change})
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
	table.SetHeader([]string{"Symbol", "%", "PRICE", "/"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	TextReply(update, "<pre>"+tableString.String()+"</pre>")
	//return "```"+tableString.String()+"```"
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
			Code:  raw.Fundcode,
			Name:  raw.Name,
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
				h.TodayEarn = h.Shares * (h.Price - h.LastPrice)
				h.Cap = h.Price * h.Shares
			}
		}

		reply = append(reply, []string{
			// h.Code,
			fmt.Sprintf("%.1f", h.TodayEarn),
			fmt.Sprintf("%.1f", h.Cost*h.Shares),
			h.Name,
		})
	}

	sort.Slice(reply, func(i, j int) bool {
		iF, _ := strconv.ParseFloat(reply[i][0], 64)
		jF, _ := strconv.ParseFloat(reply[j][0], 64)
		return iF > jF
	})

	tableString := &strings.Builder{}
	table := tb.NewWriter(tableString)
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("+")
	table.SetHeader([]string{"EARN", "COST", "NAME"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	TextReply(update, "<pre>"+tableString.String()+"</pre>")
	//return "```"+tableString.String()+"```"
}

func BondReply(update tgbotapi.Update) {
	var reply [][]string
	bond := data.GetChina10YearBondYield().Records

	reply = append(reply, []string{"10 Rate", bond[0].TenRate})
	reply = append(reply, []string{"10 RateLast", bond[1].TenRate})
	reply = append(reply, []string{"1 Rate", bond[0].OneRate})
	reply = append(reply, []string{"1 RateLast", bond[1].OneRate})
	reply = append(reply, []string{"Date", bond[0].DateString})

	tableString := &strings.Builder{}
	table := tb.NewWriter(tableString)
	table.SetAlignment(tb.ALIGN_LEFT)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")
	table.SetHeader([]string{"ITEM", "VALUE"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	TextReply(update, "<pre>"+tableString.String()+"</pre>")
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

func FundUnwatch(update tgbotapi.Update)  {
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
