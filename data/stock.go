package data

import (
	"encoding/json"
	"fmt"
	"fund/db"
	"fund/global"
	"fund/log"
	"io/ioutil"
	"net/http"
	"strings"
)

// pass a hint to get stock symbol
// hint can be either stock name or stock code like "万科A" or "sz002081"
func GetSymbol(arg string) string {
	if res := global.MgoDB.SearchStock(arg, false); res != "" {
		return res
	}

	return global.MgoDB.SearchStock(arg, true)
}

func GetStock(code string) RealTimeStockData {
	// call stock api with code
	// return stock data
	c := &http.Client{}
	url := fmt.Sprintf("https://stock.xueqiu.com/v5/stock/quote.json?symbol=%s&extend=detail", code)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("make request error %v", err)
	}
	req.Header.Add("Accept", " */*")
	req.Header.Add("Origin", " https://xueqiu.com")
	req.Header.Add("Cookie", "xq_a_token=ad254175b8f79f3ce1be51812b24adb083dc9851 ; s=c0159r1h9d")
	// req.Header.Add("Accept-Encoding", " gzip, deflate, br")
	req.Header.Add("Host", " stock.xueqiu.com")
	req.Header.Add("User-Agent", " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15")
	req.Header.Add("Accept-Language", " zh-CN,zh-Hans;q=0.9")
	req.Header.Add("Connection", " keep-alive")

	resp, err := c.Do(req)
	if err != nil {
		log.Error("Init client error %v", err)
	}

	defer resp.Body.Close()

	// return stock data
	// parse resp body to RealtimeStockData
	var d RealTimeStockData
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Read stock error", err)
	}

	err = json.Unmarshal(body, &d)
	if err != nil {
		log.Info(string(body))
		log.Error("Read stock data error %v", err)
	}

	return d
}

func UpdateStockList(market Market) {
	page := 1
	size := 200
	client := &http.Client{}

	global.MgoDB.DeleteStockList(strings.ToLower(market.Country))

	// request data until all stocks are fetched
	for {
		url := fmt.Sprintf("https://xueqiu.com/service/v5/stock/screener/quote/list?page=%d&size=%d&order=desc&orderby=percent&order_by=percent&market=%s&type=%s", page, size, market.Country, market.Board)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error("request init error %v", err)
		}

		req.Header.Add("Cookie", "acw_tc=2760829d16374557306416084eb60952ba53cd7be6edbdb9dbed52354505c6; s=c0159r1h9d; xq_a_token=ad254175b8f79f3ce1be51812b24adb083dc9851; xq_id_token=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ1aWQiOi0xLCJpc3MiOiJ1YyIsImV4cCI6MTYzOTc2NDY2MCwiY3RtIjoxNjM3NDU1NzEwMDQwLCJjaWQiOiJkOWQwbjRBWnVwIn0.h50zlh-D8e7bAckr1HYbAn8eXyQ30p1Q4xGP1Uvu8F0FtTnEwbketJh-ioaj_RonipQyue_Eu4rQI26cOYB6dfWRcMhSmeieGgQ5723y9lcbyEqAIF5WJ25gEUgmEBXcPmRzCKW1VlFHiQe4kBM3HAhAnOHz0dt50d24ccKGP3cfE7NRjjWWdv1NBm0ch3pKJ_XtYV9epOPKA-fqUuekOfukwOQJeT-jhAKs93EY0yNhRjkfqMGbaiZqEx9D0R1t1eKIusZ6zcMqbF5TaWFzlAnUVAiaGoOxBlC-HWLZhfrBa_OFa5FTg3KoDG534rVet6JUxCoqOfoOw36jcSJPpA; xq_r_token=55944e6d0310d70bf0039e421a9a722032a84077")

		res, err := client.Do(req)
		if err != nil {
			log.Error("request error %v", err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error("read body error %v", err)
		}

		stocks := db.StockList{}
		err = json.Unmarshal(body, &stocks)
		if err != nil {
			log.Info(string(body))
			log.Error("unmarshal error %v", err)
		}

		global.MgoDB.InsertStockList(stocks, strings.ToLower(market.Country))

		log.Info("%v*********** %v ********** %v", market.Country, stocks.Data.Count, page*size)

		if page*size > stocks.Data.Count {
			break
		}

		page = page + 1
	}
}

type Market struct {
	Country string `json:"market"`
	Board   string `json:"board"`
}

type RealTimeStockData struct {
	Data struct {
		Market struct {
			StatusID     int         `json:"status_id"`
			Region       string      `json:"region"`
			Status       string      `json:"status"`
			TimeZone     string      `json:"time_zone"`
			TimeZoneDesc interface{} `json:"time_zone_desc"`
			DelayTag     int         `json:"delay_tag"`
		} `json:"market"`
		Quote struct {
			CurrentExt               interface{} `json:"current_ext"`
			Symbol                   string      `json:"symbol"`
			VolumeExt                interface{} `json:"volume_ext"`
			High52W                  float64     `json:"high52w"`
			Delayed                  int         `json:"delayed"`
			Type                     int         `json:"type"`
			TickSize                 float64     `json:"tick_size"`
			FloatShares              int64       `json:"float_shares"`
			LimitDown                float64     `json:"limit_down"`
			NoProfit                 interface{} `json:"no_profit"`
			High                     float64     `json:"high"`
			FloatMarketCapital       int64       `json:"float_market_capital"`
			TimestampExt             interface{} `json:"timestamp_ext"`
			LotSize                  int         `json:"lot_size"`
			LockSet                  interface{} `json:"lock_set"`
			WeightedVotingRights     interface{} `json:"weighted_voting_rights"`
			Chg                      float64     `json:"chg"`
			Eps                      float64     `json:"eps"`
			LastClose                float64     `json:"last_close"`
			ProfitFour               float64     `json:"profit_four"`
			Volume                   int         `json:"volume"`
			VolumeRatio              float64     `json:"volume_ratio"`
			ProfitForecast           int64       `json:"profit_forecast"`
			TurnoverRate             float64     `json:"turnover_rate"`
			Low52W                   float64     `json:"low52w"`
			Name                     string      `json:"name"`
			Exchange                 string      `json:"exchange"`
			PeForecast               float64     `json:"pe_forecast"`
			TotalShares              int64       `json:"total_shares"`
			Status                   int         `json:"status"`
			IsVieDesc                interface{} `json:"is_vie_desc"`
			SecurityStatus           interface{} `json:"security_status"`
			Code                     string      `json:"code"`
			GoodwillInNetAssets      float64     `json:"goodwill_in_net_assets"`
			AvgPrice                 float64     `json:"avg_price"`
			Percent                  float64     `json:"percent"`
			WeightedVotingRightsDesc interface{} `json:"weighted_voting_rights_desc"`
			Amplitude                float64     `json:"amplitude"`
			Current                  float64     `json:"current"`
			IsVie                    interface{} `json:"is_vie"`
			CurrentYearPercent       float64     `json:"current_year_percent"`
			IssueDate                int64       `json:"issue_date"`
			SubType                  string      `json:"sub_type"`
			Low                      float64     `json:"low"`
			IsRegistrationDesc       interface{} `json:"is_registration_desc"`
			NoProfitDesc             interface{} `json:"no_profit_desc"`
			MarketCapital            int64       `json:"market_capital"`
			Dividend                 float64     `json:"dividend"`
			DividendYield            float64     `json:"dividend_yield"`
			Currency                 string      `json:"currency"`
			Navps                    float64     `json:"navps"`
			Profit                   float64     `json:"profit"`
			Timestamp                int64       `json:"timestamp"`
			PeLyr                    float64     `json:"pe_lyr"`
			Amount                   float64     `json:"amount"`
			PledgeRatio              float64     `json:"pledge_ratio"`
			TradedAmountExt          interface{} `json:"traded_amount_ext"`
			IsRegistration           interface{} `json:"is_registration"`
			Pb                       float64     `json:"pb"`
			LimitUp                  float64     `json:"limit_up"`
			PeTtm                    float64     `json:"pe_ttm"`
			Time                     int64       `json:"time"`
			Open                     float64     `json:"open"`
		} `json:"quote"`
		Others struct {
			PankouRatio float64 `json:"pankou_ratio"`
			CybSwitch   bool    `json:"cyb_switch"`
		} `json:"others"`
		Tags []struct {
			Description string `json:"description"`
			Value       int    `json:"value"`
		} `json:"tags"`
	} `json:"data"`
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}
