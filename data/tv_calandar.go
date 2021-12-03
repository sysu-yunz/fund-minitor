package data

import (
	"fmt"
	ll "fund/log"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetTVData(tvName string) string {
	url := "https://www.episodate.com/tv-show/" + tvName
	doc := reqTVHTML(url)

	return findLatestEpisode(doc)
}

// copied from https://github.com/sysu-yunz/doubanAnalysis/blob/main/main.go
func reqTVHTML(url string) *goquery.Document {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ll.Error("[NewRequest]错误")
	}

	if req != nil {
		req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")
	} else {
		ll.Error("[req.Header] error")
	}

	ll.Info("正在请求网页: %s", url)
	res, err := client.Do(req)

	if err != nil {
		ll.Error("Get tv error", err)
	}

	if res == nil || res.Body == nil {
		ll.Error("[tv res] error")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ll.Error("Read tv error", err)
	}
	htmlContent := fmt.Sprintf("%s\n", body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func findLatestEpisode(doc *goquery.Document) string {
	res := ""
	doc.Find(".cd-timeline-content").Each(func(i int, s *goquery.Selection) {
		t := s.Find(".title").Text()
		d := s.Find(".episode-datetime-convert").Text()
		ll.Info("Title-----------%s", t)
		ll.Info("Date-----------%s", d)

		dd := parseTVDate(d)

		if dd.Before(time.Now().Add(-240 * time.Hour)) {
			ll.Info("[TVDate] Skip old episode %v", dd)
		} else if dd.After(time.Now()) && dd.Before(time.Now().Add(168*time.Hour)) {
			res = res + "\n" + t + "\n" + d + "  《====" + "\n"
		} else {
			res = res + "\n" + t + "\n" + d + "\n"
		}
	})

	return res
}

func parseTVDate(d string) time.Time {
	layout := "2 January 2006 15:04"
	// 15 June 2020 01:00
	date, err := time.Parse(layout, d)
	if err != nil {
		ll.Error("[TV] Parse date fail ", err)
	}

	return date
}
