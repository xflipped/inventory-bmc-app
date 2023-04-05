// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishNetworkDeviceFunction struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	NetworkInterfaceId primitive.ObjectID `bson:"_network_interface_id,omitempty"`
	*redfish.NetworkDeviceFunction
}
