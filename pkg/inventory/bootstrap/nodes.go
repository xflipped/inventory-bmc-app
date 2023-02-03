// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
)

var (
	registerObjects = []*pbcmdb.RegisterObjectMessage{}
)

type InventoryFunctionContainer struct{}

type BmcContainer struct{}

func createObjects(ctx context.Context) (err error) {
	if err = createBmcContainerObject(ctx); err != nil {
		return
	}

	if err = createInventoryBmcMountpointObject(ctx); err != nil {
		return
	}

	if err = createInitFunctionObject(ctx); err != nil {
		return
	}

	return
}

func createBmcContainerObject(ctx context.Context) (err error) {
	// check if object exists
	elements, err := qdsl.Qdsl(ctx, types.BmcContainerPath)
	if err != nil {
		return
	}

	// TODO already exists
	if len(elements) > 0 {
		return
	}

	message, err := system.RegisterObject(types.RootID, types.BmcContainerID, types.BmcContainerLink, BmcContainer{}, true, false)
	if err != nil {
		return
	}
	registerObjects = append(registerObjects, message)
	return
}

func createInventoryBmcMountpointObject(ctx context.Context) (err error) {
	// check if object exists
	elements, err := qdsl.Qdsl(ctx, types.FunctionContainerPath)
	if err != nil {
		return
	}

	// TODO already exists
	if len(elements) > 0 {
		return
	}

	message, err := system.RegisterObject(types.FunctionsPath, types.FunctionContainerID, types.FunctionContainerLink, InventoryFunctionContainer{}, false, true)
	if err != nil {
		return
	}
	registerObjects = append(registerObjects, message)

	return
}
