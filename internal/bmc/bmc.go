// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"
	"os"

	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const databaseName = "bmc-app"

var (
	mongoAddr = "mongodb://localhost:27017/"

	log = logrus.New()
)

func init() {
	if value, ok := os.LookupEnv("MONGO_ADDR"); ok {
		mongoAddr = value
	}
}

type BmcApp struct {
	pbbmc.UnimplementedBmcServiceServer

	mongoClient *mongo.Client

	database *mongo.Database
}

func New(ctx context.Context) (bmcApp *BmcApp, err error) {
	bmcApp = &BmcApp{}

	clientOptions := options.Client().ApplyURI(mongoAddr)
	if bmcApp.mongoClient, err = mongo.Connect(ctx, clientOptions); err != nil {
		return
	}

	if err = bmcApp.mongoClient.Ping(ctx, nil); err != nil {
		return
	}

	bmcApp.database = bmcApp.mongoClient.Database(databaseName)

	return
}

func (b *BmcApp) FindOneAndReplace(ctx context.Context, colName string, filter bson.D, body any) (err error) {
	collection := b.database.Collection(colName)

	singleResult := collection.FindOneAndReplace(ctx, filter, body, options.FindOneAndReplace().SetUpsert(true).SetReturnDocument(options.After))
	if err = singleResult.Err(); err != nil {
		return
	}
	return singleResult.Decode(body)
}
