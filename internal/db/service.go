// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/stmcginnis/gofish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishService struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	DeviceId primitive.ObjectID `bson:"_device_id,omitempty"`
	*gofish.Service
}
