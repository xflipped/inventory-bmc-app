// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

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
	log.Infof("exec switch led")

	id, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return
	}

	var (
		match = bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "_id", Value: id},
				},
			},
		}

		ledProject = bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: bson.D{{Key: "$first", Value: "$chassis._id"}}},
					{Key: "url", Value: 1},
					{Key: "chassis", Value: bson.D{{Key: "$first", Value: "$chassis.chassis"}}},
				}},
		}
	)
	cur, err := b.database.Collection(devicesColName).Aggregate(ctx, mongo.Pipeline{match, lookupService, lookupChasseez, ledProject})
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
		err = errDeviceNotFound
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
		{Key: "$set", Value: redfishChassis},
	}

	_, err = b.database.Collection(chasseezColName).UpdateByID(ctx, s.Id, update)
	if err != nil {
		return
	}

	deviceCur, err := b.database.Collection(devicesColName).Aggregate(ctx, mongo.Pipeline{lookupService, lookupSystem, lookupManager, lookupChasseez, project})
	if err != nil {
		return
	}
	defer deviceCur.Close(ctx)

	if !deviceCur.Next(ctx) {
		err = errDeviceNotFound
		return
	}

	var result db.RedfishDevice
	if err = deviceCur.Decode(&result); err != nil {
		return
	}

	device = result.ToProto()
	return
}
