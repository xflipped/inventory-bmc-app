package main

import (
	"context"
	"fmt"

	"github.com/stmcginnis/gofish/redfish"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := info(context.Background()); err != nil {
		fmt.Println(err)
		return
	}
}

func info(ctx context.Context) (err error) {
	const colName = "devices"
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	if err = mongoClient.Ping(ctx, nil); err != nil {
		return
	}

	database := mongoClient.Database("bmc-app")

	id, _ := primitive.ObjectIDFromHex("642d31745b01bb3669de8dda")

	var (
		match = bson.D{
			{"$match",
				bson.D{
					{"_id", id},
				},
			},
		}

		lookupService = bson.D{
			{"$lookup",
				bson.D{
					{"from", "services"},
					{"localField", "_id"},
					{"foreignField", "_device_id"},
					{"as", "service"},
				},
			},
		}

		lookupSystem = bson.D{
			{"$lookup",
				bson.D{
					{"from", "systems"},
					{"localField", "service._id"},
					{"foreignField", "_service_id"},
					{"as", "system"},
				},
			},
		}

		project = bson.D{
			{"$project", bson.D{
				{"_id", 1},

				{"url", 1},
				{"system", bson.D{{"$first", "$system.computersystem"}}},
			}},
		}
	)
	cur, err := database.Collection(colName).Aggregate(ctx, mongo.Pipeline{match, lookupService, lookupSystem, project})
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	var s = struct {
		Id                      primitive.ObjectID `bson:"_id,omitempty"`
		Url                     string             `bson:"url,omitempty"`
		*redfish.ComputerSystem `bson:"system,omitempty"`
	}{}

	if !cur.Next(ctx) {
		err = fmt.Errorf("device not found")
		return
	}

	if err = cur.Decode(&s); err != nil {
		return
	}

	fmt.Println(s.ComputerSystem)

	return
}
