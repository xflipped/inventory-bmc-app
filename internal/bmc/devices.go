// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"
	"fmt"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"go.mongodb.org/mongo-driver/mongo"
)

const devicesColName = "devices"

var errDeviceNotFound = fmt.Errorf("device not found")

func (b *BmcApp) ListDevices(ctx context.Context, empty *pbbmc.Empty) (devices *pbredfish.Devices, err error) {
	cur, err := b.database.Collection(devicesColName).Aggregate(ctx, mongo.Pipeline{lookupService, lookupSystem, lookupManager, lookupChasseez, project})
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	var results []db.RedfishDevice
	if err = cur.All(ctx, &results); err != nil {
		return
	}

	devices = &pbredfish.Devices{}
	for _, result := range results {
		devices.Items = append(devices.Items, result.ToProto())
	}

	return
}
