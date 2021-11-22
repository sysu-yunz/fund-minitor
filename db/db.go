package db

import (
	"context"
	"fund/log"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MgoC struct {
	*mongo.Client
}

func NewDB(pwd string) *MgoC {
	uri := "mongodb+srv://chengqian" + ":" + pwd + "@cluster0-01hyt.azure.mongodb.net/fund?retryWrites=true&w=majority"
	//uri := "mongodb://127.0.0.1:27017"
	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Error("New Client %+v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Error("Connect %+v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error("Ping %+v", err)
	}

	log.Debug("Connected to MongoDB!")

	return &MgoC{client}
}

func (c *MgoC) ValidFundCode(w string) (FundInfoDB, bool) {
	col := c.Database("fund").Collection("basic")
	fund := FundInfoDB{}
	err := col.FindOne(context.TODO(), bson.M{"fundCode": w}).Decode(&fund)
	if err != nil {
		log.Error("InsertWatch %+v", err)
		return FundInfoDB{}, false
	}

	return fund, true
}

func (c *MgoC) FundWatched(cid int64, w string) bool {
	col := c.Database("fund").Collection("watch")
	count, err := col.CountDocuments(context.TODO(), bson.M{"chatID": cid, "watch": w})
	if err != nil {
		log.Error("Fund watch count %+v", err)
	}

	return count > 0
}

func (c *MgoC) InsertWatch(chatID int64, w string) {
	col := c.Database("fund").Collection("watch")
	res, err := col.InsertOne(context.TODO(), WatchDB{
		ChatID: chatID,
		Watch:  w,
	})

	if err != nil {
		log.Error("Inserting watch %+v ", err)
	}

	log.Debug("Inserted watch %+v ", res)
}

func (c *MgoC) DeleteWatch(chatID int64, w string) {
	col := c.Database("fund").Collection("watch")
	res, err := col.DeleteOne(context.TODO(), WatchDB{
		ChatID: chatID,
		Watch:  w,
	})

	if err != nil {
		log.Error("Deleting watch %+v ", err)
	}

	log.Debug("Deleted watch %+v ", res)
}

func (c *MgoC) GetWatchList(chatID int64) []WatchDB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := c.Database("fund").Collection("watch")
	cur, err := col.Find(ctx, bson.M{"chatID": chatID})
	if err != nil {
		log.Error("Finding watches %+v", err)
	}

	var results []WatchDB
	for cur.Next(ctx) {
		var result WatchDB
		err := cur.Decode(&result)
		if err != nil {
			log.Error("Decode watch %+v", err)
		}
		results = append(results, result)
	}

	return results
}

func (c *MgoC) GetHolding(chatID int64) Holdings {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := c.Database("fund").Collection("hold")
	res := col.FindOne(ctx, bson.M{"chatID": chatID})

	var result Holdings
	err := res.Decode(&result)
	if err != nil {
		log.Error("Decode watch %+v", err)
	}

	return result
}

func (c *MgoC) InsertHold(hold Holdings) {
	col := c.Database("fund").Collection("hold")
	res, err := col.InsertOne(context.TODO(), hold)

	if err != nil {
		log.Error("Inserting hold %+v ", err)
	}

	log.Debug("Inserted hold %+v ", res)
}

func (c *MgoC) DeleteStockList() {
	col := c.Database("fund").Collection("stock_us")
	deleteRes, err := col.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		log.Error("Deleting stock %+v ", err)
	}

	log.Debug("Deleted stock %+v ", deleteRes)
}

func (c *MgoC) UpdateStockList(StockList StockList) {

	col := c.Database("fund").Collection("stock_us")
	// avoid type error in col.InsertMany
	iStockList := make([]interface{}, len(StockList.Data.List))
	for i, v := range StockList.Data.List {
		iStockList[i] = v
	}

	insertRes, err := col.InsertMany(context.TODO(), iStockList)
	if err != nil {
		log.Error("Inserting stock %+v ", err)
	}

	log.Debug("Inserted stock %+v ", insertRes)
}

// fuzzy search stock in stock list
func (c *MgoC) SearchStock(arg string, fuzzy bool) string {
	if s := c.findStock("", arg, fuzzy); s != "" {
		return s
	} else if s := c.findStock("hk", arg, fuzzy); s != "" {
		return s
	} else if s := c.findStock("us", arg, fuzzy); s != "" {
		return s
	} else {
		return ""
	}
}

func (c *MgoC) findStock(market string, arg string, fuzzy bool) string {
	stock := StockInfo{}
	col := c.Database("fund").Collection("stock")

	if market != "" {
		col = c.Database("fund").Collection("stock_" + market)
	}

	if fuzzy {
		res := col.FindOne(context.TODO(), bson.M{"$or": []bson.M{
			{"name": bson.M{"$regex": arg}},
			{"symbol": bson.M{"$regex": arg}},
		}}).Decode(&stock)
		if res != nil {
			log.Error("Stock not in %s market, finding it in %s %+v", market, market, res)
			return ""
		}

		return stock.Symbol
	}

	// find doc name or symbol match
	res := col.FindOne(context.TODO(), bson.M{"$or": []bson.M{
		{"name": arg},
		{"symbol": arg},
	}}).Decode(&stock)

	if res != nil {
		log.Error("Stock not in %s market, finding it in %s %+v", market, market, res)
		return ""
	}

	return stock.Symbol
}
