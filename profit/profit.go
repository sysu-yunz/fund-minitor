package profit

import "fund/data"

func MonthProfitRate(fc string) float64 {
	 return data.GetMonthProfitRate(fc)
}
