// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson"
)

func (b *BmcApp) inventorySystems(ctx context.Context, redfishService db.RedfishService) (err error) {
	log.Infof("exec inventorySystems")

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
	log.Infof("exec inventorySystem")

	const colName = "systems"

	redfishSystem := db.RedfishSystem{
		ServiceId:      redfishService.Id,
		ComputerSystem: computerSystem,
	}

	filter := bson.D{{"_service_id", redfishSystem.ServiceId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishSystem); err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return b.inventoryBIOS(ctx, redfishSystem) })

	return p.Wait()
}

func (b *BmcApp) inventoryBIOS(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryBIOS")

	const colName = "bios"

	bios, err := redfishSystem.Bios()
	if err != nil {
		return
	}

	redfishBIOS := db.RedfishBIOS{
		SystemId: redfishSystem.Id,
		Bios:     bios,
	}

	filter := bson.D{{"_system_id", redfishBIOS.SystemId}}
	return b.FindOneAndReplace(ctx, colName, filter, &redfishBIOS)
}
