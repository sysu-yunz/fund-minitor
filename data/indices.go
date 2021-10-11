package data

import (
	"encoding/json"
	"fmt"
	"fund/log"
	"io/ioutil"
	"net/http"
)

// 000001.SS
// 399001.SZ
// ^GSPC
// ^DJI
// ^IXIC
// ^RUT
// ^VIX
// ^HSI

func IndexData(indexCode string, ch chan Meta) {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/" + indexCode + "?range=2m"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Cookie", "B=8t0kdm1fdretn&b=3&s=lf")

	res, err := client.Do(req)
	if err != nil {
		log.Error("Get indices error", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Read indices error", err)
	}

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
