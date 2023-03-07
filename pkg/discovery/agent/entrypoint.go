// Copyright 2023 NJWS Inc.

package agent

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/koron/go-ssdp"
	"github.com/stmcginnis/gofish"
)

func (a *Agent) entrypoint() (err error) {
	// wait router, need to register port
	time.Sleep(time.Millisecond * 50)

	return a.createOrUpdateFunctionLink(types.FunctionContainerPath, types.DiscoveryFunctionPath, types.DiscoveryFunctionLink)
}

func (a *Agent) createOrUpdateAliveMessage(m *ssdp.AliveMessage) (err error) {
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

	// FIXME fix if http
	u.Scheme = "https"

	redfishDevice := description.ToDevice(u)

	functionContext, err := PrepareDiscoveryFunc(redfishDevicesObject.Id.String(), redfishDevice)
	if err != nil {
		return
	}

	return a.executor.ExecAsync(a.ctx, functionContext)
}

func (a *Agent) createOrUpdate(redfishDevice device.RedfishDevice, parentID string) (err error) {
	u, err := url.Parse(redfishDevice.Api)
	if err != nil {
		return
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(fmt.Sprintf("%s://%s/redfish/v1", u.Scheme, u.Host))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var service gofish.Service

	if err = json.NewDecoder(resp.Body).Decode(&service); err != nil {
		return
	}

	uuid := service.UUID

	var functionContext *pbtypes.FunctionContext
	redfishDeviceObject, err := a.getDocument("%s.redfish-devices.root", uuid)
	if err == nil {
		log.Infof("update uuid: %s cmdb id: %s", uuid, redfishDeviceObject.Id)
		// pass/login from: update, not replace
		functionContext, err = system.UpdateObject(redfishDeviceObject.Id.String(), redfishDevice)
		if err != nil {
			return
		}
	} else {
		log.Infof("create uuid: %s", uuid)
		functionContext, err = system.CreateChild(parentID, types.RedfishDeviceID, uuid, redfishDevice)
		if err != nil {
			return
		}
	}

	return a.executor.ExecSync(a.ctx, functionContext)
}
