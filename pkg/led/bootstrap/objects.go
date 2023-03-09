// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"github.com/foliagecp/inventory-bmc-app/pkg/led/agent/types"
)

var (
	registerObjects = []*pbcmdb.RegisterObjectMessage{}
)

type RedfishFunctionContainer struct{}

func createObjects(ctx context.Context) (err error) {
	if err = createLedRedFishMountpointObject(ctx); err != nil {
		return
	}

	if err = createLedFunctionObject(ctx); err != nil {
		return
	}

	return
}

func createLedRedFishMountpointObject(ctx context.Context) (err error) {
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
