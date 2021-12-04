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
	uri := "mongodb+srv://chengqian" + ":" + pwd + "@cluster0-01hyt.azure.mongodb.net/?retryWrites=true&w=majority"
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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

// delete all data in collection cryptos
func (c *MgoC) DeleteCryptoList() {
	col := c.Database("fund").Collection("crypto")
	_, err := col.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		log.Error("Deleting cryptos %+v", err)
	}
}

// insert a list of cryptos
func (c *MgoC) InsertCryptoList(cryptos []CoinData) {
	col := c.Database("fund").Collection("crypto")
	// avoid type error in col.InsertMany
	iCryptoList := make([]interface{}, len(cryptos))
	for i, v := range cryptos {
		iCryptoList[i] = v
	}

	_, err := col.InsertMany(context.TODO(), iCryptoList)
	if err != nil {
		log.Error("Inserting cryptos %+v ", err)
	}
}

func stockCol(market string) string {
	if market != "" {
		return "stock_" + market
	}

	return "stock"
}

func (c *MgoC) DeleteStockList(market string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	col_name := stockCol(market)
	col := c.Database("fund").Collection(col_name)
	_, err := col.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Error("Deleting stock %+v ", err)
	}

	log.Debug("Deleted stock %+v ", col_name)
}

func (c *MgoC) InsertStockList(StockList StockList, market string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	col_name := stockCol(market)
	col := c.Database("fund").Collection(col_name)

	// avoid type error in col.InsertMany
	iStockList := make([]interface{}, len(StockList.Data.List))
	for i, v := range StockList.Data.List {
		iStockList[i] = v
	}

	_, err := col.InsertMany(ctx, iStockList)
	if err != nil {
		log.Error("Inserting stock %+v ", err)
	}
}

func (c *MgoC) GetCryptoCount() int64 {
	// get doc count of crypto collection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	col := c.Database("fund").Collection("crypto")
	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Error("Getting crypto count %+v", err)
	}
	return count
}

func (c *MgoC) GetLagestCryptoID() int {
	// get lagest id of crypto collection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	col := c.Database("fund").Collection("crypto")
	var result CoinData
	err := col.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"id": -1})).Decode(&result)
	if err != nil {
		log.Error("Getting lagest crypto id %+v", err)
	}

	return result.ID
}

func (c *MgoC) GetNewCryptos(oldID int) []CoinData {
	// find cryptos which id larger than oldID
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	col := c.Database("fund").Collection("crypto")
	cur, err := col.Find(ctx, bson.M{"id": bson.M{"$gt": oldID}})
	if err != nil {
		log.Error("Getting new cryptos %+v", err)
	}

	var results []CoinData
	for cur.Next(ctx) {
		var result CoinData
		err := cur.Decode(&result)
		if err != nil {
			log.Error("Decode new crypto %+v", err)
		}
		results = append(results, result)
	}

	return results
}

func (c *MgoC) GetNewCryptosCount(oldID int) int64 {
	// get count of cryptos which id larger than oldID
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	col := c.Database("fund").Collection("crypto")
	count, err := col.CountDocuments(ctx, bson.M{"id": bson.M{"$gt": oldID}})
	if err != nil {
		log.Error("Getting new cryptos count %+v", err)
	}

	return count
}

func (c *MgoC) GetStockCount(market string) int64 {
	// get doc count of stock collection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	var col_name string
	if market == "us" || market == "hk" {
		col_name = "stock_" + market
	} else if market == "" {
		col_name = "stock"
	}
	col := c.Database("fund").Collection(col_name)
	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Error("Getting stock count %+v", err)
	}
	return count
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

func (c *MgoC) InsertMoviesBasic(ms []Movie) {
	col := c.Database("douban").Collection("movie")
	var msi []interface{}
	for _, t := range ms {
		msi = append(msi, t)
	}
	res, err := col.InsertMany(context.TODO(), msi)

	if err != nil {
		log.Error("Inserting movie %+v ", err)
	}

	log.Debug("Inserted movie %+v ", res)
}

func (c *MgoC) GetMovies() *mongo.Cursor {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := c.Database("douban").Collection("movie")
	cur, err := col.Find(ctx, bson.M{"date": bson.M{
		"$lt": "2013-10-31",
	}})
	if err != nil {
		log.Error("Finding watches %+v", err)
	}

	return cur

}

func (c *MgoC) GetUnmarkedMovies() *mongo.Cursor {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := c.Database("douban").Collection("movie")

	// find all movies which runtime is ""
	cur, err := col.Find(ctx, bson.M{"runtime": ""})
	if err != nil {
		log.Error("Finding unmarked movies %+v", err)
	}

	return cur
}

func (c *MgoC) GetAllMovies() []Movie {

	var ms []Movie

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	col := c.Database("douban").Collection("movie")
	cur, err := col.Find(ctx, bson.M{})
	if err != nil {
		log.Error("Finding watches %+v", err)
	}

	for cur.Next(ctx) {
		var result Movie
		err := cur.Decode(&result)
		if err != nil {
			log.Error("Decode watch %+v", err)
		}
		ms = append(ms, result)
	}

	return ms
}

func (c *MgoC) UpdateMovieRT(m Movie) {
	filter := bson.M{"subject": m.Subject}
	update := bson.M{"$set": bson.M{"ep": m.Ep, "runtime": m.RunTime}}

	col := c.Database("douban").Collection("movie")

	updateResult, err := col.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal("%s", err)
	}

	log.Info("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}
