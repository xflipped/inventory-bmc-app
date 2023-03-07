// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"

	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
)

// discoveryFunction execute on id 'redfish-devices.root'
// with 'RedfishDevice' body
func (a *Agent) discoveryFunction(ctx module.Context) (err error) {
	var redfishDevice device.RedfishDevice
	if err = json.Unmarshal(ctx.Message(), &redfishDevice); err != nil {
		return
	}

	return a.createOrUpdate(redfishDevice, ctx.Self().Id)
}
