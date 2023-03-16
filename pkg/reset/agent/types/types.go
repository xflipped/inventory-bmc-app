// Copyright 2023 NJWS Inc.

package types

const (
	App       = "reset-bmc"
	Namespace = "proxy.foliage"

	// FunctionType will be as default qdsl to function
	ResetFunctionType = ResetFunctionPath
	Description       = "reset system"
)

const (
	FunctionContainerID = "types/function-container"
	FunctionID          = "types/function"
)

const (
	FunctionContainerLink = "reset-bmc"
	ResetFunctionLink     = "reset"
)

const (
	FunctionsPath         = "functions.root"
	FunctionContainerPath = "reset-bmc.functions.root"
	ResetFunctionPath     = "reset.reset-bmc.functions.root"
)
