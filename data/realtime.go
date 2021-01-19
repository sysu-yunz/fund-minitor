package data

import (
	"encoding/json"
	"fmt"
	"fund/global"
	"fund/log"
	r "fund/reply"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tb "github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func RealTimeFundReply(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	watchFunds := global.MgoDB.GetWatchList(chatID)

	var reply [][]string

	ch := make(chan realTimeRaw, len(watchFunds))

	for _, f := range watchFunds {
		go getRealTime(f.Watch, ch)
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

	r.TextReply(update, "<pre>"+tableString.String()+"</pre>")
	//return "```"+tableString.String()+"```"
}

func HoldReply(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	holdFunds := global.MgoDB.GetHolding(chatID)

	var hr []holdReply
	ch := make(chan realTimeRaw, len(holdFunds.Shares))

	for _, f := range holdFunds.Shares {
		go getRealTime(f.Code, ch)
	}

	for range holdFunds.Shares {
		raw := <-ch
		estimateValue := cast.ToFloat64(raw.Gsz)
		estimateRate := cast.ToFloat64(raw.Gszzl)

		hr = append(hr, holdReply{
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
				h.Earn = h.Shares * (h.Price - h.Cost)
				h.Cap = h.Price*h.Shares
			}
		}

		reply = append(reply, []string{
			h.Code,
			fmt.Sprintf("%.1f", h.Earn),
			fmt.Sprintf("%.1f", h.Cap),
			fmt.Sprintf("%.4f", h.Cost),
			fmt.Sprintf("%.4f", h.Price),
			h.Name,
		})
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
	table.SetHeader([]string{"CODE", "EARN", "COST", "PRICE", "NAME"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	r.TextReply(update, "<pre>"+tableString.String()+"</pre>")
	//return "```"+tableString.String()+"```"
}

func getRealTime(fundCode string, ch chan realTimeRaw) {
	realTimeData := realTimeRaw{Fundcode: fundCode}

	url := fmt.Sprintf("http://fundgz.1234567.com.cn/js/%v.js", fundCode)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	if res != nil && res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)

		log.Debug("Real time data resp %+v", string(body))

		r, _ := regexp.Compile(`\((.*)\)`)
		//fmt.Println(r.MatchString(string(body)))
		s := r.FindStringSubmatch(string(body))

		err = json.Unmarshal([]byte(s[1]), &realTimeData)
		if err != nil {
			log.Error("Realtime data of %+v Unmarshal failed! %+v %+v", fundCode, err, s[1])
			fdb, _ := global.MgoDB.ValidFundCode(fundCode)
			realTimeData.Name = fdb.FundName
		}
	} else {
		log.Error("Http response %+v ", res)
		fdb, _ := global.MgoDB.ValidFundCode(fundCode)
		realTimeData.Name = fdb.FundName
	}

	ch <- realTimeData
}

type realTimeRaw struct {
	Fundcode string `json:"fundcode"`
	Name     string `json:"name"`
	Jzrq     string `json:"jzrq"`  // 净值日期
	Dwjz     string `json:"dwjz"`  // 单位净值
	Gsz      string `json:"gsz"`   // 估算值
	Gszzl    string `json:"gszzl"` // 估算增长率
	Gztime   string `json:"gztime"`
}

type holdReply struct {
	Code      string
	Rate      float64
	Name      string
	Price     float64
	Shares    float64
	Cost      float64
	Total     float64
	Earn      float64
	TodayEarn float64
	Cap float64
	CapPercent float64
}
