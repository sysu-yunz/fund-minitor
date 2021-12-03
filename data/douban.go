package data

import (
	"context"
	"fmt"
	"fund/config"
	"fund/db"
	"fund/global"
	ll "fund/log"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func calTotalRT() {
	var total int
	ms := global.MgoDB.GetAllMovies()
	re := regexp.MustCompile("[0-9]+")
	for _, m := range ms {

		var ep, rt int
		var err error
		eps := re.FindAllString(m.Ep, -1)
		rts := re.FindAllString(m.RunTime, -1)

		if eps != nil {
			ep, err = strconv.Atoi(eps[0])
			if err != nil {
				ll.Error("%d", err)
			}
		} else {
			ep = 10
		}

		if rts != nil {
			rt, err = strconv.Atoi(rts[0])
			if err != nil {
				ll.Error("%d", err)
			}
		} else {
			rt = 60
		}

		total = total + ep*rt
	}

	fmt.Println(total)
}

func getMovies() []db.Movie {

	var ms []db.Movie

	url := `https://movie.douban.com/people/` + config.EnvVariable("DoubanID") + `/collect`
	doc := getHTML(url, `div[class="info"]`)
	start := findBasicSubjectInfo(doc)
	ms = append(ms, start...)

	totalPage := findTotalNum(doc)/15 + 1
	// totalPage := 1
	// 翻页-组装URL
	// https://movie.douban.com/people/dukeyunz/collect?start=15&sort=time&rating=all&filter=all&mode=grid
	// 因为第一页已经跑过一次了，直接从第二页开始

	for i := 1; i < totalPage; i++ {
		currentURL := fmt.Sprintf("%s?start=%d&sort=time&rating=all&filter=all&mode=grid", url, i*15)
		currentDoc := getHTML(currentURL, `div[class="info"]`)
		ms = append(ms, findBasicSubjectInfo(currentDoc)...)
	}

	return ms
}

func getRTMovies() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur := global.MgoDB.GetMovies()
	for cur.Next(ctx) {
		var result db.Movie
		err := cur.Decode(&result)
		if err != nil {
			ll.Error("Decode watch %+v", err)
		}

		// doc := getHTML(result.Link, "div#info")
		doc := reqHTML(result.Link)
		result.RunTime, result.Ep = findSubjectRunTime(doc)
		global.MgoDB.UpdateMovieRT(result)
	}
}

func getHTML(url string, wait interface{}) *goquery.Document {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("blink-settings", "imageEnable=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko)`),
	}

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	_ = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	timeOutCtx, cancel := context.WithTimeout(chromeCtx, 60*time.Second)
	defer cancel()

	var htmlContent string

	log.Printf("chrome visit page %s\n", url)
	err := chromedp.Run(timeOutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(wait),
		chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func reqHTML(url string) *goquery.Document {
	client := &http.Client{}

	var req *http.Request

	// req = config.ReqChrome(url)

	//if RandBool() {
	//	fmt.Println("------------------------------------------------------")
	//	req = reqChrome(url)
	//} else {
	//	fmt.Println("********************************************************")
	//	req = reqSafari(url)
	//}

	ll.Info("正在请求网页: %s", url)
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	htmlContent := fmt.Sprintf("%s\n", body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func findSubjectRunTime(doc *goquery.Document) (rt string, eps string) {

	//doc.Find("#info").Each(func(i int, s *goquery.Selection) {
	//	op, _ := s.Attr("property")
	//	con, _ := s.Attr("content")
	//	if op == "v:runtime" {
	//		fmt.Println(con)
	//	}
	//})
	var mark int
	doc.Find("div#info").Contents().Each(func(i int, s *goquery.Selection) {
		if s.Text() == "片长:" {
			rt = s.Next().Text()
			eps = "1"
			fmt.Printf("片长:%s\n", rt)
		}

		if s.Text() == "集数:" {
			goquery.NodeName(s.Next())
			fmt.Printf("集数: ")
			mark = 1
		}
		if s.Text() == "单集片长:" {
			fmt.Printf("单集片长:")
			mark = 2
		}

		if goquery.NodeName(s) == "#text" && (mark == 1 || mark == 2) {
			if mark == 1 {
				eps = s.Text()
			}

			if mark == 2 {
				rt = s.Text()
			}

			fmt.Println(s.Text())
			mark = 0
		}
	})

	return
}

func findTotalNum(doc *goquery.Document) int {
	s := doc.Find("h1").Text()
	re := regexp.MustCompile(`(?s)\((.*)\)`)
	m := re.FindAllStringSubmatch(s, -1)
	fmt.Printf(m[0][1])

	if num, err := strconv.Atoi(m[0][1]); err == nil {
		return num
	}

	return 0
}

func findBasicSubjectInfo(doc *goquery.Document) []db.Movie {
	// 获取内容
	// title 标题
	// <li class="title">
	//                        <a href="https://movie.douban.com/subject/26413293/" class="">
	//                            <em>大秦赋</em>
	//                             / 大秦帝国4：东出 / 大秦帝国之东出
	//                        </a>
	//                            <span class="playable">[可播放]</span>
	//                    </li>

	// link 链接
	// rate 评分
	// <span class="rating1-t"></span>

	// date 日期
	// <span class="date">2020-12-16</span>

	// comment 评论
	// <span class="comment">本来还说快进随便看看，弃剧了。以后国产剧一定放凉了再看，真nm坑。</span>

	// img 图片
	// <img alt="Warrior" src="https://img9.doubanio.com/view/photo/s_ratio_poster/public/p2619810129.webp" class="">

	var ms []db.Movie
	doc.Find(".item").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		img, _ := s.Find("img").Attr("src")
		titleSel := s.Find(".title")
		title := titleSel.Find("em").Text()
		link, _ := titleSel.Find("a").Attr("href")
		dateSel := s.Find(".date")
		rate, _ := dateSel.Prev().Attr("class")
		date := dateSel.Text()
		comment := s.Find(".comment").Text()
		la := strings.Split(link, "/")
		subject := la[len(la)-2]
		ms = append(ms, db.Movie{
			Subject: subject,
			Title:   title,
			Link:    link,
			Rate:    rate,
			Date:    date,
			Comment: comment,
			Img:     img,
		})
	})

	return ms
}
