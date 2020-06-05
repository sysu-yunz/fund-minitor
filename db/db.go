package db

import (
	"context"
	"fund/log"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MgoC struct {
	*mongo.Client
}

func NewDB(pwd string) *MgoC {
	uri := "mongodb+srv://chengqian"+":"+pwd+"@cluster0-01hyt.azure.mongodb.net/fund?retryWrites=true&w=majority"
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
	err := col.FindOne(context.TODO(), bson.M{"fundCode":w}).Decode(&fund)
	if err != nil {
		log.Error("InsertWatch %+v", err)
		return FundInfoDB{}, false
	}

	return fund, true
}

func (c *MgoC) FundWatched(cid int64, w string)  bool {
	col := c.Database("fund").Collection("watch")
	count, err := col.CountDocuments(context.TODO(), bson.M{"chatID":cid, "watch":w})
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
	cur, err := col.Find(ctx, bson.M{"chatID":chatID})
	if err != nil {
		log.Error("Finding watches %+v", err)
	}

	var results []WatchDB
	for cur.Next(ctx) {
		var result WatchDB
		err := cur.Decode(&result)
		if err != nil { log.Error("Decode watch %+v", err) }
		results = append(results, result)
	}

	return results
}
