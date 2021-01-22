package bot

import (
	"fund/data"
	"fund/log"
	"github.com/spf13/cast"
)

func PrevMonthProfitRate(fc string) float64 {
	return IntervalProfitRate(fc, 20)
}

func LastProfitRate(fc string) float64 {
	return IntervalProfitRate(fc, 1)
}

func LastProfit(fc string) float64 {
	lastTwo := data.GetFundHistoryData(fc, 2)
	return cast.ToFloat64(lastTwo[0].DWJZ) - cast.ToFloat64(lastTwo[1].DWJZ)
}

func IntervalProfitRate(fc string, days int) float64 {
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
