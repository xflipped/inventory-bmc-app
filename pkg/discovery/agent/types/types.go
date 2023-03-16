// Copyright 2023 NJWS Inc.

package types

const (
	App       = "discovery-bmc"
	Namespace = "proxy.foliage"

	// FunctionType will be as default qdsl to function
	DiscoveryFunctionType = DiscoveryFunctionPath
	Description           = "discovery redfish function"
)

const (
	RedfishDevicesContainerID = "types/redfish-device-container"
	FunctionContainerID       = "types/function-container"
	FunctionID                = "types/function"

	RedfishDeviceID = "types/redfish-device"

	RootID = "system/root"
)

const (
	RedfishDevicesLink    = "redfish-devices"
	FunctionContainerLink = "discovery-bmc"
	DiscoveryFunctionLink = "discovery"
)

const (
	RedfishDevicesPath    = "redfish-devices.root"
	FunctionsPath         = "functions.root"
	FunctionContainerPath = "discovery-bmc.functions.root"
	DiscoveryFunctionPath = "discovery.discovery-bmc.functions.root"
)
