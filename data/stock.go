package data

import (
	"context"
	"encoding/json"
	"fmt"
	"fund/db"
	"fund/global"
	"fund/log"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
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
	req.Header.Add("Cookie", global.MgoDB.GetCookie())
	// req.Header.Add("Accept-Encoding", " gzip, deflate, br")
	req.Header.Add("Host", " stock.xueqiu.com")
	req.Header.Add("User-Agent", " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15")
	req.Header.Add("Accept-Language", " zh-CN,zh-Hans;q=0.9")
	// req.Header.Add("Connection", " keep-alive")

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
	baseURL := "https://xueqiu.com/service/v5/stock/screener/quote/list?order=desc&orderby=percent&order_by=percent&"

	// request data until all stocks are fetched
	for {
		url := fmt.Sprintf(baseURL+"page=%d&size=%d&market=%s&type=%s", page, size, market.Country, market.Board)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error("request init error %v", err)
		}

		req.Header.Add("Cookie", global.MgoDB.GetCookie())

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

// update cookie and store it in global.Cookie
func UpdateCookie() {
	url := "https://xueqiu.com/S/SZ000002"
	getCookie(url, nil)
}

func getCookie(url string, wait interface{}) {
	chromeBin := os.Getenv("GOOGLE_CHROME_SHIM")
	log.Info("chrome path: %+v", chromeBin)

	options := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath(chromeBin),
		chromedp.Flag("headless", true),
		chromedp.Flag("blink-settings", "imageEnable=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko)`),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	chromeCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Info))
	defer cancel()

	timeOutCtx, cancel := context.WithTimeout(chromeCtx, 240*time.Second)
	defer cancel()

	// listenForNetworkEvent(chromeCtx)

	log.Info("chrome visit page %s\n", url)
	err := chromedp.Run(timeOutCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		saveCookies(),
	)

	if err != nil {
		log.Error("chrome error %v", err)
	}
}

func saveCookies() chromedp.ActionFunc {
	return func(c context.Context) error {
		cookies, err := network.GetAllCookies().Do(c)
		if err != nil {
			return err
		}

		cookiesStr := ""

		for _, cookie := range cookies {
			cookiesStr = cookiesStr + cookie.Name + "=" + cookie.Value + ";"
		}

		global.MgoDB.InsertCookie(cookiesStr)

		return nil
	}

}

func listenForNetworkEvent(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {

		case *network.EventResponseReceived:
			resp := ev.Response
			if strings.HasPrefix(resp.URL, "https://stock.xueqiu.com/v5/stock/quote.json") {
				log.Info("Network reply receeived")
			}
		}
	})
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
			FloatShares              float64     `json:"float_shares"`
			LimitDown                float64     `json:"limit_down"`
			NoProfit                 interface{} `json:"no_profit"`
			High                     float64     `json:"high"`
			FloatMarketCapital       float64     `json:"float_market_capital"`
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
			ProfitForecast           float64     `json:"profit_forecast"`
			TurnoverRate             float64     `json:"turnover_rate"`
			Low52W                   float64     `json:"low52w"`
			Name                     string      `json:"name"`
			Exchange                 string      `json:"exchange"`
			PeForecast               float64     `json:"pe_forecast"`
			TotalShares              float64     `json:"total_shares"`
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
			IssueDate                float64     `json:"issue_date"`
			SubType                  string      `json:"sub_type"`
			Low                      float64     `json:"low"`
			IsRegistrationDesc       interface{} `json:"is_registration_desc"`
			NoProfitDesc             interface{} `json:"no_profit_desc"`
			MarketCapital            float64     `json:"market_capital"`
			Dividend                 float64     `json:"dividend"`
			DividendYield            float64     `json:"dividend_yield"`
			Currency                 string      `json:"currency"`
			Navps                    float64     `json:"navps"`
			Profit                   float64     `json:"profit"`
			Timestamp                float64     `json:"timestamp"`
			PeLyr                    float64     `json:"pe_lyr"`
			Amount                   float64     `json:"amount"`
			PledgeRatio              float64     `json:"pledge_ratio"`
			TradedAmountExt          interface{} `json:"traded_amount_ext"`
			IsRegistration           interface{} `json:"is_registration"`
			Pb                       float64     `json:"pb"`
			LimitUp                  float64     `json:"limit_up"`
			PeTtm                    float64     `json:"pe_ttm"`
			Time                     float64     `json:"time"`
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
