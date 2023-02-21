// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types"
)

func createDiscoveryFunctionObject(ctx context.Context) (err error) {
	// check if object exists
	elements, err := qdsl.Qdsl(ctx, types.FunctionPath)
	if err != nil {
		return
	}

	// already exists
	if len(elements) > 0 {
		return
	}

	function := pbtypes.Function{
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.FunctionType,
		},
		Description: types.Description,
	}

	message, err := system.RegisterObject(types.FunctionContainerPath, types.FunctionID, types.FunctionLink, function, true, true)
	if err != nil {
		return
	}

	registerObjects = append(registerObjects, message)
	return
}
