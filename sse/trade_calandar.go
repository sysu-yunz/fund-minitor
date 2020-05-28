package sse

/*
	Shanghai Stock Exchange
 */

type publicHoliday struct  {
	start string
	duration int
}

// 到期日前的连续两个交易日
func getRawData() []publicHoliday {
	holidays := []publicHoliday{
		{"2020.1.1", 1},
		{"2020.1.24", 7},
		{"2020.4.4", 3},
		{"2020.5.1", 7},
		{"2020.6.25", 3},
		{"2020.10.1", 8},
	}

	return holidays
}
