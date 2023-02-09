// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/bootstrap"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

func (a *Agent) workerFunction(ctx module.Context) (err error) {
	var redfishDevice device.RedfishDevice
	if err = json.Unmarshal(ctx.CmdbContext(), &redfishDevice); err != nil {
		return
	}

	u, err := url.Parse(redfishDevice.Api)
	if err != nil {
		return
	}

	config := gofish.ClientConfig{
		Endpoint: fmt.Sprintf("%s://%s", u.Scheme, u.Host),
		Username: redfishDevice.Login,
		Password: redfishDevice.Password,
		Insecure: true,
	}

	client, err := gofish.ConnectContext(ctx, config)
	if err != nil {
		return
	}
	defer client.Logout()

	service := client.GetService()

	// create/update main (root) service
	if err = a.createOrUpdateService(ctx, redfishDevice, service); err != nil {
		return
	}

	return nil
}

func (a *Agent) createOrUpdateService(ctx module.Context, redfishDevice device.RedfishDevice, service *gofish.Service) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("service.%s.redfish-devices.root", redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(ctx.Self().Id, "types/redfish-service", "service", service)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), service)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	// create/update systems service
	return a.createOrUpdateSystems(ctx, redfishDevice, service)
}

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

func (a *Agent) createOrUpdateChasseez(ctx module.Context, redfishDevice device.RedfishDevice, service *gofish.Service) (err error) {
	parentNode, err := a.getDocument("service.%s.redfish-devices.root", redfishDevice.UUID())
	if err != nil {
		return
	}

	chassis, err := service.Chassis()
	if err != nil {
		return
	}

	for _, chs := range chassis {
		if err = a.createOrUpdateChassee(ctx, redfishDevice, parentNode, chs); err != nil {
			return
		}
	}

	// create/update managers
	return a.createOrUpdateManagers(ctx, redfishDevice, service)
}

// TODO: check Chassis & RedfishDevice UUID
func (a *Agent) createOrUpdateChassee(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())

	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-chassis", chassis.UUID, chassis)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), chassis)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateThermal(ctx, redfishDevice, chassis)
}

func (a *Agent) createOrUpdateThermal(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis) (err error) {
	rfThermal, err := chassis.Thermal()
	if err != nil {
		return
	}

	thermal := &bootstrap.RedfishThermal{
		ID:          rfThermal.ID,
		Name:        rfThermal.Name,
		Description: rfThermal.Description,
	}

	parentNode, err := a.getDocument("%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("thermal.%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-thermal", "thermal", thermal)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), thermal)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateThermalSubsystem(ctx, redfishDevice, chassis, rfThermal)
}

func (a *Agent) createOrUpdateThermalSubsystem(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, thermal *redfish.Thermal) (err error) {
	parentNode, err := a.getDocument("thermal.%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	for _, temp := range thermal.Temperatures {
		if err = a.createOrUpdateThermalTemperature(ctx, redfishDevice, chassis, parentNode, &temp); err != nil {
			return
		}
	}

	for _, fan := range thermal.Fans {
		if err = a.createOrUpdateThermalFan(ctx, redfishDevice, chassis, parentNode, &fan); err != nil {
			return
		}
	}

	return a.createOrUpdatePower(ctx, redfishDevice, chassis)
}

func (a *Agent) createOrUpdateThermalTemperature(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, rfTemperature *redfish.Temperature) (err error) {
	temperatureLink := strings.ToLower(rfTemperature.MemberID)
	temperature := &bootstrap.RedfishTemperature{Temperature: rfTemperature}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.thermal.%s.service.%s.redfish-devices.root", temperatureLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-temperature", temperatureLink, temperature)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), temperature)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	if err = a.createOrUpdateThermalTemperatureStatus(ctx, redfishDevice, chassis, temperatureLink, rfTemperature); err != nil {
		return
	}

	return
}

func (a *Agent) createOrUpdateThermalTemperatureStatus(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, temperatureLink string, rfTemperature *redfish.Temperature) (err error) {
	status := &bootstrap.RedfishStatus{Status: rfTemperature.Status}
	parentNode, err := a.getDocument("%s.thermal.%s.service.%s.redfish-devices.root", temperatureLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.thermal.%s.service.%s.redfish-devices.root", temperatureLink, chassis.UUID, redfishDevice.UUID())
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

	return
}

func (a *Agent) createOrUpdateThermalFan(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, rfFan *redfish.Fan) (err error) {
	fanLink := strings.ToLower(rfFan.MemberID)
	fan := &bootstrap.RedfishFan{Fan: rfFan}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.thermal.%s.service.%s.redfish-devices.root", fanLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-fan", fanLink, fan)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), fan)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return
}

func (a *Agent) createOrUpdatePower(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis) (err error) {
	rfPower, err := chassis.Power()
	if err != nil {
		return
	}

	power := &bootstrap.RedfishPower{
		ID:          rfPower.ID,
		Name:        rfPower.Name,
		Description: rfPower.Description,
	}

	parentNode, err := a.getDocument("%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("power.%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-power", "power", power)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), power)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdatePowerSubsystem(ctx, redfishDevice, chassis, rfPower)
}

func (a *Agent) createOrUpdatePowerSubsystem(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, power *redfish.Power) (err error) {
	parentNode, err := a.getDocument("power.%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	if err = a.createOrUpdatePowerIndicatorLED(ctx, redfishDevice, chassis, parentNode, power.IndicatorLED); err != nil {
		return
	}

	for _, pc := range power.PowerControl {
		if err = a.createOrUpdatePowerControl(ctx, redfishDevice, chassis, parentNode, &pc); err != nil {
			return
		}
	}

	for _, ps := range power.PowerSupplies {
		if err = a.createOrUpdatePowerSupply(ctx, redfishDevice, chassis, parentNode, &ps); err != nil {
			return
		}
	}

	for _, v := range power.Voltages {
		if err = a.createOrUpdateVoltage(ctx, redfishDevice, chassis, parentNode, &v); err != nil {
			return
		}
	}

	return
}

func (a *Agent) createOrUpdatePowerIndicatorLED(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, indicatorLED common.IndicatorLED) (err error) {
	led := &bootstrap.RedfishLed{Led: indicatorLED}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("led.power.%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
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

	return
}

// TODO: unique link
// TODO: device with Id is present ? create : ignore
func (a *Agent) createOrUpdatePowerControl(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, rfPowerControl *redfish.PowerControl) (err error) {
	powerControlLink := fmt.Sprintf("pcontrol%s", strings.ToLower(rfPowerControl.MemberID))
	powerControl := &bootstrap.RedfishPowerControl{PowerControl: rfPowerControl}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.power.%s.service.%s.redfish-devices.root", powerControlLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-power-control", powerControlLink, powerControl)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), powerControl)
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

// TODO: unique link
// TODO: device with Id is present ? create : ignore
func (a *Agent) createOrUpdatePowerSupply(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, rfPowerSupply *redfish.PowerSupply) (err error) {
	powerSupplyLink := fmt.Sprintf("psupply%s", strings.ToLower(rfPowerSupply.MemberID))
	powerSupply := &bootstrap.RedfishPowerSupply{PowerSupply: rfPowerSupply}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.power.%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-power-supply", powerSupplyLink, powerSupply)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), powerSupply)
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

// TODO: unique link
// TODO: device with Id is present ? create : ignore
func (a *Agent) createOrUpdateVoltage(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, rfVoltage *redfish.Voltage) (err error) {
	voltageLink := strings.ToLower(rfVoltage.MemberID)
	voltage := &bootstrap.RedfishVoltage{Voltage: rfVoltage}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.power.%s.service.%s.redfish-devices.root", voltageLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-voltage", voltageLink, voltage)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), voltage)
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
		functionContext, err = system.UpdateObject(document.Id.String(), managerLink)
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
