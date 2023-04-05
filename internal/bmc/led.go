// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"
	"fmt"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbled"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// device - remove if re-discovery and uuid updated
func (b *BmcApp) SwitchLed(ctx context.Context, request *pbled.Request) (device *pbredfish.Device, err error) {
	const colName = "devices"

	log.Infof("exec switch led")

	id, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return
	}

	var (
		match = bson.D{
			{"$match",
				bson.D{
					{"_id", id},
				},
			},
		}

		ledProject = bson.D{
			{"$project", bson.D{
				{"_id", bson.D{{"$first", "$chassis._id"}}},
				{"url", 1},
				{"chassis", bson.D{{"$first", "$chassis.chassis"}}},
			}},
		}
	)
	cur, err := b.database.Collection(colName).Aggregate(ctx, mongo.Pipeline{match, lookupService, lookupChasseez, ledProject})
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	var s = struct {
		// chassis id
		Id               primitive.ObjectID `bson:"_id,omitempty"`
		Url              string             `bson:"url,omitempty"`
		*redfish.Chassis `bson:"chassis,omitempty"`
	}{}

	if !cur.Next(ctx) {
		err = fmt.Errorf("device not found")
		return
	}

	if err = cur.Decode(&s); err != nil {
		return
	}

	config := gofish.ClientConfig{
		Endpoint: s.Url,
		Username: request.GetUsername(),
		Password: request.GetPassword(),
		Insecure: true,
	}

	redfishClient, err := gofish.ConnectContext(ctx, config)
	if err != nil {
		return
	}
	defer redfishClient.Logout()

	chassis, err := redfish.GetChassis(redfishClient, s.Chassis.ODataID)
	if err != nil {
		return
	}
	chassis.IndicatorLED = common.IndicatorLED(request.GetState())

	if err = chassis.Update(); err != nil {
		return
	}

	redfishChassis := db.RedfishChassis{
		Chassis: chassis,
	}

	update := bson.D{
		{"$set", redfishChassis},
	}

	_, err = b.database.Collection("chasseez").UpdateByID(ctx, s.Id, update)
	if err != nil {
		return
	}

	cur, err = b.database.Collection(colName).Aggregate(ctx, mongo.Pipeline{lookupService, lookupSystem, lookupManager, lookupChasseez, project})
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	if !cur.Next(ctx) {
		err = fmt.Errorf("device not found")
		return
	}

	var result db.RedfishDevice
	if err = cur.Decode(&result); err != nil {
		return
	}

	device = result.ToProto()
	return
}
