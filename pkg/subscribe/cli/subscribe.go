// Copyright 2023 NJWS Inc.

package cli

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/subscribe/agent"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func Subscribe(ctx context.Context, executor executor.Executor, query string, subscribePayload agent.SubscribePayload) (err error) {
	log.Infof("Query: %s", query)

	nodes, err := qdsl.Qdsl(ctx, query, qdsl.WithId(), qdsl.WithType())
	if err != nil {
		return
	}

	for _, node := range nodes {
		log.Infof("document: %s", node.Id)

		if node.Type != types.RedfishEventServiceKey {
			log.Infof("document: %s, skip...", node.Id)
			continue
		}

		if err = executeSubscribe(ctx, executor, node, subscribePayload); err != nil {
			return
		}
	}

	return
}

func executeSubscribe(ctx context.Context, executor executor.Executor, node *documents.Node, subscribePayload agent.SubscribePayload) (err error) {
	functionContext, err := agent.PrepareSubscribeFunc(node.Id.String(), subscribePayload)
	if err != nil {
		return
	}
	return executor.ExecAsync(ctx, functionContext)
}
