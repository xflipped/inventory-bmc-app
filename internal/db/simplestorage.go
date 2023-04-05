// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishSimpleStorage struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	SystemId primitive.ObjectID `bson:"_system_id,omitempty"`
	*redfish.SimpleStorage
}
