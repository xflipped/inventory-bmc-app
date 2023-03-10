// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/reset/agent/types"
)

func createResetFunctionObject(ctx context.Context) (err error) {
	// check if object exists
	elements, err := qdsl.Qdsl(ctx, types.ResetFunctionPath)
	if err != nil {
		return
	}

	// already exists
	if len(elements) > 0 {
		return
	}

	function := &pbtypes.Function{
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.ResetFunctionType,
		},
		Description: types.Description,
		Grounded:    false,
	}

	message, err := system.RegisterObject(types.FunctionContainerPath, types.FunctionID, types.ResetFunctionLink, function, true, true)
	if err != nil {
		return
	}

	registerObjects = append(registerObjects, message)
	return
}
