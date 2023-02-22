// Copyright 2023 NJWS Inc.

package types

const (
	App       = "inventory-bmc"
	Namespace = "proxy.foliage"

	// FunctionType will be as default qdsl to function
	InventoryFunctionType = InventoryFunctionPath
	Description           = "inventory redfish function"
)

const (
	FunctionContainerID          = "types/function-container"
	FunctionID                   = "types/function"
	RedfishServiceID             = "types/redfish-service"
	RedfishSystemID              = "types/redfish-system"
	RedfishBiosID                = "types/redfish-bios"
	RedfishLedID                 = "types/redfish-led"
	RedfishStatusID              = "types/redfish-status"
	RedfishBootID                = "types/redfish-boot"
	RedfishBootOptionID          = "types/redfish-boot-option"
	RedfishSecureBootOptionID    = "types/redfish-secure-boot"
	RedfishPcieDeviceID          = "types/redfish-pcie-device"
	RedfishPcieInterfaceID       = "types/redfish-pcie-interface"
	RedfishPcieFunctionID        = "types/redfish-pcie-function"
	RedfishPowerStateID          = "types/redfish-power-state"
	RedfishManagerID             = "types/redfish-manager"
	RedfishPowerID               = "types/redfish-power"
	RedfishCommandShellID        = "types/redfish-command-shell"
	RedfishEthernetInterfaceID   = "types/redfish-ethernet-interface"
	RedfishHostInterfaceID       = "types/redfish-host-interface"
	RedfishHostInterfaceTypeID   = "types/redfish-host-interface-type"
	RedfishChassisID             = "types/redfish-chassis"
	RedfishThermalID             = "types/redfish-thermal"
	RedfishLocationID            = "types/redfish-location"
	RedfishTemperatureID         = "types/redfish-temperature"
	RedfishFanID                 = "types/redfish-fan"
	RedfishPowerControlID        = "types/redfish-power-control"
	RedfishVoltageID             = "types/redfish-voltage"
	RedfishPowerSupplyID         = "types/redfish-power-supply"
	RedfishPartLocationID        = "types/redfish-part-location"
	RedfishPowerRestorePolicyID  = "types/redfish-power-restore-policy"
	RedfishProcessorSummaryID    = "types/redfish-processor-summary"
	RedfishProcessorID           = "types/redfish-processor"
	RedfishMemorySummaryID       = "types/redfish-memory-summary"
	RedfishMemoryID              = "types/redfish-memory"
	RedfishHostWatchdogTimerID   = "types/redfish-host-watchdog-timer"
	RedfishPhysicalSecurityID    = "types/redfish-physical-security"
	RedfishPlacementID           = "types/redfish-placement"
	RedfishPhysicalContextID     = "types/redfish-physical-context"
	RedfishPowerMetricID         = "types/redfish-power-metric"
	RedfishPowerLimitID          = "types/redfish-power-limit"
	RedfishPostalAddressID       = "types/redfish-postal-address"
	RedfishSupportedResetTypesID = "types/redfish-supported-reset-types"

	RootID = "system/root"
)

const (
	FunctionContainerLink = "inventory-bmc"
	InventoryFunctionLink = "inventory"

	RedfishServiceLink             = "service"
	RedfishBiosLink                = "bios"
	RedfishLedLink                 = "led"
	RedfishStatusLink              = "status"
	RedfishBootLink                = "boot"
	RedfishSecureBootLink          = "secure-boot"
	RedfishPcieInterfaceLink       = "pcie-interface"
	RedfishPowerStateLink          = "power-state"
	RedfishLocationLink            = "location"
	RedfishCommandShellLink        = "command-shell"
	RedfishTypeLink                = "type"
	RedfishThermalLink             = "thermal"
	RedfishPowerLink               = "power"
	RedfishDeviceKey               = "redfish-device"
	RedfishPartLocationLink        = "part-location"
	RedfishPowerRestorePolicyLink  = "power-restore-policy"
	RedfishProcessorSummaryLink    = "processor-summary"
	RedfishMemorySummaryLink       = "memory-summary"
	RedfishHostWatchdogTimerLink   = "host-watchdog-timer"
	RedfishPhysicalSecurityLink    = "physical-security"
	RedfishPlacementLink           = "placement"
	RedfishPhysicalContextLink     = "physical-context"
	RedfishPowerMetricLink         = "power-metric"
	RedfishPowerLimitLink          = "power-limit"
	RedfishPostalAddressLink       = "postal-address"
	RedfishSupportedResetTypesLink = "supported-reset-types"
)

const (
	FunctionsPath         = "functions.root"
	FunctionContainerPath = "inventory-bmc.functions.root"
	InventoryFunctionPath = "inventory.inventory-bmc.functions.root"
)
