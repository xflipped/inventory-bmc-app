// Copyright 2023 NJWS Inc.

package cli

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/reset/agent"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func Reset(ctx context.Context, executor executor.Executor, query string, resetPayload agent.ResetPayload) (err error) {
	log.Infof("Query: %s", query)

	nodes, err := qdsl.Qdsl(ctx, query, qdsl.WithId(), qdsl.WithType())
	if err != nil {
		return
	}

	for _, node := range nodes {
		log.Infof("document: %s", node.Id)

		if node.Type != types.RedfishSystemKey {
			log.Infof("document: %s, skip...", node.Id)
			continue
		}

		if err = executeReset(ctx, executor, node, resetPayload); err != nil {
			return
		}
	}

	return
}

func executeReset(ctx context.Context, executor executor.Executor, node *documents.Node, resetPayload agent.ResetPayload) (err error) {
	functionContext, err := agent.PrepareResetFunc(node.Id.String(), resetPayload)
	if err != nil {
		return
	}
	return executor.ExecAsync(ctx, functionContext)
}
