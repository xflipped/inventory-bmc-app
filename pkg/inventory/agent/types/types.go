// Copyright 2023 NJWS Inc.

package types

const (
	App       = "inventory-bmc"
	Namespace = "proxy.foliage"

	// FunctionType will be as default qdsl to function
	FunctionType = FunctionPath
	Description  = "inventory redfish function"
)

const (
	FunctionContainerID = "types/function-container"
	FunctionID          = "types/function"

	RootID = "system/root"
)

const (
	FunctionContainerLink = "redfish"
	FunctionLink          = "inventory"

	RedfishDeviceKey = "redfish-device"
)

const (
	FunctionsPath         = "functions.root"
	FunctionContainerPath = "redfish.functions.root"
	FunctionPath          = "inventory.redfish.functions.root"
)
