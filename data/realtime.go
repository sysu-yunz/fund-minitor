package data

import (
	"encoding/json"
	"fmt"
	"fund/log"
	"io/ioutil"
	"net/http"
	"regexp"
)

func GetRealTime(fundCode string) {
	realTime := getRealTime(fundCode)
	fmt.Println(realTime.Name, realTime.Gszzl)
}

func getRealTime(fundCode string) RealTimeRaw {
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

	r, _ := regexp.Compile(`\((.*?)\)`)
	//fmt.Println(r.MatchString(string(body)))
	s := r.FindStringSubmatch(string(body))

	realTimeData := RealTimeRaw{}
	err = json.Unmarshal([]byte(s[1]), &realTimeData)
	if err != nil {
		log.Error("Unmarshal failed! ", err)
	}

	return realTimeData
}

type RealTimeRaw struct {
	Fundcode string `json:"fundcode"`
	Name     string `json:"name"`
	Jzrq     string `json:"jzrq"`
	Dwjz     string `json:"dwjz"`
	Gsz      string `json:"gsz"`
	Gszzl    string `json:"gszzl"`
	Gztime   string `json:"gztime"`
}
