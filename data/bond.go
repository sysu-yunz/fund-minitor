package data

import (
	"encoding/json"
	"fmt"
	"fund/log"
	tb "github.com/olekukonko/tablewriter"
	"io/ioutil"
	"net/http"
	"strings"
)

func BondReply() string {
	var reply [][]string
	bond := getChina10YearBondYield().Records

	reply = append(reply, []string{"10 Rate", bond[0].TenRate})
	reply = append(reply, []string{"10 RateLast", bond[1].TenRate})
	reply = append(reply, []string{"1 Rate", bond[0].OneRate})
	reply = append(reply, []string{"1 RateLast", bond[1].OneRate})
	reply = append(reply, []string{"Date", bond[0].DateString})

	tableString := &strings.Builder{}
	table := tb.NewWriter(tableString)
	table.SetAlignment(tb.ALIGN_LEFT)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("+")
	table.SetHeader([]string{"ITEM", "VALUE"})

	for _, v := range reply {
		table.Append(v)
	}

	table.Render()

	return "<pre>"+tableString.String()+"</pre>"
}

func getChina10YearBondYield() BondDataRaw {
	url := "http://www.chinamoney.com.cn/ags/ms/cm-u-bk-currency/SddsIntrRateGovYldHis?lang=CN&pageNum=1&pageSize=15"
	method := "GET"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	var bond BondDataRaw
	err = json.Unmarshal(body, &bond)
	if err != nil{
		log.Error("Bond Unmarshal Error %+v", err)
		log.Debug(string(body))
	}

	return bond
}

type BondDataRaw struct {
	Records []Records `json:"records"`
}
type Records struct {
	OneRate    string      `json:"oneRate"`
	DateString string      `json:"dateString"`
	TenRate    string      `json:"tenRate"`
}