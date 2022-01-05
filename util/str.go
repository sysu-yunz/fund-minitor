package util

import (
	"fund/log"
	"strings"

	"github.com/yanyiwu/gojieba"
)

func ShortenFundName(fundName string) string {
	type name struct {
		Full  string
		Short string
	}

	specials := []name{
		{"天弘中证食品饮料指数C", "食品饮料"},
		{"华夏沪港通恒生ETF联接C", "恒指ETF"},
		{"汇添富中证生物科技指数A", "生物科技"},
		{"易方达裕祥回报债券", "高回报债"},
		{"中银金融地产混合", "金融地产"},
		{"交银创新成长混合", "交银成长"},
		{"广发中证全指家用电器指数C", "家电ETF"},
		{"工银新能源汽车混合C", "新能车"},
		{"广发中证基建工程指数C", "基建工程"},
		{"天弘中证电子ETF联接C", "电子ETF"},
		{"汇添富全球消费混合人民币C", "全球消费"},
		{"广发纳斯达克100指数C", "纳指100"},
		{"天弘中证500指数C", "中证500"},
		{"工银战略转型股票A", "D-工银战略转型股票A"},
		{"天弘中证银行ETF联接C", "D-天弘中证银行C"},
		{"天弘中证食品饮料ETF联接C", "M-天弘食品饮料C"},
		{"易方达上证50增强C", "M-易方达上证50增强C"},
		{"景顺长城鼎益混合(LOF)", "D-景顺长城鼎益混合"},
		{"嘉实农业产业股票", "M-嘉实农业产业股票"},
		{"易方达安心回馈混合", "D-易方达安心回馈混合"},
		{"易方达瑞程灵活配置混合C", "D-易方达瑞程混合C"},
		{"中融新经济混合C", "M-中融新经济混合C"},
	}

	for _, n := range specials {
		if fundName == n.Full {
			return n.Short
		}
	}
	x := gojieba.NewJieba()
	defer x.Free()

	// 耗时
	words := x.Cut(fundName, true)
	log.Debug(strings.Join(words, "-"))
	if len(words) >= 2 {
		return strings.Join(words[1:2], "")
	}

	return strings.Join(words, "")
}
