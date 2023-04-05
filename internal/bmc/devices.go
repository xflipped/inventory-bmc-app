// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	lookupService = bson.D{
		{"$lookup",
			bson.D{
				{"from", "services"},
				{"localField", "_id"},
				{"foreignField", "_device_id"},
				{"as", "service"},
			},
		},
	}

	lookupSystem = bson.D{
		{"$lookup",
			bson.D{
				{"from", "systems"},
				{"localField", "service._id"},
				{"foreignField", "_service_id"},
				{"as", "system"},
			},
		},
	}

	lookupManager = bson.D{
		{"$lookup",
			bson.D{
				{"from", "managers"},
				{"localField", "service._id"},
				{"foreignField", "_service_id"},
				{"as", "manager"},
			},
		},
	}

	lookupChasseez = bson.D{
		{"$lookup",
			bson.D{
				{"from", "chasseez"},
				{"localField", "service._id"},
				{"foreignField", "_service_id"},
				{"as", "chassis"},
			},
		},
	}

	project = bson.D{
		{"$project", bson.D{
			{"_id", 1},

			{"url", 1},
			{"uuid", bson.D{{"$first", "$service.service.uuid"}}},
			{"serial_number", bson.D{{"$first", "$system.computersystem.serialnumber"}}},
			{"name", bson.D{{"$first", "$service.service.product"}}},

			// FIXME
			{"mac_address", bson.D{{"$first", "$manager.manager.mac"}}},

			{"model", bson.D{{"$first", "$manager.manager.model"}}},
			{"vendor", bson.D{{"$first", "$service.service.vendor"}}},
			{"power_state", bson.D{{"$first", "$manager.manager.powerstate"}}},
			{"health_status", bson.D{{"$first", "$system.computersystem.status.health"}}},
			{"indicator_led", bson.D{{"$first", "$chassis.chassis.indicatorled"}}},

			// FIXME
			{"min_temp", bson.D{{"$first", "$manager.manager.temp"}}},
			// FIXME
			{"max_temp", bson.D{{"$first", "$manager.manager.temp"}}},
		}},
	}
)

func (b *BmcApp) ListDevices(ctx context.Context, empty *pbbmc.Empty) (devices *pbredfish.Devices, err error) {
	const colName = "devices"

	devices = &pbredfish.Devices{}

	cur, err := b.database.Collection(colName).Aggregate(ctx, mongo.Pipeline{lookupService, lookupSystem, lookupManager, lookupChasseez, project})
	if err != nil {
		return
	}
	defer cur.Close(ctx)

	var results []db.RedfishDevice
	if err = cur.All(ctx, &results); err != nil {
		return
	}

	for _, result := range results {
		devices.Items = append(devices.Items, result.ToProto())
	}

	return
}
