package job

import (
	"fmt"
	"fund/data"
	"fund/db"
	"fund/log"
	"fund/notifier"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// find schedulers in https://www.easycron.com/
func DailyReport(c *gin.Context) {
	username := os.Getenv("username")
	password := os.Getenv("password")
	u, p, ok := c.Request.BasicAuth()
	if !ok {
		fmt.Println("Error parsing basic auth")
		c.String(http.StatusUnauthorized, "Error parsing basic auth")
		return
	}
	if u != username {
		log.Info("Username provided is correct: %s\n", u)
		c.String(http.StatusUnauthorized, "Error parsing basic auth")
		return
	}
	if p != password {
		log.Info("Password provided is correct: %s\n", u)
		c.String(http.StatusUnauthorized, "Error parsing basic auth")
		return
	}
	log.Info("Username: %s\n", u)
	log.Info("Password: %s\n", p)
	go SendReport()
	c.String(http.StatusOK, "Sending email...")
}

func SendReport() {
	// save last data summary

	// update the database
	// CN sh_zs
	// HK hk
	// US us
	cnCount := data.GetStockCount("")
	hkCount := data.GetStockCount("hk")
	usCount := data.GetStockCount("us")

	// find largest cmc_id in old collection
	cryptoLargest := data.GetLagestCryptoID()

	// compare data summary and send main changes

	chs := make(chan string, 4)

	go func() {
		data.UpdateCoinList()
		chs <- "updated coin list"
	}()

	go func() {
		data.UpdateStockList(data.Market{Country: "CN", Board: "sh_sz"})
		chs <- "updated CN stock list"
	}()

	go func() {
		data.UpdateStockList(data.Market{Country: "HK", Board: "hk"})
		chs <- "updated HK stock list"
	}()
	go func() {
		data.UpdateStockList(data.Market{Country: "US", Board: "us"})
		chs <- "updated US stock list"
	}()

	for i := 0; i < 4; i++ {
		log.Info(<-chs)
	}

	e := &notifier.Email{
		To:      "dukeyunz@hotmail.com",
		Subject: "bot daily report",
	}

	templateData := struct {
		NewCryptoCount  int64
		NewCNStockCount int64
		NewHKStockCount int64
		NewUSStockCount int64
		NewCryptos      []db.CoinData
	}{
		NewCryptoCount:  data.GetNewCryptosCount(cryptoLargest),
		NewCNStockCount: data.GetStockCount("") - cnCount,
		NewHKStockCount: data.GetStockCount("hk") - hkCount,
		NewUSStockCount: data.GetStockCount("us") - usCount,
		NewCryptos:      data.GetNewCryptos(cryptoLargest),
	}

	if err := e.ParseTemplate("job/template.html", templateData); err == nil {
		e.Send()
	} else {
		log.Error("Parse template failed: %s\n", err.Error())
	}
}

// func reportChart() {
// 	graph := chart.Chart{
// 		Series: []chart.Series{
// 			chart.ContinuousSeries{
// 				XValues: []float64{1.0, 2.0, 3.0, 4.0},
// 				YValues: []float64{1.0, 2.0, 3.0, 4.0},
// 			},
// 		},
// 	}

// 	f, _ := os.Create("job/output.png")
// 	defer f.Close()
// 	graph.Render(chart.PNG, f)
// }
