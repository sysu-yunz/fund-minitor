package data

import (
	"encoding/json"
	"fmt"
	"fund/log"
	r "fund/reply"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tb "github.com/olekukonko/tablewriter"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func GlobalIndexReply(update tgbotapi.Update) {
	//indices := []string{"000001.SS","399001.SZ","^GSPC","^DJI","^IXIC","^RUT","^VIX","^HSI"}
	indices := []string{"000001.SS","399001.SZ", "^IXIC","^HSI"}

	var reply [][]string
	ch := make(chan Meta, len(indices))

	for _, f := range indices {
		go indexData(f, ch)
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
	table.SetHeader([]string{"Symbol", "%", "PRICE", "CHANGE"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	r.TextReply(update, "<pre>"+tableString.String()+"</pre>")
	//return "```"+tableString.String()+"```"
}

// 000001.SS
// 399001.SZ
// ^GSPC
// ^DJI
// ^IXIC
// ^RUT
// ^VIX
// ^HSI

func indexData(indexCode string, ch chan Meta) {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/"+indexCode+"?range=2m"
	method := "GET"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Cookie", "B=8t0kdm1fdretn&b=3&s=lf")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	//fmt.Println(string(body))

	raw := IndexChartRaw{}
	err = json.Unmarshal(body, &raw)
	if err != nil {
		log.Error("Indices Unmarshal %+v", err)
	}

	ch <- raw.Chart.Result[0].MetaData
}

type IndexChartRaw struct {
	Chart Chart `json:"chart"`
}
type Meta struct {
	Currency             string  `json:"currency"`
	Symbol               string  `json:"symbol"`
	ExchangeName         string  `json:"exchangeName"`
	InstrumentType       string  `json:"instrumentType"`
	FirstTradeDate       int     `json:"firstTradeDate"`
	RegularMarketTime    int     `json:"regularMarketTime"`
	Timezone             string  `json:"timezone"`
	ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
	RegularMarketPrice   float64 `json:"regularMarketPrice"`
	ChartPreviousClose   float64 `json:"chartPreviousClose"`
	PreviousClose        float64 `json:"previousClose"`
	Scale                int     `json:"scale"`
	PriceHint            int     `json:"priceHint"`
}
type Result struct {
	MetaData Meta `json:"meta"`
}
type Chart struct {
	Result []Result `json:"result"`
}