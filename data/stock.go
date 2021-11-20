package data

import (
	"encoding/json"
	"fmt"
	"fund/log"
	"io/ioutil"
	"net/http"
)

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
		log.Error("Read stock data error %v", err)
	}

	return d
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
