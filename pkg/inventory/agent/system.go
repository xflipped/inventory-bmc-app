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

func (a *Agent) createOrUpdateSystems(ctx module.Context, redfishDevice device.RedfishDevice, service *gofish.Service) (err error) {
	parentNode, err := a.getDocument("service.%s.redfish-devices.root", redfishDevice.UUID())
	if err != nil {
		return
	}

	systems, err := service.Systems()
	if err != nil {
		return
	}

	for _, computerSystem := range systems {
		if err = a.createOrUpdateSystem(ctx, redfishDevice, parentNode, computerSystem); err != nil {
			return
		}
	}

	// create/update chassis
	// TODO: run in parallel
	return a.createOrUpdateChasseez(ctx, redfishDevice, service)
}

func (a *Agent) createOrUpdateSystem(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	systemLink := fmt.Sprintf("system-%s", computerSystem.UUID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-system", systemLink, computerSystem)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), computerSystem)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateSystemBIOS(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemBIOS(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	bios, err := computerSystem.Bios()
	if err != nil {
		return
	}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("bios.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-bios", "bios", bios)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), bios)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateSystemBiosAttributes(ctx, redfishDevice, computerSystem, bios)
}

func (a *Agent) createOrUpdateSystemBiosAttributes(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem, bios *redfish.Bios) (err error) {
	parentNode, err := a.getDocument("bios.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	biosAttributes := bios.Attributes

	for biosAttributeName, _ := range biosAttributes {
		biosAttributeValue := biosAttributes.String(biosAttributeName)
		if err = a.createOrUpdateBiosAttribute(ctx, redfishDevice, parentNode, computerSystem, biosAttributeName, biosAttributeValue); err != nil {
			return
		}
	}
	return a.createOrUpdateSystemLed(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateBiosAttribute(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, biosAttributeName, biosAttributeValue string) (err error) {
	biosAttribute := &bootstrap.RedfishBiosAttribute{BiosAttributeValue: biosAttributeValue}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.bios.system-%s.service.%s.redfish-devices.root", biosAttributeName, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-bios-attribute", biosAttributeName, biosAttribute)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), biosAttribute)
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

func (a *Agent) createOrUpdateSystemLed(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	led := &bootstrap.RedfishLed{Led: computerSystem.IndicatorLED}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("led.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-led", "led", led)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), led)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateSystemStatus(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemStatus(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	status := &bootstrap.RedfishStatus{Status: computerSystem.Status}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
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

	return a.createOrUpdateSystemBoot(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemBoot(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	boot := computerSystem.Boot

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("boot.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-boot", "boot", boot)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), boot)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateSystemBootOptions(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemBootOptions(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	parentNode, err := a.getDocument("boot.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	bootOptions, err := computerSystem.BootOptions()
	if err != nil {
		return
	}

	for _, bootOption := range bootOptions {
		if err = a.createOrUpdateBootOption(ctx, redfishDevice, parentNode, computerSystem, bootOption); err != nil {
			return
		}
	}
	return a.createOrUpdateSystemSecureBoot(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateBootOption(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, bootOption *redfish.BootOption) (err error) {
	bootOptionLink := fmt.Sprintf("option-%s", bootOption.ID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.boot.system-%s.service.%s.redfish-devices.root", bootOptionLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-boot-option", bootOptionLink, bootOption)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), bootOption)
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

func (a *Agent) createOrUpdateSystemSecureBoot(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	secureBoot, err := computerSystem.SecureBoot()
	if err != nil {
		return
	}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("secure-boot.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-secure-boot", "secure-boot", secureBoot)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), secureBoot)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdatePCIeDevices(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdatePCIeDevices(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	devices, err := computerSystem.PCIeDevices()
	if err != nil {
		return
	}

	for _, device := range devices {
		if err = a.createOrUpdatePCIeDevice(ctx, redfishDevice, parentNode, computerSystem, device); err != nil {
			return
		}
	}

	// TODO: return createMemories, createProcessors, etc.
	return
}

func (a *Agent) createOrUpdatePCIeDevice(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, device *redfish.PCIeDevice) (err error) {
	deviceLink := fmt.Sprintf("pcie-device-%s", device.ID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-pcie-device", deviceLink, device)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), device)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdatePCIeDeviceStatus(ctx, redfishDevice, computerSystem, deviceLink, device)
}

func (a *Agent) createOrUpdatePCIeDeviceStatus(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem, deviceLink string, device *redfish.PCIeDevice) (err error) {
	status := &bootstrap.RedfishStatus{Status: device.Status}

	parentNode, err := a.getDocument("%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
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

	return a.createOrUpdatePCIeDeviceInterface(ctx, redfishDevice, computerSystem, deviceLink, device)
}

func (a *Agent) createOrUpdatePCIeDeviceInterface(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem, deviceLink string, device *redfish.PCIeDevice) (err error) {
	deviceInterface := &device.PCIeInterface

	parentNode, err := a.getDocument("%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("pcie-interface.%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-pcie-interface", "pcie-interface", deviceInterface)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), deviceInterface)
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

// TODO: currently unavailable structures, pending to implement
func (a *Agent) createOrUpdateMemories(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	memories, err := computerSystem.Memory()
	if err != nil {
		return
	}

	for _, memory := range memories {
		if err = a.createOrUpdateMemory(ctx, redfishDevice, parentNode, computerSystem, memory); err != nil {
			return
		}
	}
	return
}

func (a *Agent) createOrUpdateMemory(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, memory *redfish.Memory) (err error) {
	fmt.Println(memory.DeviceLocator, memory.MemoryLocation)
	return
}
