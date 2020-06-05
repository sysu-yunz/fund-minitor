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