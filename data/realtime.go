package data

import (
	"encoding/json"
	"fmt"
	"fund/global"
	"fund/log"
	r "fund/reply"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	tb "github.com/olekukonko/tablewriter"
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
		//reply = append(reply, []string{raw.Fundcode, raw.Gszzl, raw.Name})
		reply = append(reply, []string{raw.Gszzl, raw.Name})
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
	//table.SetHeader([]string{"CODE", "RATE", "NAME"})
	table.SetHeader([]string{"RATE", "NAME"})

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
	Jzrq     string `json:"jzrq"`
	Dwjz     string `json:"dwjz"`
	Gsz      string `json:"gsz"`
	Gszzl    string `json:"gszzl"`
	Gztime   string `json:"gztime"`
}
