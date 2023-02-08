// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex/types"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

var (
	registerTypes = []*pbcmdb.RegisterTypeMessage{}
)

// FIXME will be custom structs?
// Root Service
type RedfishService struct {
	*gofish.Service
}

// Core structs - system, chassis
type RedfishSystem struct {
	*redfish.ComputerSystem
}

type RedfishChassis struct {
	*redfish.Chassis
}

// Nested structs
type RedfishBios struct {
	*redfish.Bios
}

type RedfishLed struct {
	Led common.IndicatorLED `json:"led"`
}

type RedfishStatus struct {
	Status common.Status `json:"status"`
}

type RedfishBoot struct {
	Boot *redfish.Boot `json:"boot"`
}

type RedfishBootOption struct {
	BootOption *redfish.BootOption `json:"bootOption"`
}

type RedfishSecureBoot struct {
	SecureBoot *redfish.SecureBoot `json:"secureBoot"`
}

type RedfishPcieDevice struct {
	PCIeDevice *redfish.PCIeDevice `json:"PCIeDevice"`
}

type RedfishPcieInterface struct {
	PCIeInterface *redfish.PCIeInterface `json:"PCIeInterface"`
}

type RedfishThermal struct {
	ID          string `json:"Id,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
}

type RedfishTemperature struct {
	Temperature *redfish.Temperature `json:"temperature"`
}

type RedfishFan struct {
	Fan *redfish.Fan `json:"fan"`
}

type RedfishPower struct {
	ID          string `json:"Id,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
}

type RedfishPowerControl struct {
	PowerControl *redfish.PowerControl `json:"powerControl"`
}

type RedfishPowerSupply struct {
	PowerSupply *redfish.PowerSupply `json:"powerSupply"`
}

type RedfishVoltage struct {
	Voltage *redfish.Voltage `json:"voltage"`
}

func createType(ctx context.Context, pt *types.Type) (err error) {
	query := fmt.Sprintf("%s.types.root", pt.Schema.Title)
	elements, err := qdsl.Qdsl(ctx, query)
	if err != nil {
		return
	}

	// TODO already exists
	if len(elements) > 0 {
		return
	}

	message, err := system.RegisterType(pt, true)
	if err != nil {
		return
	}

	registerTypes = append(registerTypes, message)
	return
}

func createRedfishServiceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishService{})
	return createType(ctx, pt)
}

func createRedfishSystemType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishSystem{})
	return createType(ctx, pt)
}

func createRedfishChassisType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishChassis{})
	return createType(ctx, pt)
}

func createRedfishBiosType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishBios{})
	return createType(ctx, pt)
}

func createRedfishLedType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishLed{})
	return createType(ctx, pt)
}

func createRedfishStatus(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishStatus{})
	return createType(ctx, pt)
}

func createRedfishBoot(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishBoot{})
	return createType(ctx, pt)
}

func createRedfishBootOption(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishBootOption{})
	return createType(ctx, pt)
}

func createRedfishSecureBoot(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishSecureBoot{})
	return createType(ctx, pt)
}

func createRedfishPCIeDevice(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPcieDevice{})
	return createType(ctx, pt)
}

func createRedfishPCIeInterface(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPcieInterface{})
	return createType(ctx, pt)
}

func createRedfishThermal(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishThermal{})
	return createType(ctx, pt)
}

func createRedfishTemperature(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishTemperature{})
	return createType(ctx, pt)
}

func createRedfishFan(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishFan{})
	return createType(ctx, pt)
}

func createRedfishPower(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPower{})
	return createType(ctx, pt)
}

func createRedfishPowerControl(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerControl{})
	return createType(ctx, pt)
}

func createRedfishPowerSupply(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerSupply{})
	return createType(ctx, pt)
}

func createRedfishVoltage(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishVoltage{})
	return createType(ctx, pt)
}

func createTypes(ctx context.Context) (err error) {
	if err = createRedfishServiceType(ctx); err != nil {
		return
	}
	if err = createRedfishSystemType(ctx); err != nil {
		return
	}
	if err = createRedfishChassisType(ctx); err != nil {
		return
	}
	if err = createRedfishBiosType(ctx); err != nil {
		return
	}
	if err = createRedfishLedType(ctx); err != nil {
		return
	}
	if err = createRedfishStatus(ctx); err != nil {
		return
	}
	if err = createRedfishBoot(ctx); err != nil {
		return
	}
	if err = createRedfishBootOption(ctx); err != nil {
		return
	}
	if err = createRedfishSecureBoot(ctx); err != nil {
		return
	}
	if err = createRedfishPCIeDevice(ctx); err != nil {
		return
	}
	if err = createRedfishPCIeInterface(ctx); err != nil {
		return
	}
	if err = createRedfishThermal(ctx); err != nil {
		return
	}
	if err = createRedfishTemperature(ctx); err != nil {
		return
	}
	if err = createRedfishFan(ctx); err != nil {
		return
	}
	if err = createRedfishPower(ctx); err != nil {
		return
	}
	if err = createRedfishPowerControl(ctx); err != nil {
		return
	}
	if err = createRedfishPowerSupply(ctx); err != nil {
		return
	}
	if err = createRedfishVoltage(ctx); err != nil {
		return
	}
	return
}
