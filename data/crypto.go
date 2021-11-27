package data

import (
	"encoding/json"
	"fmt"
	"fund/db"
	"fund/global"
	"fund/log"
	"io/ioutil"
	"net/http"
)

func GetCoinQuote(id string) db.CoinQuoteRaw {
	url := "https://web-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?id=" + id
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Cookie", "__cfduid=d90bb5e0c15aa0f9f956eca4afd34692e1590303661")

	res, err := client.Do(req)
	if err != nil {
		log.Error("Get crypto error", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Read bitcoin error", err)
	}

	quote := db.CoinQuoteRaw{}
	err = json.Unmarshal(body, &quote)
	if err != nil {
		log.Error("JSON unmarshal error: ", err)
	}

	return quote
}

func GetCryptoCount() int64 {
	// return doc count in crypto collection
	return global.MgoDB.GetCryptoCount()
}

func GetStockCount(market string) int64 {
	// return doc count in stock collection
	return global.MgoDB.GetStockCount(market)
}

func UpdateCoinList() {
	start := 1
	limit := 200
	client := &http.Client{}

	global.MgoDB.DeleteCryptoList()

	// request data until all crypto are fetched
	for {
		// format url with start and limit
		url := fmt.Sprintf("https://web-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?start=%d&limit=%d", start, limit)
		method := "GET"
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Cookie", "__cfduid=d90bb5e0c15aa0f9f956eca4afd34692e1590303661")

		res, err := client.Do(req)
		if err != nil {
			log.Error("Get bitcoin error", err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error("Read bitcoin error", err)
		}

		//fmt.Println(string(body))

		coinList := db.CoinsRaw{}
		err = json.Unmarshal(body, &coinList)
		if err != nil {
			log.Error("JSON unmarshal error: ", err)
		}

		global.MgoDB.InsertCryptoList(coinList.Coins)

		log.Info("********************** %v *************** %v", coinList.Status.TotalCount, start+limit)

		if start > coinList.Status.TotalCount {
			break
		}

		start += limit
	}

}
