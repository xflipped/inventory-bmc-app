// Copyright 2023 NJWS Inc.

package cli

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func Inventory(ctx context.Context, executor executor.Executor, query string, inventoryPayload agent.InventoryPayload) (err error) {
	log.Infof("Query: %s", query)

	nodes, err := qdsl.Qdsl(ctx, query, qdsl.WithId(), qdsl.WithType())
	if err != nil {
		return
	}

	for _, node := range nodes {
		log.Infof("document: %s", node.Id)

		if node.Type != types.RedfishDeviceKey {
			log.Infof("document: %s, skip...", node.Id)
			continue
		}

		if err = executeInventory(ctx, executor, node, inventoryPayload); err != nil {
			return
		}
	}

	return
}

func executeInventory(ctx context.Context, executor executor.Executor, node *documents.Node, inventoryPayload agent.InventoryPayload) (err error) {
	functionContext, err := agent.PrepareInventoryFunc(node.Id.String(), inventoryPayload)
	if err != nil {
		return
	}
	return executor.ExecAsync(ctx, functionContext)
}
