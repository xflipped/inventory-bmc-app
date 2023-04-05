// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"
	"fmt"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbinventory"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"github.com/stmcginnis/gofish"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// device - remove if re-discovery and uuid updated
func (b *BmcApp) Inventory(ctx context.Context, request *pbinventory.Request) (device *pbredfish.Device, err error) {
	const colName = "devices"

	log.Infof("exec inventory: %s", request.GetId())

	id, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return
	}

	collection := b.database.Collection(colName)

	filter := bson.D{{Key: "_id", Value: id}}

	singleResult := collection.FindOne(ctx, filter)
	if err = singleResult.Err(); err != nil {
		return
	}

	var redfishDevice db.RedfishDevice
	if err = singleResult.Decode(&redfishDevice); err != nil {
		return
	}

	config := gofish.ClientConfig{
		Endpoint: redfishDevice.Url,
		Username: request.GetUsername(),
		Password: request.GetPassword(),
		Insecure: true,
	}

	redfishClient, err := gofish.ConnectContext(ctx, config)
	if err != nil {
		return
	}
	defer redfishClient.Logout()

	if err = b.inventoryService(ctx, redfishDevice, redfishClient.Service); err != nil {
		return
	}

	cur, err := b.database.Collection(colName).Aggregate(ctx, mongo.Pipeline{lookupService, lookupSystem, lookupManager, lookupChasseez, project})
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
