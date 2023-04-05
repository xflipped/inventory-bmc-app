// Copyright 2023 NJWS Inc.

package bmc

import (
	"go.mongodb.org/mongo-driver/bson"
)

var (
	lookupService = bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: servicesColName},
				{Key: "localField", Value: "_id"},
				{Key: "foreignField", Value: "_device_id"},
				{Key: "as", Value: "service"},
			},
		},
	}

	lookupSystem = bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: systemsColName},
				{Key: "localField", Value: "service._id"},
				{Key: "foreignField", Value: "_service_id"},
				{Key: "as", Value: "system"},
			},
		},
	}

	lookupManager = bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: managersColName},
				{Key: "localField", Value: "service._id"},
				{Key: "foreignField", Value: "_service_id"},
				{Key: "as", Value: "manager"},
			},
		},
	}

	lookupChasseez = bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: chasseezColName},
				{Key: "localField", Value: "service._id"},
				{Key: "foreignField", Value: "_service_id"},
				{Key: "as", Value: "chassis"},
			},
		},
	}

	project = bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 1},

				{Key: "url", Value: 1},
				{Key: "uuid", Value: bson.D{{Key: "$first", Value: "$service.service.uuid"}}},
				{Key: "serial_number", Value: bson.D{{Key: "$first", Value: "$system.computersystem.serialnumber"}}},
				{Key: "name", Value: bson.D{{Key: "$first", Value: "$service.service.product"}}},

				// FIXME
				{Key: "mac_address", Value: bson.D{{Key: "$first", Value: "$manager.manager.mac"}}},

				{Key: "model", Value: bson.D{{Key: "$first", Value: "$manager.manager.model"}}},
				{Key: "vendor", Value: bson.D{{Key: "$first", Value: "$service.service.vendor"}}},
				{Key: "power_state", Value: bson.D{{Key: "$first", Value: "$manager.manager.powerstate"}}},
				{Key: "health_status", Value: bson.D{{Key: "$first", Value: "$system.computersystem.status.health"}}},
				{Key: "indicator_led", Value: bson.D{{Key: "$first", Value: "$chassis.chassis.indicatorled"}}},

				// FIXME
				{Key: "min_temp", Value: bson.D{{Key: "$first", Value: "$manager.manager.temp"}}},
				// FIXME
				{Key: "max_temp", Value: bson.D{{Key: "$first", Value: "$manager.manager.temp"}}},
			}},
	}
)
