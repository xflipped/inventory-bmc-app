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

func getDocument(ctx context.Context, query string) (document *documents.Node, err error) {
	documents, err := qdsl.Qdsl(ctx, query, qdsl.WithKey(), qdsl.WithId(), qdsl.WithType())
	if err != nil {
		return
	}
	for _, document = range documents {
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

	documents, err := qdsl.Qdsl(ctx, query, qdsl.WithId(), qdsl.WithObject(), qdsl.WithType(), qdsl.WithLink())
	if err != nil {
		return
	}

	for _, document := range documents {
		if document.Type != types.RedfishDeviceKey {
			continue
		}

		var redfishDevice device.RedfishDevice
		if err = json.Unmarshal(document.Object, &redfishDevice); err != nil {
			return
		}

		redfishDevice.Login = login
		redfishDevice.Password = password

		log.Infof("update uuid: %s cmdb id: %s", redfishDevice.UUID(), document.Id)

		// pass/login from: update, not replace
		functionContext, err := system.UpdateObject(document.Id.String(), redfishDevice)
		if err != nil {
			return err
		}

		if err = executor.ExecSync(ctx, functionContext); err != nil {
			return err
		}

		functionContext, err = createLink(ctx, document.Id)
		if err != nil {
			return err
		}

		if err = executor.ExecSync(ctx, functionContext); err != nil {
			return err
		}

	}

	return
}

func createLink(ctx context.Context, id documents.DocumentID) (functionContext *pbtypes.FunctionContext, err error) {
	route := &pbtypes.FunctionRoute{
		Url:             "http://inventory-bmc:31001/statefun",
		ExecuteOnCreate: true,
		ExecuteOnUpdate: true,
	}

	function, err := getDocument(ctx, types.FunctionPath)
	if err != nil {
		return
	}

	return system.CreateLink(function.Id.String(), id.String(), id.Key(), function.Type, route)
}
