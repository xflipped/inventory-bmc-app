// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"
	"fmt"
	"net/url"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish"
)

const serviceMask = "service.*[?@._id == '%s'?].objects.root"

func (a *Agent) inventoryFunction(ctx module.Context) (err error) {
	var redfishDevice device.RedfishDevice
	if err = json.Unmarshal(ctx.CmdbContext(), &redfishDevice); err != nil {
		return
	}

	u, err := url.Parse(redfishDevice.Api)
	if err != nil {
		return
	}

	// TODO: check options
	config := gofish.ClientConfig{
		Endpoint:  fmt.Sprintf("%s://%s", u.Scheme, u.Host),
		Username:  redfishDevice.Login,
		Password:  redfishDevice.Password,
		Insecure:  true,
		BasicAuth: true,
	}

	client, err := gofish.ConnectContext(ctx, config)
	if err != nil {
		return
	}
	defer client.Logout()
	// create/update main (root) service
	return a.createOrUpdateService(ctx, client.GetService())
}

func (a *Agent) createOrUpdateService(ctx module.Context, service *gofish.Service) (err error) {
	document, err := a.syncCreateOrUpdateChild(ctx, ctx.Self().Id, types.RedfishServiceID, types.RedfishServiceLink, service, serviceMask, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return a.createOrUpdateSystems(ctx, service, document) })
	p.Exec(func() error { return a.createOrUpdateChasseez(ctx, service, document) })
	p.Exec(func() error { return a.createOrUpdateManagers(ctx, service, document) })
	return p.Wait()
}

func (a *Agent) createSyncCreateOrUpdateChild(from, moType, name string, payload any, format string, args ...any) (functionContext *pbtypes.FunctionContext, err error) {
	document, err := a.getDocument(format, args...)
	if err != nil {
		return system.CreateChild(from, moType, name, payload)
	}

	return system.UpdateObject(document.Id.String(), payload)
}

func (a *Agent) syncCreateOrUpdateChild(ctx module.Context, from, moType, name string, payload any, format string, args ...any) (document *documents.Node, err error) {
	functionContext, err := a.createSyncCreateOrUpdateChild(from, moType, name, payload, format, args...)
	if err != nil {
		return
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.getDocument(format, args...)
}

func (a *Agent) asyncCreateOrUpdateChild(ctx module.Context, from, moType, name string, payload any, format string, args ...any) (err error) {
	functionContext, err := a.createSyncCreateOrUpdateChild(from, moType, name, payload, format, args...)
	if err != nil {
		return
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)
	return
}
