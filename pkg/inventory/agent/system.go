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
	return a.createOrUpdateChasseez(ctx, redfishDevice, service)
}

func (a *Agent) createOrUpdateSystem(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-system", computerSystem.UUID, computerSystem)
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

	parentNode, err := a.getDocument("%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("bios.%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
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

	return a.createOrUpdateSystemLed(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemLed(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	led := &bootstrap.RedfishLed{Led: computerSystem.IndicatorLED}

	parentNode, err := a.getDocument("%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("led.%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
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

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateSystemStatus(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemStatus(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	status := &bootstrap.RedfishStatus{Status: computerSystem.Status}

	parentNode, err := a.getDocument("%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
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

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateSystemBoot(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemBoot(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	boot := &bootstrap.RedfishBoot{Boot: &computerSystem.Boot}

	parentNode, err := a.getDocument("%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("boot.%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
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
	parentNode, err := a.getDocument("boot.%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	options, err := computerSystem.BootOptions()
	if err != nil {
		return
	}

	for _, opt := range options {
		if err = a.createOrUpdateBootOption(ctx, redfishDevice, parentNode, computerSystem, opt); err != nil {
			return
		}
	}
	return a.createOrUpdateSystemSecureBoot(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateBootOption(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, bootOption *redfish.BootOption) (err error) {
	optionLink := strings.ToLower(bootOption.ID)
	option := &bootstrap.RedfishBootOption{BootOption: bootOption}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.boot.%s.service.%s.redfish-devices.root", optionLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-boot-option", optionLink, option)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), option)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return
}

func (a *Agent) createOrUpdateSystemSecureBoot(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	secureBoot, err := computerSystem.SecureBoot()
	if err != nil {
		return
	}

	parentNode, err := a.getDocument("%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("secure-boot.%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
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

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdatePCIeDevices(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdatePCIeDevices(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	parentNode, err := a.getDocument("%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
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
	return
}

func (a *Agent) createOrUpdatePCIeDevice(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, rfDevice *redfish.PCIeDevice) (err error) {
	deviceLink := strings.ToLower(rfDevice.ID)
	device := &bootstrap.RedfishPcieDevice{PCIeDevice: rfDevice}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
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

	return a.createOrUpdatePCIeDeviceStatus(ctx, redfishDevice, computerSystem, deviceLink, rfDevice)
}

func (a *Agent) createOrUpdatePCIeDeviceStatus(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem, deviceLink string, device *redfish.PCIeDevice) (err error) {
	status := &bootstrap.RedfishStatus{Status: device.Status}
	parentNode, err := a.getDocument("%s.%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
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
	iface := &bootstrap.RedfishPcieInterface{PCIeInterface: &device.PCIeInterface}
	parentNode, err := a.getDocument("%s.%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("pcie-interface.%s.%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-pcie-interface", "pcie-interface", iface)
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

	return
}

func (a *Agent) createOrUpdateMemories(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	parentNode, err := a.getDocument("%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
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
