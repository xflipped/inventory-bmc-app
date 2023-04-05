// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish"
	"go.mongodb.org/mongo-driver/bson"
)

const servicesColName = "services"

func (b *BmcApp) inventoryService(ctx context.Context, redfishDevice db.RedfishDevice, service *gofish.Service) (err error) {
	log.Infof("exec inventoryService")

	redfishService := db.RedfishService{
		DeviceId: redfishDevice.Id,
		Service:  service,
	}

	filter := bson.D{{Key: "_device_id", Value: redfishService.DeviceId}}
	if err = b.FindOneAndReplace(ctx, servicesColName, filter, &redfishService); err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return b.inventorySystems(ctx, redfishService) })
	p.Exec(func() error { return b.inventoryManagers(ctx, redfishService) })
	p.Exec(func() error { return b.inventoryChasseez(ctx, redfishService) })
	err = p.Wait()
	return
}
