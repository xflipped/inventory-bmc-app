// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *BmcApp) inventoryService(ctx context.Context, redfishDevice db.RedfishDevice, service *gofish.Service) (redfishService db.RedfishService, err error) {
	const dbName = "services"

	redfishService = db.RedfishService{
		DeviceId: redfishDevice.Id,
		Service:  service,
	}

	filter := bson.D{{"_device_id", redfishService.DeviceId}}

	update := bson.D{{"$set", redfishService}}

	collection := b.database.Collection(dbName)

	singleResult := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))
	if err = singleResult.Err(); err != nil {
		return
	}

	if err = singleResult.Decode(&redfishService); err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return b.inventorySystems(ctx, redfishService) })
	err = p.Wait()
	return
}
