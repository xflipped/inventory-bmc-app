// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
)

func (a *Agent) getDocument(query string) (document *documents.Node, err error) {
	documents, err := qdsl.Qdsl(a.ctx, query, qdsl.WithKey(), qdsl.WithId(), qdsl.WithType())
	if err != nil {
		return
	}
	for _, document = range documents {
		return
	}
	err = fmt.Errorf("document '%s' not found", query)
	return
}

func (a *Agent) getFunction() (document *documents.Node, err error) {
	// search function_type init 'init.exmt.functions.root'
	return a.getDocument(types.FunctionPath)
}

func (a *Agent) getNodes() (document *documents.Node, err error) {
	// search 'nodes.root'
	return a.getDocument(types.BmcContainerPath)
}
