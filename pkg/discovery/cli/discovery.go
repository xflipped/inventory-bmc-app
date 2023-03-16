// Copyright 2023 NJWS Inc.

package cli

import (
	"context"
	"net/url"

	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
)

func Discovery(ctx context.Context, executor executor.Executor, addr string) (err error) {
	u, err := url.Parse(addr)
	if err != nil {
		return
	}
	var description device.Description

	redfishDevice := description.ToDevice(u)

	redfishDevicesObject, err := utils.GetDocument(ctx, types.RedfishDevicesPath)
	if err != nil {
		return
	}

	functionContext, err := agent.PrepareDiscoveryFunc(redfishDevicesObject.Id.String(), redfishDevice)
	if err != nil {
		return
	}

	return executor.ExecSync(ctx, functionContext)
}
