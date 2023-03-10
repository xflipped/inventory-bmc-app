// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"

	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
)

// TODO:
// Possible update:
// 1. execute on system object
// 2. get redfish device object using cmdb finder.Links()
// 3. get redfish client using api, login and password from redfish device object
// 4. get system using client to fill etag value for PATCH requests
// 5. reset system with reset type
// 6. rerun inventory

// resetFunction executes on '[device-uuid].redfish-devices.root'
func (a *Agent) resetFunction(ctx module.Context) (err error) {
	var redfishDevice device.RedfishDevice
	if err = json.Unmarshal(ctx.CmdbContext(), &redfishDevice); err != nil {
		return
	}
	var resetType redfish.ResetType
	if err = json.Unmarshal(ctx.Message(), &resetType); err != nil {
		return
	}

	client, err := utils.ConnectToRedfish(ctx, redfishDevice)
	if err != nil {
		return err
	}
	defer client.Logout()

	if err = a.resetSystem(client.Service, resetType); err != nil {
		return
	}

	// rerun inventory
	return executeInventory(ctx)
}

func (a *Agent) resetSystem(service *gofish.Service, resetType redfish.ResetType) (err error) {
	systems, err := service.Systems()
	if err != nil {
		return err
	}

	for _, system := range systems {
		if err = system.Reset(resetType); err != nil {
			return
		}
	}
	return nil
}

func executeInventory(ctx module.Context) (err error) {
	functionContext, err := agent.PrepareInventoryFunc(ctx.Self().Id)
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
