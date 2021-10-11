package data

import (
	"encoding/json"
	"fmt"
	"fund/global"
	"fund/log"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/spf13/cast"
)

func LastPrice(fc string) float64 {
	lastValue := GetFundHistoryData(fc, 1)

	return cast.ToFloat64(lastValue[0].DWJZ)
}

func GetRealTime(fundCode string, ch chan RealTimeRaw) {
	realTimeData := RealTimeRaw{Fundcode: fundCode}

	url := fmt.Sprintf("http://fundgz.1234567.com.cn/js/%v.js", fundCode)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)

	if err != nil {
		log.Error("Get tv error", err)
	}

	if res == nil || res.Body == nil {
		log.Error("[tv res] error")
	}

	defer res.Body.Close()

	if res != nil && res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error("Read realtime error", err)
		}

		log.Debug("Real time data resp %+v", string(body))

		if string(body) == "jsonpgz();" {
			fdb, _ := global.MgoDB.ValidFundCode(fundCode)
			realTimeData.Name = fdb.FundName
		} else {
			r, _ := regexp.Compile(`\((.*)\)`)
			s := r.FindStringSubmatch(string(body))
			err = json.Unmarshal([]byte(s[1]), &realTimeData)
			if err != nil {
				log.Error("Realtime data of %+v Unmarshal failed! %+v %+v", fundCode, err, s[1])
			}
			fdb, _ := global.MgoDB.ValidFundCode(fundCode)
			realTimeData.Name = fdb.FundName
		}
	} else {
		log.Error("Http response %+v ", res)
	}

	ch <- realTimeData
}

type RealTimeRaw struct {
	Fundcode string `json:"fundcode"`
	Name     string `json:"name"`
	Jzrq     string `json:"jzrq"`  // 净值日期
	Dwjz     string `json:"dwjz"`  // 单位净值
	Gsz      string `json:"gsz"`   // 估算值
	Gszzl    string `json:"gszzl"` // 估算增长率
	Gztime   string `json:"gztime"`
}

type HoldReply struct {
	Code       string
	Rate       float64
	Name       string
	Price      float64
	LastPrice  float64
	Shares     float64
	Cost       float64
	Total      float64
	Earn       float64
	TodayEarn  float64
	Cap        float64
	CapPercent float64
}
