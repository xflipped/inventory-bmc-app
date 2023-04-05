// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishStorageVolume struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	StorageId primitive.ObjectID `bson:"_storage_id,omitempty"`
	*redfish.Volume
}
