// Copyright 2023 NJWS Inc.

package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
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

	log.Infof("Query: %s", query)

	nodes, err := qdsl.Qdsl(ctx, query, qdsl.WithId(), qdsl.WithObject(), qdsl.WithType(), qdsl.WithKey())
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

		if err = createOrUpdateLink(ctx, executor, node); err != nil {
			return
		}
	}

	return
}

func updateCredentials(ctx context.Context, executor executor.Executor, node *documents.Node, login, password string) (err error) {
	var redfishDevice device.RedfishDevice
	if err = json.Unmarshal(node.Object, &redfishDevice); err != nil {
		return
	}

	redfishDevice.Login = login
	redfishDevice.Password = password

	log.Infof("update uuid: %s cmdb id: %s", redfishDevice.UUID(), node.Id)

	// pass/login from: update, not replace
	functionContext, err := system.UpdateObject(node.Id.String(), redfishDevice)
	if err != nil {
		return
	}

	return executor.ExecSync(ctx, functionContext)
}

func createOrUpdateLink(ctx context.Context, executor executor.Executor, node *documents.Node) (err error) {
	var functionContext *pbtypes.FunctionContext

	route := &pbtypes.FunctionRoute{
		Url:             "http://inventory-bmc:31001/statefun",
		ExecuteOnCreate: true,
		ExecuteOnUpdate: true,
	}

	query := fmt.Sprintf("%s.%s", node.Key, types.FunctionPath)

	if linkDocument, err := getDocument(ctx, query); err == nil {
		functionContext, err = system.UpdateAdvancedLink(linkDocument.LinkId.String(), route)
		if err != nil {
			return err
		}
	} else {
		log.Error(err)
		function, err := getDocument(ctx, types.FunctionPath)
		if err != nil {
			return err
		}

		functionContext, err = system.CreateLink(function.Id.String(), node.Id.String(), node.Key, function.Type, route)
		if err != nil {
			return err
		}
	}

	return executor.ExecSync(ctx, functionContext)
}
