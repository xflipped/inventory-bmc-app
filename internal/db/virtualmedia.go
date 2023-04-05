// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishVirtualMedia struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	ManagerId primitive.ObjectID `bson:"_manager_id,omitempty"`
	*redfish.VirtualMedia
}
