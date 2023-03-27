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

func Subscribe(ctx context.Context, query, destinationUrl string) (err error) {
	executor, err := executor.New()
	if err != nil {
		return
	}
	defer executor.Close()

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

		if err = executeSubscribe(ctx, executor, node, destinationUrl); err != nil {
			return
		}
	}

	return
}

func executeSubscribe(ctx context.Context, executor executor.Executor, node *documents.Node, destinationUrl string) (err error) {
	functionContext, err := agent.PrepareSubscribeFunc(node.Id.String(), destinationUrl)
	if err != nil {
		return
	}
	return executor.ExecAsync(ctx, functionContext)
}
