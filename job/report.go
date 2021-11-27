package job

import (
	"fmt"
	"fund/data"
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
	go UpdateReport()
	c.String(http.StatusOK, "ok")
}

func UpdateReport() {
	// save last data summary

	// update the database
	// CN sh_zs
	// HK hk
	// US us
	cryptoCount := data.GetCryptoCount()
	cnCount := data.GetStockCount("")
	hkCount := data.GetStockCount("hk")
	usCount := data.GetStockCount("us")

	// data.UpdateStockList(data.Market{Country: "CN", Board: "sh_zs"})
	// data.UpdateStockList(data.Market{Country: "HK", Board: "hk"})
	// data.UpdateStockList(data.Market{Country: "US", Board: "us"})

	// data.UpdateCoinList()

	// compare data summary and send main changes

	e := &notifier.Email{
		To:      "dukeyunz@hotmail.com",
		Subject: "bot daily report",
	}

	templateData := struct {
		NewCryptoCount  int64
		NewCNStockCount int64
		NewHKStockCount int64
		NewUSStockCount int64
	}{
		NewCryptoCount:  data.GetCryptoCount() - cryptoCount,
		NewCNStockCount: data.GetStockCount("") - cnCount,
		NewHKStockCount: data.GetStockCount("hk") - hkCount,
		NewUSStockCount: data.GetStockCount("us") - usCount,
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
