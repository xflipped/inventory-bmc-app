// Copyright 2023 NJWS Inc.

package agent

import (
	"net/url"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/koron/go-ssdp"
)

const (
	redfishFilter = "urn:dmtf-org:service:redfish-rest:1"
)

func (a *Agent) entrypoint() (err error) {
	m := &ssdp.Monitor{
		Alive: a.onAlive,
		Bye:   a.onBye,
	}

	if err = m.Start(); err != nil {
		return
	}

	<-a.ctx.Done()
	return a.ctx.Err()
}

func (a *Agent) onAlive(m *ssdp.AliveMessage) {
	if m.Type != redfishFilter {
		return
	}

	if err := a.createOrUpdate(m); err != nil {
		log.Error(err)
		return
	}
}

func (a *Agent) createOrUpdate(m *ssdp.AliveMessage) (err error) {
	redfishDevicesObject, err := a.getDocument(types.RedfishDevicesPath)
	if err != nil {
		return
	}

	description, err := device.GetDescription(m.Location)
	if err != nil {
		return
	}

	u, err := url.Parse(description.Device.PresentationURL)
	if err != nil {
		return
	}

	u.Scheme = "https"

	redfishDevice := description.ToDevice(u.String())

	var functionContext *pbtypes.FunctionContext
	redfishDeviceObject, err := a.getDocument("%s.redfish-devices.root", redfishDevice.UUID())
	if err == nil {
		log.Infof("update uuid: %s cmdb id: %s", redfishDevice.UUID(), redfishDeviceObject.Id)
		// pass/login from: update, not replace
		functionContext, err = system.UpdateObject(redfishDeviceObject.Id.String(), redfishDevice)
		if err != nil {
			return
		}
	} else {
		log.Infof("create uuid: %s", redfishDevice.UUID())
		functionContext, err = system.CreateChild(redfishDevicesObject.Id.String(), types.RedfishDeviceID, redfishDevice.UUID(), redfishDevice)
		if err != nil {
			return
		}
	}

	return a.executor.ExecAsync(a.ctx, functionContext)
}

func (a *Agent) onBye(m *ssdp.ByeMessage) {
	if m.Type != redfishFilter {
		return
	}

	log.Infof("Bye: From=%s Type=%s USN=%s", m.From.String(), m.Type, m.USN)
}
