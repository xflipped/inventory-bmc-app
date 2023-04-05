// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbpower"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (b *BmcApp) SwitchPower(ctx context.Context, request *pbpower.Request) (device *pbredfish.Device, err error) {
	log.Infof("exec switch power")

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

		powerProject = bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: bson.D{{Key: "$first", Value: "$system._id"}}},
					{Key: "url", Value: 1},
					{Key: "system", Value: bson.D{{Key: "$first", Value: "$system.computersystem"}}},
				}},
		}
	)
	cur, err := b.database.Collection(devicesColName).Aggregate(ctx, mongo.Pipeline{match, lookupService, lookupSystem, powerProject})
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	var s = struct {
		// chassis id
		Id                     primitive.ObjectID `bson:"_id,omitempty"`
		Url                    string             `bson:"url,omitempty"`
		redfish.ComputerSystem `bson:"system,omitempty"`
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

	system, err := redfish.GetComputerSystem(redfishClient, s.ComputerSystem.ODataID)
	if err != nil {
		return
	}

	if err = system.Reset(redfish.ResetType(request.GetType())); err != nil {
		return
	}

	system, err = redfish.GetComputerSystem(redfishClient, s.ComputerSystem.ODataID)
	if err != nil {
		return
	}

	redfishSystem := db.RedfishSystem{
		ComputerSystem: system,
	}

	update := bson.D{
		{Key: "$set", Value: redfishSystem},
	}

	_, err = b.database.Collection(systemsColName).UpdateByID(ctx, s.Id, update)
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
