package db

type WatchDB struct {
	ChatID int64  `json:"chatID" bson:"chatID"`
	Watch  string `json:"watch" bson:"watch"`
}

type FundInfoDB struct {
	FundCode   string `json:"fundCode"`
	FundShort  string `json:"fundShort"`
	FundName   string `json:"fundName"`
	FundType   string `json:"fundType"`
	FundPinYin string `json:"fundPinYin"`
}

type Holdings struct {
	ChatID  int64   `json:"chatID" bson:"chatID"`
	Shares  []Share `json:"Shares"`
	Bitcoin float64 `json:"bitcoin"`

	Cryptos []CryptoHolding `json:"cryptos"`
	Stocks  []StockHolding  `json:"stocks"`
}

type Share struct {
	SecurityType int64
	Code         string  `json:"Code"`
	Shares       float64 `json:"Shares"`
	Cost         float64 `json:"Cost"`
}

type CryptoHolding struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Cost   float64 `json:"cost"`
	Amount float64 `json:"shares"`
}

type StockHolding struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Cost   float64 `json:"cost"`
	Amount float64 `json:"shares"`
}
type StockInfo struct {
	Symbol             string      `json:"symbol"`
	NetProfitCagr      float64     `json:"net_profit_cagr"`
	NorthNetInflow     interface{} `json:"north_net_inflow"`
	Ps                 float64     `json:"ps"`
	Type               int         `json:"type"`
	Percent            float64     `json:"percent"`
	HasFollow          bool        `json:"has_follow"`
	TickSize           float64     `json:"tick_size"`
	PbTtm              float64     `json:"pb_ttm"`
	FloatShares        int         `json:"float_shares"`
	Current            float64     `json:"current"`
	Amplitude          float64     `json:"amplitude"`
	Pcf                float64     `json:"pcf"`
	CurrentYearPercent float64     `json:"current_year_percent"`
	FloatMarketCapital interface{} `json:"float_market_capital"`
	NorthNetInflowTime interface{} `json:"north_net_inflow_time"`
	MarketCapital      int64       `json:"market_capital"`
	DividendYield      interface{} `json:"dividend_yield"`
	LotSize            int         `json:"lot_size"`
	RoeTtm             float64     `json:"roe_ttm"`
	TotalPercent       float64     `json:"total_percent"`
	Percent5M          float64     `json:"percent5m"`
	IncomeCagr         float64     `json:"income_cagr"`
	Amount             float64     `json:"amount"`
	Chg                float64     `json:"chg"`
	IssueDateTs        int64       `json:"issue_date_ts"`
	Eps                float64     `json:"eps"`
	MainNetInflows     interface{} `json:"main_net_inflows"`
	Volume             int         `json:"volume"`
	VolumeRatio        interface{} `json:"volume_ratio"`
	Pb                 float64     `json:"pb"`
	Followers          int         `json:"followers"`
	TurnoverRate       float64     `json:"turnover_rate"`
	FirstPercent       float64     `json:"first_percent"`
	Name               string      `json:"name"`
	PeTtm              float64     `json:"pe_ttm"`
	TotalShares        int         `json:"total_shares"`
	LimitupDays        int         `json:"limitup_days"`
}

type StockList struct {
	Data struct {
		Count int         `json:"count"`
		List  []StockInfo `json:"list"`
	} `json:"data"`
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}
