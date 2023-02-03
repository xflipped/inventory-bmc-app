// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"

	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
)

type Request struct {
	Query string `json:"query"`
	Name  string `json:"name"`
}

func prepareFunc(id string, r Request) (fc *pbtypes.FunctionContext, err error) {
	ft := &pbtypes.FunctionType{
		Namespace: types.Namespace,
		Type:      types.FunctionPath,
	}

	fc = &pbtypes.FunctionContext{
		Id:           id,
		FunctionType: ft,
	}
	fc.Value, err = json.Marshal(r)
	return
}

// genFunction generate function call with object uuid and qdsl
func genFunction(id, query string) (*pbtypes.FunctionContext, error) {
	r := Request{Query: query}
	return prepareFunc(id, r)
}
