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
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

const (
	systemMask                            = "system-%s.service.*[?@._id == '%s'?].objects.root"
	biosMask                              = "bios.system-%s.service.*[?@._id == '%s'?].objects.root"
	ledMask                               = "led.system-%s.service.*[?@._id == '%s'?].objects.root"
	statusMask                            = "status.system-%s.service.*[?@._id == '%s'?].objects.root"
	bootMask                              = "boot.system-%s.service.*[?@._id == '%s'?].objects.root"
	bootOptionMask                        = "%s.boot.system-%s.service.*[?@._id == '%s'?].objects.root"
	secureBootMask                        = "secure-boot.system-%s.service.*[?@._id == '%s'?].objects.root"
	subSystemMask                         = "%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	statusSubSystemMask                   = "status.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	pcieInterfaceMask                     = "pcie-interface.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	powerStateMask                        = "power-state.system-%s.service.*[?@._id == '%s'?].objects.root"
	powerRestoreStateMask                 = "power-restore-policy.system-%s.service.*[?@._id == '%s'?].objects.root"
	processorSummaryMask                  = "processor-summary.system-%s.service.*[?@._id == '%s'?].objects.root"
	statusProcessorSummaryMask            = "status.processor-summary.system-%s.service.*[?@._id == '%s'?].objects.root"
	memorySummaryMask                     = "memory-summary.system-%s.service.*[?@._id == '%s'?].objects.root"
	statusMemorySummaryMask               = "status.memory-summary.system-%s.service.*[?@._id == '%s'?].objects.root"
	hostWatchdogTimerMask                 = "host-watchdog-timer.system-%s.service.*[?@._id == '%s'?].objects.root"
	subSubSystemMask                      = "%s.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	statusSubSubSystemMask                = "status.%s.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	ledSubSubSystemMask                   = "led.%s.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	locationSubSubSystemMask              = "location.%s.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	partLocationLocationSubSubSystemMask  = "part-location.location.%s.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	placementLocationSubSubSystemMask     = "placement.location.%s.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
	postalAddressLocationSubSubSystemMask = "postal-address.location.%s.%s.system-%s.service.*[?@._id == '%s'?].objects.root"
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
	p.Exec(func() error { return a.createOrUpdatePCIeFunctions(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemPowerState(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemPowerRestorePolicy(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateProcessorSummary(ctx, document, computerSystem) })
	// p.Exec(func() error { return a.createOrUpdateProcessors(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateMemorySummary(ctx, document, computerSystem) })
	// p.Exec(func() error { return a.createOrUpdateMemories(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateMemoryDomains(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateHostWatchdogTimer(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSimpleStorages(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateStorages(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemNetworkInterfaces(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemEthernetInterfaces(ctx, document, computerSystem) })
	p.Exec(func() error { return a.createOrUpdateSystemLogServices(ctx, document, computerSystem) })
	// TODO: add new entities if available etc.
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

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPcieDeviceID, deviceLink, device, subSystemMask, deviceLink, computerSystem.UUID, ctx.Self().Id)
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
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSystemMask, deviceLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePCIeDeviceInterface(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, deviceLink string, device *redfish.PCIeDevice) (err error) {
	deviceInterface := &device.PCIeInterface
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPcieInterfaceID, types.RedfishPcieInterfaceLink, deviceInterface, pcieInterfaceMask, deviceLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePCIeFunctions(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	functions, err := computerSystem.PCIeFunctions()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, function := range functions {
		function := function
		p.Exec(func() error { return a.createOrUpdatePCIeFunction(ctx, parentNode, computerSystem, function) })
	}

	return p.Wait()
}

func (a *Agent) createOrUpdatePCIeFunction(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, function *redfish.PCIeFunction) (err error) {
	functionLink := fmt.Sprintf("pcie-function-%s", function.ID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPcieFunctionID, functionLink, function, subSystemMask, functionLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdatePCIeFunctionStatus(ctx, document, computerSystem, functionLink, function)
}

func (a *Agent) createOrUpdatePCIeFunctionStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, functionLink string, function *redfish.PCIeFunction) (err error) {
	status := &bootstrap.RedfishStatus{Status: function.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSystemMask, functionLink, computerSystem.UUID, ctx.Self().Id)
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

func (a *Agent) createOrUpdateProcessors(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	processors, err := computerSystem.Processors()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, processor := range processors {
		processor := processor
		p.Exec(func() error { return a.createOrUpdateProcessor(ctx, parentNode, computerSystem, processor) })
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateProcessor(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, processor *redfish.Processor) (err error) {
	processorLink := processor.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishProcessorID, processorLink, processor, subSystemMask, processorLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateProcessorStatus(ctx, document, computerSystem, processorLink, processor)
}

func (a *Agent) createOrUpdateProcessorStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, processorLink string, processor *redfish.Processor) (err error) {
	status := &bootstrap.RedfishStatus{Status: processor.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSystemMask, processorLink, computerSystem.UUID, ctx.Self().Id)
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
	memoryLink := memory.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishMemoryID, memoryLink, memory, subSystemMask, memoryLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateMemoryStatus(ctx, document, computerSystem, memoryLink, memory)
}

func (a *Agent) createOrUpdateMemoryStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, memoryLink string, memory *redfish.Memory) (err error) {
	status := &bootstrap.RedfishStatus{Status: memory.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSystemMask, memoryLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateMemoryDomains(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	memoryDomains, err := computerSystem.MemoryDomains()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, memoryDomain := range memoryDomains {
		memoryDomain := memoryDomain
		p.Exec(func() error { return a.createOrUpdateMemoryDomain(ctx, parentNode, computerSystem, memoryDomain) })
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateMemoryDomain(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, memoryDomain *redfish.MemoryDomain) (err error) {
	memoryDomainLink := memoryDomain.ID
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishMemoryDomainID, memoryDomainLink, memoryDomain, subSystemMask, memoryDomainLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateHostWatchdogTimer(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	hostWatchdogTimer := computerSystem.HostWatchdogTimer
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishHostWatchdogTimerID, types.RedfishHostWatchdogTimerLink, hostWatchdogTimer, hostWatchdogTimerMask, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSimpleStorages(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	simpleStorages, err := computerSystem.SimpleStorages()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, simpleStorage := range simpleStorages {
		simpleStorage := simpleStorage
		p.Exec(func() error { return a.createOrUpdateSimpleStorage(ctx, parentNode, computerSystem, simpleStorage) })
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateSimpleStorage(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, simpleStorage *redfish.SimpleStorage) (err error) {
	simpleStorageLink := fmt.Sprintf("simple-storage-%s", simpleStorage.ID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishSimpleStorageID, simpleStorageLink, simpleStorage, subSystemMask, simpleStorageLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdateSimpleStorageStatus(ctx, document, computerSystem, simpleStorageLink, simpleStorage)
	})
	for _, device := range simpleStorage.Devices {
		device := device
		p.Exec(func() error {
			return a.createOrUpdateSimpleStorageDevice(ctx, document, computerSystem, simpleStorageLink, device)
		})
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateSimpleStorageStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, simpleStorageLink string, simpleStorage *redfish.SimpleStorage) (err error) {
	status := &bootstrap.RedfishStatus{Status: simpleStorage.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSystemMask, simpleStorageLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSimpleStorageDevice(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, simpleStorageLink string, device redfish.Device) (err error) {
	deviceLink := device.Name

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStorageDeviceID, deviceLink, device, subSubSystemMask, deviceLink, simpleStorageLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateSimpleStorageDeviceStatus(ctx, document, computerSystem, simpleStorageLink, deviceLink, device)
}

func (a *Agent) createOrUpdateSimpleStorageDeviceStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, simpleStorageLink, deviceLink string, device redfish.Device) (err error) {
	status := &bootstrap.RedfishStatus{Status: device.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSubSystemMask, deviceLink, simpleStorageLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateStorages(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	storages, err := computerSystem.Storage()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, storage := range storages {
		storage := storage
		p.Exec(func() error { return a.createOrUpdateStorage(ctx, parentNode, computerSystem, storage) })
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateStorage(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storage *redfish.Storage) (err error) {
	storageLink := storage.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStorageID, storageLink, storage, subSystemMask, storageLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdateStorageStatus(ctx, document, computerSystem, storageLink, storage)
	})
	p.Exec(func() error {
		return a.createOrUpdateDrives(ctx, document, computerSystem, storageLink, storage)
	})
	p.Exec(func() error {
		return a.createOrUpdateVolumes(ctx, document, computerSystem, storageLink, storage)
	})

	return p.Wait()
}

func (a *Agent) createOrUpdateStorageStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink string, storage *redfish.Storage) (err error) {
	status := &bootstrap.RedfishStatus{Status: storage.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSystemMask, storageLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateDrives(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink string, storage *redfish.Storage) (err error) {
	drives, err := storage.Drives()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, drive := range drives {
		drive := drive
		p.Exec(func() error { return a.createOrUpdateDrive(ctx, parentNode, computerSystem, storageLink, drive) })
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateDrive(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink string, drive *redfish.Drive) (err error) {
	driveLink := drive.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishDriveID, driveLink, drive, subSubSystemMask, driveLink, storageLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdateDriveStatus(ctx, document, computerSystem, storageLink, driveLink, drive)
	})
	p.Exec(func() error {
		return a.createOrUpdateDriveIndicatorLED(ctx, document, computerSystem, storageLink, driveLink, drive)
	})
	p.Exec(func() error {
		return a.createOrUpdateDriveLocation(ctx, document, computerSystem, storageLink, driveLink, drive)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdateDriveStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink, driveLink string, drive *redfish.Drive) (err error) {
	status := &bootstrap.RedfishStatus{Status: drive.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSubSystemMask, driveLink, storageLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateDriveIndicatorLED(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink, driveLink string, drive *redfish.Drive) (err error) {
	led := &bootstrap.RedfishLed{Led: drive.IndicatorLED}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLedID, types.RedfishLedLink, led, ledSubSubSystemMask, driveLink, storageLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateDriveLocation(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink, driveLink string, drive *redfish.Drive) (err error) {
	location := drive.PhysicalLocation

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLocationID, types.RedfishLocationLink, location, locationSubSubSystemMask, driveLink, storageLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdateDrivePartLocation(ctx, document, computerSystem, storageLink, driveLink, location.PartLocation)
	})
	p.Exec(func() error {
		return a.createOrUpdateDrivePlacement(ctx, document, computerSystem, storageLink, driveLink, location.Placement)
	})
	p.Exec(func() error {
		return a.createOrUpdateDrivePostalAddress(ctx, document, computerSystem, storageLink, driveLink, location.PostalAddress)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdateDrivePartLocation(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink, driveLink string, partLocation common.PartLocation) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPartLocationID, types.RedfishPartLocationLink, partLocation, partLocationLocationSubSubSystemMask, driveLink, storageLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateDrivePlacement(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink, driveLink string, placement common.Placement) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPlacementID, types.RedfishPlacementLink, placement, placementLocationSubSubSystemMask, driveLink, storageLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateDrivePostalAddress(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink, driveLink string, postalAddress common.PostalAddress) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPostalAddressID, types.RedfishPostalAddressLink, postalAddress, postalAddressLocationSubSubSystemMask, driveLink, storageLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateVolumes(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink string, storage *redfish.Storage) (err error) {
	volumes, err := storage.Volumes()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, volume := range volumes {
		volume := volume
		p.Exec(func() error { return a.createOrUpdateVolume(ctx, parentNode, computerSystem, storageLink, volume) })
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateVolume(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink string, volume *redfish.Volume) (err error) {
	volumeLink := volume.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishVolumeID, volumeLink, volume, subSubSystemMask, volumeLink, storageLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateVolumeStatus(ctx, document, computerSystem, storageLink, volumeLink, volume)
}

func (a *Agent) createOrUpdateVolumeStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, storageLink, volumeLink string, volume *redfish.Volume) (err error) {
	status := &bootstrap.RedfishStatus{Status: volume.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSubSystemMask, volumeLink, storageLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSystemNetworkInterfaces(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	networkInterfaces, err := computerSystem.NetworkInterfaces()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, networkInterface := range networkInterfaces {
		networkInterface := networkInterface
		p.Exec(func() error {
			return a.createOrUpdateSystemNetworkInterface(ctx, parentNode, computerSystem, networkInterface)
		})
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateSystemNetworkInterface(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, networkInterface *redfish.NetworkInterface) (err error) {
	networkInterfaceLink := networkInterface.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishNetworkInterfaceID, networkInterfaceLink, networkInterface, subSystemMask, networkInterfaceLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateSystemNetworkInterfaceStatus(ctx, document, computerSystem, networkInterfaceLink, networkInterface)
}

func (a *Agent) createOrUpdateSystemNetworkInterfaceStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, networkInterfaceLink string, networkInterface *redfish.NetworkInterface) (err error) {
	status := &bootstrap.RedfishStatus{Status: networkInterface.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSystemMask, networkInterfaceLink, computerSystem.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateSystemEthernetInterfaces(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	ethernetInterfaces, err := computerSystem.EthernetInterfaces()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, ethernetInterface := range ethernetInterfaces {
		ethernetInterface := ethernetInterface
		p.Exec(func() error {
			return a.createOrUpdateSystemEthernetInterface(ctx, parentNode, computerSystem, ethernetInterface)
		})
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateSystemEthernetInterface(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, ethernetInterface *redfish.EthernetInterface) (err error) {
	ethernetInterfaceLink := ethernetInterface.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishEthernetInterfaceID, ethernetInterfaceLink, ethernetInterface, subSystemMask, ethernetInterfaceLink, computerSystem.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateSystemEthernetInterfaceStatus(ctx, document, computerSystem, ethernetInterfaceLink, ethernetInterface)
}

func (a *Agent) createOrUpdateSystemEthernetInterfaceStatus(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, ethernetInterfaceLink string, ethernetInterface *redfish.EthernetInterface) (err error) {
	status := &bootstrap.RedfishStatus{Status: ethernetInterface.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSystemMask, ethernetInterfaceLink, computerSystem.UUID, ctx.Self().Id)
}
