// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson"
)

func (b *BmcApp) inventoryChasseez(ctx context.Context, redfishService db.RedfishService) (err error) {
	log.Infof("exec inventoryChasseez")

	chasseez, err := redfishService.Chassis()
	if err != nil {
		return
	}
	p := utils.NewParallel()
	for _, chassis := range chasseez {
		chassis := chassis
		p.Exec(func() error { return b.inventoryChassis(ctx, redfishService, chassis) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryChassis(ctx context.Context, redfishService db.RedfishService, chassis *redfish.Chassis) (err error) {
	log.Infof("exec inventoryChassis")

	const colName = "chasseez"

	redfishChassis := db.RedfishChassis{
		ServiceId: redfishService.Id,
		Chassis:   chassis,
	}

	filter := bson.D{{"_service_id", redfishChassis.ServiceId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishChassis); err != nil {
		return
	}

	return
}