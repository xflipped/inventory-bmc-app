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

type RedfishPcieFunction struct {
	*redfish.PCIeFunction
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

type RedfishProcessor struct {
	*redfish.Processor
}

type RedfishMemorySummary struct {
	*redfish.MemorySummary
}

type RedfishMemory struct {
	*redfish.Memory
}

type RedfishMemoryDomain struct {
	*redfish.MemoryDomain
}

type RedfishHostWatchdogTimer struct {
	*redfish.WatchdogTimer
}

type RedfishSimpleStorage struct {
	*redfish.SimpleStorage
}

type RedfishStorage struct {
	*redfish.Storage
}

type RedfishStorageDevice struct {
	*redfish.Device
}

type RedfishDrive struct {
	*redfish.Drive
}

type RedfishVolume struct {
	*redfish.Volume
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

type RedfishPhysicalContext struct {
	PhysicalContext *common.PhysicalContext `json:"physicalContext"`
}

type RedfishPowerMetric struct {
	*redfish.PowerMetric
}

type RedfishPowerLimit struct {
	*redfish.PowerLimit
}

type RedfishPowerSupply struct {
	*redfish.PowerSupply
}

type RedfishVoltage struct {
	*redfish.Voltage
}

type RedfishPhysicalSecurity struct {
	*redfish.PhysicalSecurity
}

type RedfishLocation struct {
	*common.Location
}

type RedfishPartLocation struct {
	*common.PartLocation
}

type RedfishPlacement struct {
	*common.Placement
}

type RedfishPostalAddress struct {
	*common.PostalAddress
}

type RedfishSupportedResetTypes struct {
	SupportedResetTypes []redfish.ResetType `json:"resetTypes"`
}

type RedfishCommandShell struct {
	*redfish.CommandShell
}

type RedfishNetworkInterface struct {
	*redfish.NetworkInterface
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

type RedfishNetworkAdapter struct {
	*redfish.NetworkAdapter
}

type RedfishNetworkDeviceFunction struct {
	*redfish.NetworkDeviceFunction
}

type RedfishNetworkPort struct {
	*redfish.NetworkPort
}

type RedfishLogService struct {
	*redfish.LogService
}

type RedfishLogEntry struct {
	*redfish.LogEntry
}

type RedfishEventService struct {
	*redfish.EventService
}

type RedfishEventDestination struct {
	*redfish.EventDestination
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

func createRedfishStatusType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishStatus{})
	return createType(ctx, pt)
}

func createRedfishBootType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishBoot{})
	return createType(ctx, pt)
}

func createRedfishBootOptionType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishBootOption{})
	return createType(ctx, pt)
}

func createRedfishSecureBootType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishSecureBoot{})
	return createType(ctx, pt)
}

func createRedfishPCIeDeviceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPcieDevice{})
	return createType(ctx, pt)
}

func createRedfishPCIeFunctionType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPcieFunction{})
	return createType(ctx, pt)
}

func createRedfishPCIeInterfaceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPcieInterface{})
	return createType(ctx, pt)
}

func createRedfishPowerStateType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerState{})
	return createType(ctx, pt)
}

func createRedfishPowerRestorePolicyType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerRestorePolicy{})
	return createType(ctx, pt)
}

func createRedfishProcessorSummaryType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishProcessorSummary{})
	return createType(ctx, pt)
}

func createRedfishProcessorType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishProcessor{})
	return createType(ctx, pt)
}

func createRedfishMemorySummaryType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishMemorySummary{})
	return createType(ctx, pt)
}

func createRedfishMemoryType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishMemory{})
	return createType(ctx, pt)
}

func createRedfishMemoryDomainType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishMemoryDomain{})
	return createType(ctx, pt)
}

func createRedfishHostWatchdogTimerType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishHostWatchdogTimer{})
	return createType(ctx, pt)
}

func createRedfishSimpleStorageType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishSimpleStorage{})
	return createType(ctx, pt)
}

func createRedfishStorageType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishStorage{})
	return createType(ctx, pt)
}

func createRedfishStorageDeviceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishStorageDevice{})
	return createType(ctx, pt)
}

func createRedfishDriveType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishDrive{})
	return createType(ctx, pt)
}

func createRedfishVolumeType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishVolume{})
	return createType(ctx, pt)
}

func createRedfishThermalType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishThermal{})
	return createType(ctx, pt)
}

func createRedfishTemperatureType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishTemperature{})
	return createType(ctx, pt)
}

func createRedfishFanType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishFan{})
	return createType(ctx, pt)
}

func createRedfishPowerType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPower{})
	return createType(ctx, pt)
}

func createRedfishPowerControlType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerControl{})
	return createType(ctx, pt)
}

func createRedfishPhysicalContextType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPhysicalContext{})
	return createType(ctx, pt)
}

func createRedfishPowerMetricType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerMetric{})
	return createType(ctx, pt)
}

func createRedfishPowerLimitType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerLimit{})
	return createType(ctx, pt)
}

func createRedfishPowerSupplyType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPowerSupply{})
	return createType(ctx, pt)
}

func createRedfishVoltageType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishVoltage{})
	return createType(ctx, pt)
}

func createRedfishPhysicalSecurityType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPhysicalSecurity{})
	return createType(ctx, pt)
}

func createRedfishLocationType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishLocation{})
	return createType(ctx, pt)
}

func createRedfishPartLocationType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPartLocation{})
	return createType(ctx, pt)
}

func createRedfishPlacementType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPlacement{})
	return createType(ctx, pt)
}

func createRedfishPostalAddressType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishPostalAddress{})
	return createType(ctx, pt)
}

func createRedfishSupportedResetTypesType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishSupportedResetTypes{})
	return createType(ctx, pt)
}

func createRedfishCommandShellType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishCommandShell{})
	return createType(ctx, pt)
}

func createRedfishNetworkInterfaceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishNetworkInterface{})
	return createType(ctx, pt)
}

func createRedfishEthernetInterfaceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishEthernetInterface{})
	return createType(ctx, pt)
}

func createRedfishHostInterfaceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishHostInterface{})
	return createType(ctx, pt)
}

func createRedfishHostInterfaceTypeType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishHostInterfaceType{})
	return createType(ctx, pt)
}

func createRedfishNetworkAdapterType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishNetworkAdapter{})
	return createType(ctx, pt)
}

func createRedfishNetworkDeviceFunctionType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishNetworkDeviceFunction{})
	return createType(ctx, pt)
}

func createRedfishNetworkPortType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishNetworkPort{})
	return createType(ctx, pt)
}

func createRedfishLogServiceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishLogService{})
	return createType(ctx, pt)
}

func createRedfishLogEntryType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishLogEntry{})
	return createType(ctx, pt)
}

func createRedfishEventServiceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishEventService{})
	return createType(ctx, pt)
}

func createRedfishEventDestinationType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishEventDestination{})
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
	if err = createRedfishStatusType(ctx); err != nil {
		return
	}
	if err = createRedfishBootType(ctx); err != nil {
		return
	}
	if err = createRedfishBootOptionType(ctx); err != nil {
		return
	}
	if err = createRedfishSecureBootType(ctx); err != nil {
		return
	}
	if err = createRedfishPCIeDeviceType(ctx); err != nil {
		return
	}
	if err = createRedfishPCIeFunctionType(ctx); err != nil {
		return
	}
	if err = createRedfishPCIeInterfaceType(ctx); err != nil {
		return
	}
	if err = createRedfishPowerStateType(ctx); err != nil {
		return
	}
	if err = createRedfishPowerRestorePolicyType(ctx); err != nil {
		return
	}
	if err = createRedfishProcessorSummaryType(ctx); err != nil {
		return
	}
	if err = createRedfishProcessorType(ctx); err != nil {
		return
	}
	if err = createRedfishMemorySummaryType(ctx); err != nil {
		return
	}
	if err = createRedfishMemoryType(ctx); err != nil {
		return
	}
	if err = createRedfishMemoryDomainType(ctx); err != nil {
		return
	}
	if err = createRedfishHostWatchdogTimerType(ctx); err != nil {
		return
	}
	if err = createRedfishSimpleStorageType(ctx); err != nil {
		return
	}
	if err = createRedfishStorageType(ctx); err != nil {
		return
	}
	if err = createRedfishStorageDeviceType(ctx); err != nil {
		return
	}
	if err = createRedfishDriveType(ctx); err != nil {
		return
	}
	if err = createRedfishVolumeType(ctx); err != nil {
		return
	}
	if err = createRedfishThermalType(ctx); err != nil {
		return
	}
	if err = createRedfishTemperatureType(ctx); err != nil {
		return
	}
	if err = createRedfishFanType(ctx); err != nil {
		return
	}
	if err = createRedfishPowerType(ctx); err != nil {
		return
	}
	if err = createRedfishPowerControlType(ctx); err != nil {
		return
	}
	if err = createRedfishPhysicalContextType(ctx); err != nil {
		return
	}
	if err = createRedfishPowerMetricType(ctx); err != nil {
		return
	}
	if err = createRedfishPowerLimitType(ctx); err != nil {
		return
	}
	if err = createRedfishPowerSupplyType(ctx); err != nil {
		return
	}
	if err = createRedfishVoltageType(ctx); err != nil {
		return
	}
	if err = createRedfishPhysicalSecurityType(ctx); err != nil {
		return
	}
	if err = createRedfishLocationType(ctx); err != nil {
		return
	}
	if err = createRedfishPartLocationType(ctx); err != nil {
		return
	}
	if err = createRedfishPlacementType(ctx); err != nil {
		return
	}
	if err = createRedfishPostalAddressType(ctx); err != nil {
		return
	}
	if err = createRedfishSupportedResetTypesType(ctx); err != nil {
		return
	}
	if err = createRedfishCommandShellType(ctx); err != nil {
		return
	}
	if err = createRedfishNetworkInterfaceType(ctx); err != nil {
		return
	}
	if err = createRedfishEthernetInterfaceType(ctx); err != nil {
		return
	}
	if err = createRedfishHostInterfaceType(ctx); err != nil {
		return
	}
	if err = createRedfishHostInterfaceTypeType(ctx); err != nil {
		return
	}
	if err = createRedfishNetworkAdapterType(ctx); err != nil {
		return
	}
	if err = createRedfishNetworkDeviceFunctionType(ctx); err != nil {
		return
	}
	if err = createRedfishNetworkPortType(ctx); err != nil {
		return
	}
	if err = createRedfishLogServiceType(ctx); err != nil {
		return
	}
	if err = createRedfishLogEntryType(ctx); err != nil {
		return
	}
	if err = createRedfishEventServiceType(ctx); err != nil {
		return
	}
	if err = createRedfishEventDestinationType(ctx); err != nil {
		return
	}
	return
}
