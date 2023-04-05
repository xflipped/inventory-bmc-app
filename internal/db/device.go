// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishDevice struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Url  string             `bson:"url"`
	UUID string             `bson:"UUID"`
}

func (d *RedfishDevice) ToProto() (device *pbredfish.Device, err error) {
	device = &pbredfish.Device{Id: d.Id.Hex(), Url: d.Url}
	return
}
