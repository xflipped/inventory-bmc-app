// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
)

func (a *Agent) entrypoint() (err error) {
	redfishDevicesObject, err := a.getDocument(types.RedfishDevicesPath)
	if err != nil {
		return
	}

	services, err := redfish.GetServices()
	if err != nil {
		return
	}

	log.Infof("found (%d)", len(services))

	for _, service := range services {
		description, err := device.GetDescription(service.Location)
		if err != nil {
			log.Error(err)
			continue
		}

		redfishDevice := description.ToDevice(service.Header().Get("Al"))

		redfishDeviceObject, err := a.getDocument(fmt.Sprintf("%s.redfish-devices.root", redfishDevice.UUID()))
		if err == nil {
			log.Infof("update uuid: %s cmdb id: %s", redfishDevice.UUID(), redfishDeviceObject.Id)
			// pass/login from: update, not replace
			functionContext, err := system.UpdateObject(redfishDeviceObject.Id.String(), redfishDevice)
			if err != nil {
				return err
			}

			if err = a.executor.ExecAsync(a.ctx, functionContext); err != nil {
				return err
			}
			continue
		}

		log.Infof("create uuid: %s", redfishDevice.UUID())

		functionContext, err := system.CreateChild(redfishDevicesObject.Id.String(), types.RedfishDeviceID, redfishDevice.UUID(), redfishDevice)
		if err != nil {
			return err
		}

		if err = a.executor.ExecAsync(a.ctx, functionContext); err != nil {
			return err
		}
	}

	return
}
