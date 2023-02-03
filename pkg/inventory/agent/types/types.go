// Copyright 2023 NJWS Inc.

package types

const (
	App       = "inventory-bmc"
	Namespace = "proxy.foliage"

	// FunctionType will be as default qdsl to function
	FunctionType = FunctionPath
	Description  = "inventory-bmc init function"
)

const (
	BmcContainerID      = "types/bmc-container"
	FunctionContainerID = "types/function-container"
	FunctionID          = "types/function"

	RootID = "system/root"

	RedfishDeviceKey = "redfish-device"
)

const (
	BmcContainerLink      = "bmc"
	FunctionContainerLink = "inventory-bmc"
	FunctionLink          = "init"
)

const (
	BmcContainerPath      = "bmc.root"
	FunctionsPath         = "functions.root"
	FunctionContainerPath = "inventory-bmc.functions.root"
	FunctionPath          = "init.inventory-bmc.functions.root"
)
