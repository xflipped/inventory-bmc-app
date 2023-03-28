// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
)

type ResetPayload struct {
	ResetType            redfish.ResetType
	ConnectionParameters utils.ConnectionParameters
}

// resetFunction executes on 'system-[system-uuid].service.[device-uuid].redfish-devices.root'
func (a *Agent) resetFunction(ctx module.Context) (err error) {
	system := &redfish.ComputerSystem{}
	if err = json.Unmarshal(ctx.CmdbContext(), system); err != nil {
		return
	}
	var payload ResetPayload
	if err = json.Unmarshal(ctx.Message(), &payload); err != nil {
		return
	}

	client, err := utils.Connect(ctx, payload.ConnectionParameters)
	if err != nil {
		return err
	}
	defer client.Logout()

	// fetch system to update etag
	system, err = redfish.GetComputerSystem(client, system.ODataID)
	if err != nil {
		return
	}

	if err = system.Reset(payload.ResetType); err != nil {
		return
	}

	// fetch updated system
	system, err = redfish.GetComputerSystem(client, system.ODataID)
	if err != nil {
		return
	}

	// rerun inventory on computer system object
	return a.asyncUpdateObject(ctx, system)
}

func (a *Agent) asyncUpdateObject(ctx module.Context, payload any) (err error) {
	functionContext, err := system.UpdateObject(ctx.Self().Id, payload)
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
