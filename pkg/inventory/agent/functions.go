// Copyright 2023 NJWS Inc.

package agent

import (
	"context"
	"fmt"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
)

func createOrUpdateFunctionLink(ctx context.Context, fromQuery, toQuery, name string) (functionContext *pbtypes.FunctionContext, err error) {
	route := &pbtypes.FunctionRoute{
		Url: "http://inventory-bmc:31001/statefun",
	}

	query := fmt.Sprintf("%s.%s", name, fromQuery)

	if linkDocument, err := getDocument(ctx, query); err == nil {
		return system.UpdateAdvancedLink(linkDocument.LinkId.String(), route)
	}

	parent, err := getDocument(ctx, fromQuery)
	if err != nil {
		return
	}

	child, err := getDocument(ctx, toQuery)
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
