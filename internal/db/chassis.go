// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishChassis struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ServiceId primitive.ObjectID `bson:"_service_id,omitempty"`
	*redfish.Chassis
}
