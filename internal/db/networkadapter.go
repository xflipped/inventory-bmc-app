// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishNetworkAdapter struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ChassisId primitive.ObjectID `bson:"_chassis_id,omitempty"`
	*redfish.NetworkAdapter
}