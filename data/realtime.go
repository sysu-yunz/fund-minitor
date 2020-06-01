package data

import (
	"encoding/json"
	"fmt"
	"fund/config"
	"fund/log"
	tb "github.com/olekukonko/tablewriter"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func RealTimeFundReply() string {
	watchFunds := config.GetWatches()

	var reply [][]string

	ch := make(chan realTimeRaw, len(watchFunds))

	for _, f := range watchFunds {
		go getRealTime(f, ch)
	}

	for range watchFunds {
		raw := <-ch
		reply = append(reply, []string{raw.Fundcode, raw.Gszzl, raw.Name})
	}

	tableString := &strings.Builder{}
	table := tb.NewWriter(tableString)
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("+")
	table.SetHeader([]string{"CODE", "RATE", "NAME"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	return "<pre>"+tableString.String()+"</pre>"
	//return "```"+tableString.String()+"```"
}

func getRealTime(fundCode string, ch chan realTimeRaw) {
	url := fmt.Sprintf("http://fundgz.1234567.com.cn/js/%v.js", fundCode)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	//fmt.Println(string(body))

	r, _ := regexp.Compile(`\((.*)\)`)
	//fmt.Println(r.MatchString(string(body)))
	s := r.FindStringSubmatch(string(body))

	realTimeData := realTimeRaw{}
	err = json.Unmarshal([]byte(s[1]), &realTimeData)
	if err != nil {
		log.Error("%+v Unmarshal failed! %+v %+v", fundCode, err, s[1])
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
