// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"
	"strings"

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
	var functionContext *pbtypes.FunctionContext
	managerLink := fmt.Sprintf("mng-%s", manager.UUID)

	document, err := a.getDocument("mng-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
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

	return a.createOrUpdateManagerCommandShell(ctx, redfishDevice, manager)
}

func (a *Agent) createOrUpdateManagerCommandShell(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager) (err error) {
	cmdShell := &bootstrap.RedfishCommandShell{CommandShell: &manager.CommandShell}

	parentNode, err := a.getDocument("mng-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("command-shell.mng-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-command-shell", "command-shell", cmdShell)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), cmdShell)
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
	parentNode, err := a.getDocument("mng-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	ifaces, err := manager.EthernetInterfaces()
	if err != nil {
		return
	}

	for _, iface := range ifaces {
		if err = a.createOrUpdateEthernetInterface(ctx, redfishDevice, parentNode, manager, iface); err != nil {
			return
		}
	}
	return
}

func (a *Agent) createOrUpdateEthernetInterface(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, manager *redfish.Manager, rfIface *redfish.EthernetInterface) (err error) {
	ifaceLink := strings.ToLower(rfIface.ID)
	iface := &bootstrap.RedfishEthernetInterface{EthernetInterface: rfIface}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.mng-%s.service.%s.redfish-devices.root", ifaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-ethernet-interface", ifaceLink, iface)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), iface)
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
	parentNode, err := a.getDocument("mng-%s.service.%s.redfish-devices.root", manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	ifaces, err := manager.HostInterfaces()
	if err != nil {
		return
	}

	for _, iface := range ifaces {
		if err = a.createOrUpdateHostInterface(ctx, redfishDevice, parentNode, manager, iface); err != nil {
			return
		}
	}

	return
}

func (a *Agent) createOrUpdateHostInterface(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, manager *redfish.Manager, rfIface *redfish.HostInterface) (err error) {
	ifaceLink := fmt.Sprintf("host-ifs-%s", strings.ToLower(rfIface.ID))
	iface := &bootstrap.RedfishHostInterface{HostInterface: rfIface}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.mng-%s.service.%s.redfish-devices.root", ifaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-host-interface", ifaceLink, iface)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), iface)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateHostInterfaceStatus(ctx, redfishDevice, manager, ifaceLink, rfIface)
}

func (a *Agent) createOrUpdateHostInterfaceStatus(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager, ifaceLink string, iface *redfish.HostInterface) (err error) {
	status := &bootstrap.RedfishStatus{Status: iface.Status}
	parentNode, err := a.getDocument("%s.mng-%s.service.%s.redfish-devices.root", ifaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.mng-%s.service.%s.redfish-devices.root", ifaceLink, manager.UUID, redfishDevice.UUID())
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

	return a.createOrUpdateHostInterfaceType(ctx, redfishDevice, manager, ifaceLink, iface)
}

func (a *Agent) createOrUpdateHostInterfaceType(ctx module.Context, redfishDevice device.RedfishDevice, manager *redfish.Manager, ifaceLink string, iface *redfish.HostInterface) (err error) {
	ifaceType := &bootstrap.RedfishHostInterfaceType{HostInterfaceType: &iface.HostInterfaceType}
	parentNode, err := a.getDocument("%s.mng-%s.service.%s.redfish-devices.root", ifaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("type.%s.mng-%s.service.%s.redfish-devices.root", ifaceLink, manager.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-host-interface-type", "type", ifaceType)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), ifaceType)
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
