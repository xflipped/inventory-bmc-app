// Copyright 2023 NJWS Inc.

package db

import (
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RedfishDevice struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Url          string             `bson:"url,omitempty"`
	UUID         string             `bson:"uuid,omitempty"`
	SerialNumber string             `bson:"serial_number,omitempty"`
	Name         string             `bson:"name,omitempty"`
	MacAddress   string             `bson:"mac_address,omitempty"`
	Model        string             `bson:"model,omitempty"`
	Vendor       string             `bson:"vendor,omitempty"`
	PowerState   string             `bson:"power_state,omitempty"`
	HealthStatus string             `bson:"health_status,omitempty"`
	IndicatorLed string             `bson:"indicator_led,omitempty"`
	MinTemp      string             `bson:"min_temp,omitempty"`
	MaxTemp      string             `bson:"max_temp,omitempty"`
}

func (d *RedfishDevice) ToProto() (device *pbredfish.Device, err error) {
	device = &pbredfish.Device{
		Id:           d.Id.Hex(),
		Url:          d.Url,
		UUID:         d.UUID,
		SerialNumber: d.SerialNumber,
		Name:         d.Name,
		MacAddress:   d.MacAddress,
		Model:        d.Model,
		Vendor:       d.Vendor,
		PowerState:   d.PowerState,
		HealthStatus: d.HealthStatus,
		IndicatorLed: d.IndicatorLed,
		MinTemp:      d.MinTemp,
		MaxTemp:      d.MaxTemp,
	}
	return
}
