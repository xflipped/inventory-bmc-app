// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

type LedPayload struct {
	Led                  common.IndicatorLED
	ConnectionParameters utils.ConnectionParameters
}

// ledFunction executes on 'chassis-[chassis-uuid].service.[device-uuid].redfish-devices.root'
func (a *Agent) ledFunction(ctx module.Context) (err error) {
	chassis := &redfish.Chassis{}
	if err = json.Unmarshal(ctx.CmdbContext(), chassis); err != nil {
		return
	}
	var payload LedPayload
	if err = json.Unmarshal(ctx.Message(), &payload); err != nil {
		return
	}

	client, err := utils.Connect(ctx, payload.ConnectionParameters)
	if err != nil {
		return err
	}
	defer client.Logout()

	chassis, err = redfish.GetChassis(client, chassis.ODataID)
	if err != nil {
		return
	}
	chassis.IndicatorLED = payload.Led
	if err = chassis.Update(); err != nil {
		return
	}

	// TODO: rerun inventory on computer system?
	// rerun inventory on chassis object
	return a.asyncUpdateObject(ctx, chassis)
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
