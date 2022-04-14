package data

import (
	"encoding/json"
	"fmt"
	"fund/global"
	"fund/log"
	"io/ioutil"
	"os"
)

func BasicVideoInfo() {
	// 从google获取基本信息
	// 用youtube-dl获取视频时长
	// 分析数据

	// 测试获取视频时长
	// getVideoDuration("https://www.youtube.com/watch?v=QH2-TGUlwu4")
	// var testClient = youtube.Client{Debug: true}
	// v, err := testClient.GetVideo("https://www.youtube.com/watch?v\u003dEI0pHRI04VQ")
	// if err != nil {
	// 	log.Error("%+v", err)
	// }
	// log.Info("%+v", v.Duration)

	var v []interface{}
	fmt.Println(os.Getwd())
	jsonFile, err := os.Open("./yt.json")
	if err != nil {
		log.Error("%+v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err := json.Unmarshal(byteValue, &v); err != nil {
		// handle error
	}
	global.MgoDB.InsertYTV(v)
}
