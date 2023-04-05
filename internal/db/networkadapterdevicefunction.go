// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishNetworkAdapterDeviceFunction struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	NetworkAdapterId primitive.ObjectID `bson:"_network_adapter_id,omitempty"`
	*redfish.NetworkDeviceFunction
}
