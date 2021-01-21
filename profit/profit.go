package profit

import (
	"fmt"
	"fund/data"
	"fund/global"
	"fund/log"
	"github.com/spf13/cast"
)

func PrevMonthProfitRate(fc string) float64 {
	return IntervalRate(fc, 20)
}

func LastProfitRate(fc string) float64 {
	return IntervalRate(fc, 1)
}

func UpdateLastRate() {
	hs := global.MgoDB.GetHolding(481088602)
	for _, sh := range hs.Shares {
		fmt.Println(LastProfitRate(sh.Code), sh.Code, sh.Shares)
	}
}

func IntervalRate(fc string, days int) float64 {
	historyData := data.GetFundHistoryData(fc, days)

	if days == 1 {
		return cast.ToFloat64(historyData[0].JZZZL)
	} else {
		if len(historyData) > 0 {
			nowStr := historyData[0].DWJZ
			prevStr := historyData[len(historyData)-1].DWJZ

			return cast.ToFloat64(nowStr) - cast.ToFloat64(prevStr)
		} else {
			log.Error("Get history data failed!")
		}
	}

	return 0
}
