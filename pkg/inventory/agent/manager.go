// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/bootstrap"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
)

const (
	managerMask             = "manager-%s.service.*[?@._id == '%s'?].objects.root"
	statusManagerMask       = "status.manager-%s.service.*[?@._id == '%s'?].objects.root"
	powerStateManagerMask   = "power-state.manager-%s.service.*[?@._id == '%s'?].objects.root"
	commandShellManagerMask = "command-shell.manager-%s.service.*[?@._id == '%s'?].objects.root"
	subManagerMask          = "%s.manager-%s.service.*[?@._id == '%s'?].objects.root"
	statusSubManagerMask    = "status.%s.manager-%s.service.*[?@._id == '%s'?].objects.root"
	typeSubManagerMask      = "type.%s.manager-%s.service.*[?@._id == '%s'?].objects.root"
)

func (a *Agent) createOrUpdateManagers(ctx module.Context, service *gofish.Service, parentNode *documents.Node) (err error) {
	managers, err := service.Managers()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, manager := range managers {
		manager := manager
		p.Exec(func() error { return a.createOrUpdateManager(ctx, parentNode, manager) })
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateManager(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager) (err error) {
	managerLink := fmt.Sprintf("manager-%s", manager.UUID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishManagerID, managerLink, manager, managerMask, manager.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return a.createOrUpdateManagerStatus(ctx, document, manager) })
	p.Exec(func() error { return a.createOrUpdateManagerPowerState(ctx, document, manager) })
	p.Exec(func() error { return a.createOrUpdateManagerCommandShell(ctx, document, manager) })
	p.Exec(func() error { return a.createOrUpdateEthernetInterfaces(ctx, document, manager) })
	p.Exec(func() error { return a.createOrUpdateHostInterfaces(ctx, document, manager) })
	return p.Wait()
}

func (a *Agent) createOrUpdateManagerStatus(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager) (err error) {
	status := &bootstrap.RedfishStatus{Status: manager.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusManagerMask, manager.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateManagerPowerState(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager) (err error) {
	powerState := &bootstrap.RedfishPowerState{PowerState: manager.PowerState}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerStateID, types.RedfishPowerStateLink, powerState, powerStateManagerMask, manager.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateManagerCommandShell(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager) (err error) {
	commandShell := manager.CommandShell
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishCommandShellID, types.RedfishCommandShellLink, commandShell, commandShellManagerMask, manager.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateEthernetInterfaces(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager) (err error) {
	ethernetInterfaces, err := manager.EthernetInterfaces()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, ethernetInterface := range ethernetInterfaces {
		ethernetInterface := ethernetInterface
		p.Exec(func() error { return a.createOrUpdateEthernetInterface(ctx, parentNode, manager, ethernetInterface) })
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateEthernetInterface(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager, ethernetInterface *redfish.EthernetInterface) (err error) {
	ethernetInterfaceLink := ethernetInterface.ID
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishEthernetInterfaceID, ethernetInterfaceLink, ethernetInterface, subManagerMask, ethernetInterfaceLink, manager.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateHostInterfaces(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager) (err error) {
	hostInterfaces, err := manager.HostInterfaces()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, hostInterface := range hostInterfaces {
		hostInterface := hostInterface
		p.Exec(func() error { return a.createOrUpdateHostInterface(ctx, parentNode, manager, hostInterface) })
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateHostInterface(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager, hostInterface *redfish.HostInterface) (err error) {
	hostInterfaceLink := fmt.Sprintf("host-interface-%s", hostInterface.ID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishHostInterfaceID, hostInterfaceLink, hostInterface, subManagerMask, hostInterfaceLink, manager.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	document, err = a.getDocument(subManagerMask, hostInterfaceLink, manager.UUID, ctx.Self().Id)
	if err != nil {
		return
	}
	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdateHostInterfaceStatus(ctx, document, manager, hostInterfaceLink, hostInterface)
	})
	p.Exec(func() error {
		return a.createOrUpdateHostInterfaceType(ctx, document, manager, hostInterfaceLink, hostInterface)
	})
	return p.Wait()

}

func (a *Agent) createOrUpdateHostInterfaceStatus(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager, hostInterfaceLink string, hostInterface *redfish.HostInterface) (err error) {
	status := &bootstrap.RedfishStatus{Status: hostInterface.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubManagerMask, hostInterfaceLink, manager.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateHostInterfaceType(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager, hostInterfaceLink string, hostInterface *redfish.HostInterface) (err error) {
	hostInterfaceType := &bootstrap.RedfishHostInterfaceType{HostInterfaceType: &hostInterface.HostInterfaceType}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishHostInterfaceTypeID, types.RedfishTypeLink, hostInterfaceType, typeSubManagerMask, hostInterfaceLink, manager.UUID, ctx.Self().Id)
}
