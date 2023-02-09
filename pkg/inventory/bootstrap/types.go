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

// Core structs - system, chassis, manager
type RedfishSystem struct {
	*redfish.ComputerSystem
}

type RedfishChassis struct {
	*redfish.Chassis
}

type RedfishManager struct {
	*redfish.Manager
}

// Nested structs
type RedfishBios struct {
	*redfish.Bios
}

type RedfishBiosAttribute struct {
	BiosAttributeValue string `json:"biosAttributeValue"`
}

type RedfishLed struct {
	Led common.IndicatorLED `json:"led"`
}

type RedfishStatus struct {
	Status common.Status `json:"status"`
}

type RedfishBoot struct {
	*redfish.Boot
}

type RedfishBootOption struct {
	*redfish.BootOption
}

type RedfishSecureBoot struct {
	*redfish.SecureBoot
}

type RedfishPcieDevice struct {
	*redfish.PCIeDevice
}

type RedfishPcieInterface struct {
	*redfish.PCIeInterface
}

type RedfishPowerState struct {
	PowerState redfish.PowerState `json:"powerState"`
}

type RedfishPowerRestorePolicy struct {
	PowerRestorePolicy redfish.PowerState `json:"powerRestorePolicy"`
}

type RedfishProcessorSummary struct {
	*redfish.ProcessorSummary
}

type RedfishMemorySummary struct {
	*redfish.MemorySummary
}

type RedfishHostWatchdogTimer struct {
	*redfish.WatchdogTimer
}

type RedfishThermal struct {
	ID          string `json:"Id,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
}

type RedfishTemperature struct {
	*redfish.Temperature
}

type RedfishFan struct {
	*redfish.Fan
}

type RedfishPower struct {
	ID          string `json:"Id,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Description,omitempty"`
}

type RedfishPowerControl struct {
	*redfish.PowerControl
}

type RedfishPowerSupply struct {
	*redfish.PowerSupply
}

type RedfishVoltage struct {
	*redfish.Voltage
}

type RedfishCommandShell struct {
	*redfish.CommandShell
}

type RedfishEthernetInterface struct {
	*redfish.EthernetInterface
}

type RedfishHostInterface struct {
	*redfish.HostInterface
}

type RedfishHostInterfaceType struct {
	HostInterfaceType *redfish.HostInterfaceType `json:"hostInterfaceType"`
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

func createRedfishManagerType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishManager{})
	return createType(ctx, pt)
}

func createRedfishBiosType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishBios{})
	return createType(ctx, pt)
}

func createRedfishBiosAttributeType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishBiosAttribute{})
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

func createRedfishPowerState(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerState{})
	return createType(ctx, pt)
}

func createRedfishPowerRestorePolicy(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerRestorePolicy{})
	return createType(ctx, pt)
}

func createRedfishProcessorSummary(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishProcessorSummary{})
	return createType(ctx, pt)
}

func createRedfishMemorySummary(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishMemorySummary{})
	return createType(ctx, pt)
}

func createRedfishHostWatchdogTimer(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishHostWatchdogTimer{})
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

func createRedfishCommandShell(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishCommandShell{})
	return createType(ctx, pt)
}

func createRedfishEthernetInterface(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishEthernetInterface{})
	return createType(ctx, pt)
}

func createRedfishHostInterface(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishHostInterface{})
	return createType(ctx, pt)
}

func createRedfishHostInterfaceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishHostInterfaceType{})
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
	if err = createRedfishManagerType(ctx); err != nil {
		return
	}
	if err = createRedfishBiosType(ctx); err != nil {
		return
	}
	if err = createRedfishBiosAttributeType(ctx); err != nil {
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
	if err = createRedfishPowerState(ctx); err != nil {
		return
	}
	if err = createRedfishPowerRestorePolicy(ctx); err != nil {
		return
	}
	if err = createRedfishProcessorSummary(ctx); err != nil {
		return
	}
	if err = createRedfishMemorySummary(ctx); err != nil {
		return
	}
	if err = createRedfishHostWatchdogTimer(ctx); err != nil {
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
	if err = createRedfishCommandShell(ctx); err != nil {
		return
	}
	if err = createRedfishEthernetInterface(ctx); err != nil {
		return
	}
	if err = createRedfishHostInterface(ctx); err != nil {
		return
	}
	if err = createRedfishHostInterfaceType(ctx); err != nil {
		return
	}
	return
}
