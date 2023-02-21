// Copyright 2023 NJWS Inc.

package agent

import (
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
)

func PrepareInventoryFunc(id string) (fc *pbtypes.FunctionContext, err error) {
	fc = &pbtypes.FunctionContext{
		Id: id,
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.InventoryFunctionPath,
		},
	}
	return
}
