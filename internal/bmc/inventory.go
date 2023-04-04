// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbinventory"
	"github.com/stmcginnis/gofish"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// device - remove if re-discovery and uuid updated
func (b *BmcApp) Inventory(ctx context.Context, request *pbinventory.Request) (response *pbinventory.Response, err error) {
	const dbName = "devices"

	log.Infof("exec inventory: %s", request.GetId())

	id, err := primitive.ObjectIDFromHex(request.GetId())
	if err != nil {
		return
	}

	collection := b.database.Collection(dbName)

	filter := bson.D{{"_id", id}}

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

	redfishService, err := b.inventoryService(ctx, redfishDevice, redfishClient.Service)
	if err != nil {
		return
	}

	response = &pbinventory.Response{Id: redfishService.Id.Hex()}
	return
}
