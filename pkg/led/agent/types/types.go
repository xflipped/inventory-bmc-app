// Copyright 2023 NJWS Inc.

package types

const (
	App       = "led-bmc"
	Namespace = "proxy.foliage"

	// FunctionType will be as default qdsl to function
	LedFunctionType = LedFunctionPath
	Description     = "update chassis indicator LED function"
)

const (
	FunctionContainerID = "types/function-container"
	FunctionID          = "types/function"
)

const (
	FunctionContainerLink = "led-bmc"
	LedFunctionLink       = "led"
)

const (
	FunctionsPath         = "functions.root"
	FunctionContainerPath = "led-bmc.functions.root"
	LedFunctionPath       = "led.led-bmc.functions.root"
)
