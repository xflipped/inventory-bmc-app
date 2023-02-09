// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/bootstrap"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
)

func (a *Agent) createOrUpdateManagers(ctx module.Context, redfishDevice device.RedfishDevice, service *gofish.Service) (err error) {
	parentNode, err := a.getDocument("service.%s.redfish-devices.root", redfishDevice.UUID())
	if err != nil {
		return
	}

	managers, err := service.Managers()
	if err != nil {
		return
	}

	for _, manager := range managers {
		if err = a.createOrUpdateManager(ctx, redfishDevice, parentNode, manager); err != nil {
			return
		}
	}

	return
}

func (a *Agent) createOrUpdateManager(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, manager *redfish.Manager) (err error) {
	managerLink := fmt.Sprintf("manager-%s", manager.UUID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("manager-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-manager", managerLink, manager)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), manager)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateManagerStatus(ctx, redfishDevice, manager)
}

func (a *Agent) createOrUpdateManagerStatus(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager) (err error) {
	status := &bootstrap.RedfishStatus{Status: manager.Status}

	parentNode, err := a.getDocument("manager-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.manager-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-status", "status", status)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), status)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateManagerPowerState(ctx, redfishDevice, manager)
}

func (a *Agent) createOrUpdateManagerPowerState(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager) (err error) {
	powerState := &bootstrap.RedfishPowerState{PowerState: manager.PowerState}

	parentNode, err := a.getDocument("manager-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("power-state.manager-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-power-state", "power-state", powerState)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), powerState)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateManagerCommandShell(ctx, redfishDevice, manager)
}

func (a *Agent) createOrUpdateManagerCommandShell(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager) (err error) {
	commandShell := manager.CommandShell

	parentNode, err := a.getDocument("manager-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("command-shell.manager-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-command-shell", "command-shell", commandShell)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), commandShell)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateEthernetInterfaces(ctx, redfishDevice, manager)
}

func (a *Agent) createOrUpdateEthernetInterfaces(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager) (err error) {
	parentNode, err := a.getDocument("manager-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	ethernetInterfaces, err := manager.EthernetInterfaces()
	if err != nil {
		return
	}

	for _, ethernetInterface := range ethernetInterfaces {
		if err = a.createOrUpdateEthernetInterface(ctx, redfishDevice, parentNode, manager, ethernetInterface); err != nil {
			return
		}
	}
	return
}

func (a *Agent) createOrUpdateEthernetInterface(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, manager *redfish.Manager, ethernetInterface *redfish.EthernetInterface) (err error) {
	ethernetInterfaceLink := ethernetInterface.ID

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.manager-%s.service.%s.redfish-devices.root", ethernetInterfaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-ethernet-interface", ethernetInterfaceLink, ethernetInterface)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), ethernetInterface)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateHostInterfaces(ctx, redfishDevice, manager)
}

func (a *Agent) createOrUpdateHostInterfaces(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager) (err error) {
	parentNode, err := a.getDocument("manager-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	hostInterfaces, err := manager.HostInterfaces()
	if err != nil {
		return
	}

	for _, hostInterface := range hostInterfaces {
		if err = a.createOrUpdateHostInterface(ctx, redfishDevice, parentNode, manager, hostInterface); err != nil {
			return
		}
	}

	return
}

func (a *Agent) createOrUpdateHostInterface(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, manager *redfish.Manager, hostInterface *redfish.HostInterface) (err error) {
	hostInterfaceLink := fmt.Sprintf("host-interface-%s", hostInterface.ID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.manager-%s.service.%s.redfish-devices.root", hostInterfaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-host-interface", hostInterfaceLink, hostInterface)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), hostInterface)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateHostInterfaceStatus(ctx, redfishDevice, manager, hostInterfaceLink, hostInterface)
}

func (a *Agent) createOrUpdateHostInterfaceStatus(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager, hostInterfaceLink string, hostInterface *redfish.HostInterface) (err error) {
	status := &bootstrap.RedfishStatus{Status: hostInterface.Status}

	parentNode, err := a.getDocument("%s.manager-%s.service.%s.redfish-devices.root", hostInterfaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.manager-%s.service.%s.redfish-devices.root", hostInterfaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-status", "status", status)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), status)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateHostInterfaceType(ctx, redfishDevice, manager, hostInterfaceLink, hostInterface)
}

func (a *Agent) createOrUpdateHostInterfaceType(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager, hostInterfaceLink string, hostInterface *redfish.HostInterface) (err error) {
	hostInterfaceType := &bootstrap.RedfishHostInterfaceType{HostInterfaceType: &hostInterface.HostInterfaceType}

	parentNode, err := a.getDocument("%s.manager-%s.service.%s.redfish-devices.root", hostInterfaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("type.%s.manager-%s.service.%s.redfish-devices.root", hostInterfaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-host-interface-type", "type", hostInterfaceType)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), hostInterfaceType)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return
}
