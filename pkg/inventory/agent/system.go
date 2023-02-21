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
	systemMask                 = "system-%s.service.*[?@._id == '%s'?].objects.root"
	biosMask                   = "bios.system-%s.service.*[?@._id == '%s'?].objects.root"
	ledMask                    = "led.system-%s.service.*[?@._id == '%s'?].objects.root"
	statusMask                 = "status.system-%s.service.*[?@._id == '%s'?].objects.root"
	bootMask                   = "boot.system-%s.service.*[?@._id == '%s'?].objects.root"
	bootOptionMask             = "%s.boot.system-%s.service.*[?@._id == '%s'?].objects.root"
	secureBootMask             = "secure-boot.system-%s.service.*[?@._id == '%s'?].objects.root"
	pcieMask                   = "%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	pcieStatusMask             = "status.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	pcieInterfaceMask          = "pcie-interface.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	powerStateMask             = "power-state.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	powerRestoreStateMask      = "power-restore-policy.system-%s.service.*[?@._id == '%s'?].objects.root"
	processorSummaryMask       = "processor-summary.system-%s.service.*[?@._id == '%s'?].objects.root"
	statusProcessorSummaryMask = "status.processor-summary.system-%s.service.*[?@._id == '%s'?].objects.root"
	memorySummaryMask          = "memory-summary.system-%s.service.*[?@._id == '%s'?].objects.root"
	statusMemorySummaryMask    = "status.memory-summary.system-%s.service.*[?@._id == '%s'?].objects.root"
	hostWatchdogTimerMask      = "host-watchdog-timer.system-%s.service.*[?@._id == '%s'?].objects.root"
)

func (a *Agent) createOrUpdateSystems(ctx module.Context, service *gofish.Service, parentNode *documents.Node) (err error) {
	systems, err := service.Systems()
	if err != nil {
		return
	}
	p := utils.NewParallel()
	for _, computerSystem := range systems {
		computerSystem := computerSystem
		p.Exec(func() error { return a.createOrUpdateSystem(ctx, parentNode, computerSystem) })
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateSystem(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	systemLink := fmt.Sprintf("system-%s", computerSystem.UUID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishSystemID, systemLink, computerSystem, systemMask, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return a.createOrUpdateSystemBIOS(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemLed(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemStatus(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemBoot(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemSecureBoot(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdatePCIeDevices(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemPowerState(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemPowerRestorePolicy(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateProcessorSummary(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateMemorySummary(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateHostWatchdogTimer(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateMemories(ctx, document, computerSystem) })
	// TODO: return createMemories, createProcessors, etc.
	return p.Wait()
}

func (a *Agent) createOrUpdateSystemBIOS(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	bios, err := computerSystem.Bios()
	if err != nil {
		return
	}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishBiosID, types.RedfishBiosLink, bios, biosMask, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSystemLed(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	led := &bootstrap.RedfishLed{Led: computerSystem.IndicatorLED}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLedID, types.RedfishLedLink, led, ledMask, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSystemStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	status := &bootstrap.RedfishStatus{Status: computerSystem.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusMask, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSystemBoot(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	boot := computerSystem.Boot

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishBootID, types.RedfishBootLink, boot, bootMask, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	bootOptions, err := computerSystem.BootOptions()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, bootOption := range bootOptions {
		bootOption := bootOption
		p.Exec(func() error { return a.createOrUpdateBootOption(ctx, document, computerSystem, bootOption) })
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateBootOption(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, bootOption *redfish.BootOption) (err error) {
	bootOptionLink := fmt.Sprintf("option-%s", bootOption.ID)
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishBootOptionID, bootOptionLink, bootOption, bootOptionMask, bootOptionLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSystemSecureBoot(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	secureBoot, err := computerSystem.SecureBoot()
	if err != nil {
		return
	}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishSecureBootOptionID, types.RedfishSecureBootLink, secureBoot, secureBootMask, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePCIeDevices(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	devices, err := computerSystem.PCIeDevices()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, device := range devices {
		device := device
		p.Exec(func() error { return a.createOrUpdatePCIeDevice(ctx, parentNode, computerSystem, device) })
	}

	return p.Wait()

}

func (a *Agent) createOrUpdatePCIeDevice(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, device *redfish.PCIeDevice) (err error) {
	deviceLink := fmt.Sprintf("pcie-device-%s", device.ID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPcieID, deviceLink, device, pcieMask, deviceLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdatePCIeDeviceStatus(ctx, document, computerSystem, deviceLink, device)
	})
	p.Exec(func() error {
		return a.createOrUpdatePCIeDeviceInterface(ctx, document, computerSystem, deviceLink, device)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdatePCIeDeviceStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, deviceLink string, device *redfish.PCIeDevice) (err error) {
	status := &bootstrap.RedfishStatus{Status: device.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, pcieStatusMask, deviceLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePCIeDeviceInterface(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, deviceLink string, device *redfish.PCIeDevice) (err error) {
	deviceInterface := &device.PCIeInterface
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPcieInterfaceID, types.RedfishPcieInterfaceLink, deviceInterface, pcieInterfaceMask, deviceLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSystemPowerState(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	powerState := &bootstrap.RedfishPowerState{PowerState: computerSystem.PowerState}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerStateID, types.RedfishPowerStateLink, powerState, powerStateMask, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSystemPowerRestorePolicy(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	powerRestorePolicy := &bootstrap.RedfishPowerRestorePolicy{PowerRestorePolicy: computerSystem.PowerRestorePolicy}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerRestorePolicyID, types.RedfishPowerRestorePolicyLink, powerRestorePolicy, powerRestoreStateMask, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateProcessorSummary(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	processorSummary := computerSystem.ProcessorSummary
	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishProcessorSummaryID, types.RedfishProcessorSummaryLink, processorSummary, processorSummaryMask, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}
	return a.createOrUpdateProcessorSummaryStatus(ctx, document, computerSystem, &processorSummary)
}

func (a *Agent) createOrUpdateProcessorSummaryStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, processorSummary *redfish.ProcessorSummary) (err error) {
	status := &bootstrap.RedfishStatus{Status: processorSummary.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusProcessorSummaryMask, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateMemorySummary(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	memorySummary := computerSystem.MemorySummary

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishMemorySummaryID, types.RedfishMemorySummaryLink, memorySummary, memorySummaryMask, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateMemorySummaryStatus(ctx, document, computerSystem, &memorySummary)
}

func (a *Agent) createOrUpdateMemorySummaryStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, memorySummary *redfish.MemorySummary) (err error) {
	status := &bootstrap.RedfishStatus{Status: memorySummary.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusMemorySummaryMask, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateHostWatchdogTimer(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	hostWatchdogTimer := computerSystem.HostWatchdogTimer
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishHostWatchdogTimerID, types.RedfishHostWatchdogTimerLink, hostWatchdogTimer, hostWatchdogTimerMask, computerSystem.UUID, ctx.Self().Id)
}

// TODO: currently unavailable structures, pending to implement
func (a *Agent) createOrUpdateMemories(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	memories, err := computerSystem.Memory()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, memory := range memories {
		memory := memory
		p.Exec(func() error { return a.createOrUpdateMemory(ctx, parentNode, computerSystem, memory) })
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateMemory(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, memory *redfish.Memory) (err error) {
	fmt.Println(memory.DeviceLocator, memory.MemoryLocation)
	return
}
