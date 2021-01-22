package data

import (
	"encoding/json"
	"fmt"
	"fund/log"
	"io/ioutil"
	"net/http"
)

func GetChina10YearBondYield() BondDataRaw {
	url := "http://www.chinamoney.com.cn/ags/ms/cm-u-bk-currency/SddsIntrRateGovYldHis?lang=CN&pageNum=1&pageSize=15"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var bond BondDataRaw
	err = json.Unmarshal(body, &bond)
	if err != nil {
		log.Error("Bond Unmarshal Error %+v", err)
		log.Debug(string(body))
	}

	return bond
}

type BondDataRaw struct {
	Records []Records `json:"records"`
}
type Records struct {
	OneRate    string `json:"oneRate"`
	DateString string `json:"dateString"`
	TenRate    string `json:"tenRate"`
}
