// Copyright 2023 NJWS Inc.

package types

const (
	App       = "subscribe-bmc"
	Namespace = "proxy.foliage"

	// FunctionType will be as default qdsl to function
	SubscribeFunctionType = SubscribeFunctionPath
	Description           = "subcribe to BMC events function"
)

const (
	FunctionContainerID = "types/function-container"
	FunctionID          = "types/function"
)

const (
	FunctionContainerLink = "subscribe-bmc"
	SubscribeFunctionLink = "subscribe"
)

const (
	FunctionsPath         = "functions.root"
	FunctionContainerPath = "subscribe-bmc.functions.root"
	SubscribeFunctionPath = "subscribe.subscribe-bmc.functions.root"
)
