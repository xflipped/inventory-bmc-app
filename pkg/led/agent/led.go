// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"

	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/common"
)

// TODO:
// Possible update:
// 1. execute on chassis object
// 2. get redfish device object using cmdb finder.Links()
// 3. get redfish client using api, login and password from redfish device object
// 4. get chassis using client to fill etag value for PATCH requests
// 5. update chassis with indicator led
// 6. run re-inventory

// ledFunction execute on id '[device-uuid].redfish-devices.root'
func (a *Agent) ledFunction(ctx module.Context) (err error) {
	var redfishDevice device.RedfishDevice
	if err = json.Unmarshal(ctx.CmdbContext(), &redfishDevice); err != nil {
		return
	}
	var indicatorLED common.IndicatorLED
	if err = json.Unmarshal(ctx.Message(), &indicatorLED); err != nil {
		return
	}

	client, err := utils.ConnectToRedfish(ctx, redfishDevice)
	if err != nil {
		return err
	}
	defer client.Logout()

	if err = a.updateChassisIndicatorLED(client.Service, indicatorLED); err != nil {
		return
	}

	// rerun inventory
	return executeInventory(ctx)
}

func (a *Agent) updateChassisIndicatorLED(service *gofish.Service, indicatorLED common.IndicatorLED) (err error) {
	chasseez, err := service.Chassis()
	if err != nil {
		return err
	}

	for _, chassee := range chasseez {
		chassee.IndicatorLED = indicatorLED
		if err = chassee.Update(); err != nil {
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
