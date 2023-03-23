// Copyright 2023 NJWS Inc.

package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/subscribe/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
)

func createOrUpdateFunctionLink(ctx context.Context, fromQuery, toQuery, name string) (functionContext *pbtypes.FunctionContext, err error) {
	route := &pbtypes.FunctionRoute{
		Url: fmt.Sprintf("http://%s:31005/statefun", types.App),
	}

	query := fmt.Sprintf("%s.%s", name, fromQuery)

	if linkDocument, err := utils.GetDocument(ctx, query); err == nil {
		return system.UpdateAdvancedLink(linkDocument.LinkId.String(), route)
	}

	parent, err := utils.GetDocument(ctx, fromQuery)
	if err != nil {
		return
	}

	child, err := utils.GetDocument(ctx, toQuery)
	if err != nil {
		return
	}

	return system.CreateLink(parent.Id.String(), child.Id.String(), name, "function", route)
}

func (a *Agent) createOrUpdateFunctionLink(fromQuery, toQuery, name string) (err error) {
	functionContext, err := createOrUpdateFunctionLink(a.ctx, fromQuery, toQuery, name)
	if err != nil {
		return
	}

	return a.executor.ExecSync(a.ctx, functionContext)
}

func PrepareSubscribeFunc(id string, subscribePayload SubscribePayload) (fc *pbtypes.FunctionContext, err error) {
	fc = &pbtypes.FunctionContext{
		Id: id,
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.SubscribeFunctionPath,
		},
	}

	fc.Value, err = json.Marshal(subscribePayload)
	return
}
