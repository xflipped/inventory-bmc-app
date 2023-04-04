// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *BmcApp) inventorySystems(ctx context.Context, redfishService db.RedfishService) (err error) {
	systems, err := redfishService.Systems()
	if err != nil {
		return
	}

	p := utils.NewParallel()

	for _, system := range systems {
		system := system
		p.Exec(func() error { return b.inventorySystem(ctx, redfishService, system) })
	}

	err = p.Wait()
	return
}

func (b *BmcApp) inventorySystem(ctx context.Context, redfishService db.RedfishService, computerSystem *redfish.ComputerSystem) (err error) {
	const dbName = "systems"

	redfishSystem := db.RedfishSystem{
		ServiceId:      redfishService.Id,
		ComputerSystem: computerSystem,
	}

	filter := bson.D{{"_service_id", redfishSystem.ServiceId}}

	update := bson.D{{"$set", redfishSystem}}

	collection := b.database.Collection(dbName)

	singleResult := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))
	if err = singleResult.Err(); err != nil {
		return
	}

	if err = singleResult.Decode(&redfishSystem); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryBIOS(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	const dbName = "bios"

	bios, err := redfishSystem.Bios()
	if err != nil {
		return
	}

	redfishBIOS := db.RedfishBIOS{
		SystemId: redfishSystem.Id,
		Bios:     bios,
	}

	filter := bson.D{{"_system_id", redfishBIOS.SystemId}}

	update := bson.D{{"$set", redfishBIOS}}

	collection := b.database.Collection(dbName)

	singleResult := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After))
	if err = singleResult.Err(); err != nil {
		return
	}

	if err = singleResult.Decode(&redfishBIOS); err != nil {
		return
	}

	return
}
