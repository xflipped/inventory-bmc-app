package main

import (
	"context"
	"fmt"

	"github.com/foliagecp/inventory-bmc-app/internal/db"

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
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	if err = mongoClient.Ping(ctx, nil); err != nil {
		return
	}

	database := mongoClient.Database("bmc-app")

	lookupService := bson.D{
		{"$lookup",
			bson.D{
				{"from", "services"},
				{"localField", "_id"},
				{"foreignField", "_device_id"},
				{"as", "service"},
			},
		},
	}

	lookupSystem := bson.D{
		{"$lookup",
			bson.D{
				{"from", "systems"},
				{"localField", "service._id"},
				{"foreignField", "_service_id"},
				{"as", "system"},
			},
		},
	}

	lookupManager := bson.D{
		{"$lookup",
			bson.D{
				{"from", "managers"},
				{"localField", "service._id"},
				{"foreignField", "_service_id"},
				{"as", "manager"},
			},
		},
	}

	lookupChasseez := bson.D{
		{"$lookup",
			bson.D{
				{"from", "chasseez"},
				{"localField", "service._id"},
				{"foreignField", "_service_id"},
				{"as", "chassis"},
			},
		},
	}

	project := bson.D{
		{"$project", bson.D{
			{"_id", 1},

			{"url", 1},
			{"uuid", bson.D{{"$first", "$service.service.uuid"}}},
			{"serial_number", bson.D{{"$first", "$system.computersystem.serialnumber"}}},
			{"name", bson.D{{"$first", "$service.service.product"}}},

			// FIXME
			{"mac_address", bson.D{{"$first", "$manager.manager.mac"}}},

			{"model", bson.D{{"$first", "$manager.manager.model"}}},
			{"vendor", bson.D{{"$first", "$service.service.vendor"}}},
			{"power_state", bson.D{{"$first", "$manager.manager.powerstate"}}},
			{"health_status", bson.D{{"$first", "$system.computersystem.status.health"}}},
			{"indicator_led", bson.D{{"$first", "$system.computersystem.indicatorled"}}},
			// FIXME
			{"min_temp", bson.D{{"$first", "$manager.manager.temp"}}},
			// FIXME
			{"max_temp", bson.D{{"$first", "$manager.manager.temp"}}},
		}},
	}

	showInfoCursor, err := database.Collection("devices").Aggregate(ctx, mongo.Pipeline{lookupService, lookupSystem, lookupManager, lookupChasseez, project})
	if err != nil {
		return
	}

	var devices []db.RedfishDevice
	if err = showInfoCursor.All(ctx, &devices); err != nil {
		return
	}

	for _, dev := range devices {
		fmt.Println(dev)
	}

	return
}
