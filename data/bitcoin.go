package data

import (
	"encoding/json"
	"fmt"
	"fund/log"
	"io/ioutil"
	"net/http"
	"time"
)

func GetBitcoin() BitcoinRaw {
	url := "https://web-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?limit=6&start=1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Cookie", "__cfduid=d90bb5e0c15aa0f9f956eca4afd34692e1590303661")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	//fmt.Println(string(body))

	btc := BitcoinRaw{}
	err = json.Unmarshal(body, &btc)
	if err != nil {
		log.Error("JSON unmarshal error: ", err)
	}

	return btc
}

type BitcoinRaw struct {
	Status   Status     `json:"status"`
	CoinData []CoinData `json:"data"`
}

type Status struct {
	Timestamp    time.Time   `json:"timestamp"`
	ErrorCode    int         `json:"error_code"`
	ErrorMessage interface{} `json:"error_message"`
	Elapsed      int         `json:"elapsed"`
	CreditCount  int         `json:"credit_count"`
	Notice       interface{} `json:"notice"`
}

type USD struct {
	Price            float64   `json:"price"`
	Volume24H        float64   `json:"volume_24h"`
	PercentChange1H  float64   `json:"percent_change_1h"`
	PercentChange24H float64   `json:"percent_change_24h"`
	PercentChange7D  float64   `json:"percent_change_7d"`
	MarketCap        float64   `json:"market_cap"`
	LastUpdated      time.Time `json:"last_updated"`
}

type Quote struct {
	USD USD `json:"USD"`
}

type CoinData struct {
	ID                int         `json:"id"`
	Name              string      `json:"name"`
	Symbol            string      `json:"symbol"`
	Slug              string      `json:"slug"`
	NumMarketPairs    int         `json:"num_market_pairs"`
	DateAdded         time.Time   `json:"date_added"`
	Tags              []string    `json:"tags"`
	MaxSupply         int         `json:"max_supply"`
	CirculatingSupply int         `json:"circulating_supply"`
	TotalSupply       int         `json:"total_supply"`
	Platform          interface{} `json:"platform"`
	CmcRank           int         `json:"cmc_rank"`
	LastUpdated       time.Time   `json:"last_updated"`
	Quote             Quote       `json:"quote"`
}

