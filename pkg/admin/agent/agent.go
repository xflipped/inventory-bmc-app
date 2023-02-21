// Copyright 2023 NJWS Inc.

package agent

import (
	"context"
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func getDocument(ctx context.Context, query string) (node *documents.Node, err error) {
	nodes, err := qdsl.Qdsl(ctx, query, qdsl.WithKey(), qdsl.WithId(), qdsl.WithType(), qdsl.WithLinkId())
	if err != nil {
		return
	}
	for _, node = range nodes {
		return
	}
	err = fmt.Errorf("document '%s' not found", query)
	return
}

func ChangeCredentials(ctx context.Context, query, login, password string) (err error) {
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

		if err = updateCredentials(ctx, executor, node, login, password); err != nil {
			return
		}

		if err = executeInventory(ctx, executor, node); err != nil {
			return
		}
	}

	return
}

func updateCredentials(ctx context.Context, executor executor.Executor, node *documents.Node, login, password string) (err error) {
	redfishDevice := device.RedfishDevice{
		Login:    login,
		Password: password,
	}

	log.Infof("update document: %s", node.Id)

	// pass/login from: update, not replace
	functionContext, err := system.UpdateObject(node.Id.String(), redfishDevice)
	if err != nil {
		return
	}

	return executor.ExecSync(ctx, functionContext)
}

func executeInventory(ctx context.Context, executor executor.Executor, node *documents.Node) (err error) {
	functionContext, err := agent.PrepareInventoryFunc(node.Id.String())
	if err != nil {
		return
	}
	return executor.ExecAsync(ctx, functionContext)
}
