// Copyright 2023 NJWS Inc.

package cli

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/led/agent"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func Led(ctx context.Context, executor executor.Executor, query string, ledPayload agent.LedPayload) (err error) {
	log.Infof("Query: %s", query)

	nodes, err := qdsl.Qdsl(ctx, query, qdsl.WithId(), qdsl.WithType())
	if err != nil {
		return
	}

	for _, node := range nodes {
		log.Infof("document: %s", node.Id)

		if node.Type != types.RedfishChassisKey {
			log.Infof("document: %s, skip...", node.Id)
			continue
		}

		if err = executeLed(ctx, executor, node, ledPayload); err != nil {
			return
		}
	}

	return
}

func executeLed(ctx context.Context, executor executor.Executor, node *documents.Node, ledPayload agent.LedPayload) (err error) {
	functionContext, err := agent.PrepareLedFunc(node.Id.String(), ledPayload)
	if err != nil {
		return
	}
	return executor.ExecSync(ctx, functionContext)
}
