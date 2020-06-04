package db

import (
	"context"
	"fund/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDB(pwd string) *mongo.Client {
	uri := "mongodb+srv://chengqian"+":"+pwd+"@cluster0-01hyt.azure.mongodb.net/fund?retryWrites=true&w=majority"

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

	return client
}
