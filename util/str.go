package util

import (
	"fund/log"
	"github.com/yanyiwu/gojieba"
	"strings"
)

func ShortenFundName(fundName string) string {
	x := gojieba.NewJieba()
	defer x.Free()
	words := x.Cut(fundName, true)
	log.Debug(strings.Join(words, "-"))
	return strings.Join(words[:2], "")
}
