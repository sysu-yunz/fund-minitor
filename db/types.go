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

type StockInfoDB struct {
}

type Holdings struct {
	ChatID int64   `json:"chatID" bson:"chatID"`
	Shares []Share `json:"Shares"`
	Bitcoin float64 `json:"bitcoin"`
}

type Share struct {
	SecurityType int64
	Code         string  `json:"Code"`
	Shares       float64 `json:"Shares"`
	Cost         float64 `json:"Cost"`
}
