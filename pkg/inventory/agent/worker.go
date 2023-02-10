// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/stmcginnis/gofish"
)

func (a *Agent) workerFunction(ctx module.Context) (err error) {
	var redfishDevice device.RedfishDevice
	if err = json.Unmarshal(ctx.CmdbContext(), &redfishDevice); err != nil {
		return
	}

	u, err := url.Parse(redfishDevice.Api)
	if err != nil {
		return
	}

	config := gofish.ClientConfig{
		Endpoint: fmt.Sprintf("%s://%s", u.Scheme, u.Host),
		Username: redfishDevice.Login,
		Password: redfishDevice.Password,
		Insecure: true,
	}

	client, err := gofish.ConnectContext(ctx, config)
	if err != nil {
		return
	}
	defer client.Logout()

	service := client.GetService()

	// create/update main (root) service
	if err = a.createOrUpdateService(ctx, redfishDevice, service); err != nil {
		return
	}

	return nil
}

func (a *Agent) createOrUpdateService(ctx module.Context, redfishDevice device.RedfishDevice, service *gofish.Service) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("service.%s.redfish-devices.root", redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(ctx.Self().Id, "types/redfish-service", "service", service)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), service)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	var wg sync.WaitGroup

	wg.Add(3)
	go a.createOrUpdateSystems(ctx, redfishDevice, service, &wg)
	go a.createOrUpdateChasseez(ctx, redfishDevice, service, &wg)
	go a.createOrUpdateManagers(ctx, redfishDevice, service, &wg)

	wg.Wait()

	return
}
