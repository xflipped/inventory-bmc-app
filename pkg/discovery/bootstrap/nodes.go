// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types"
)

var (
	registerObjects = []*pbcmdb.RegisterObjectMessage{}
)

type RedfishFunctionContainer struct{}

type RedfishDeviceContainer struct{}

func createObjects(ctx context.Context) (err error) {
	if err = createRedfishDeviceContainerObject(ctx); err != nil {
		return
	}

	if err = createDiscoveryRedFishMountpointObject(ctx); err != nil {
		return
	}

	if err = createDiscoveryFunctionObject(ctx); err != nil {
		return
	}

	return
}

func createRedfishDeviceContainerObject(ctx context.Context) (err error) {
	// check if object exists
	elements, err := qdsl.Qdsl(ctx, types.RedfishDevicesPath)
	if err != nil {
		return
	}

	// TODO already exists
	if len(elements) > 0 {
		return
	}

	message, err := system.RegisterObject(types.RootID, types.RedfishDevicesContainerID, types.RedfishDevicesLink, RedfishDeviceContainer{}, true, false)
	if err != nil {
		return
	}
	registerObjects = append(registerObjects, message)
	return
}

func createDiscoveryRedFishMountpointObject(ctx context.Context) (err error) {
	// check if object exists
	elements, err := qdsl.Qdsl(ctx, types.FunctionContainerPath)
	if err != nil {
		return
	}

	// TODO already exists
	if len(elements) > 0 {
		return
	}

	message, err := system.RegisterObject(types.FunctionsPath, types.FunctionContainerID, types.FunctionContainerLink, RedfishFunctionContainer{}, false, true)
	if err != nil {
		return
	}
	registerObjects = append(registerObjects, message)

	return
}
