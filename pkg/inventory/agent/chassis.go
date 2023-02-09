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
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

func (a *Agent) createOrUpdateChasseez(ctx module.Context, redfishDevice device.RedfishDevice, service *gofish.Service) (err error) {
	parentNode, err := a.getDocument("service.%s.redfish-devices.root", redfishDevice.UUID())
	if err != nil {
		return
	}

	chasseez, err := service.Chassis()
	if err != nil {
		return
	}

	for _, chassee := range chasseez {
		if err = a.createOrUpdateChassee(ctx, redfishDevice, parentNode, chassee); err != nil {
			return
		}
	}

	// create/update managers
	// TODO: run in parallel
	return a.createOrUpdateManagers(ctx, redfishDevice, service)
}

// TODO: check Chassis & RedfishDevice UUID, now they are the same
func (a *Agent) createOrUpdateChassee(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	chassisLink := fmt.Sprintf("chassis-%s", chassis.UUID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())

	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-chassis", chassisLink, chassis)
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
	redfishThermal, err := chassis.Thermal()
	if err != nil {
		return
	}

	thermal := &bootstrap.RedfishThermal{
		ID:          redfishThermal.ID,
		Name:        redfishThermal.Name,
		Description: redfishThermal.Description,
	}

	parentNode, err := a.getDocument("chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("thermal.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
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

	return a.createOrUpdateThermalSubsystem(ctx, redfishDevice, chassis, redfishThermal)
}

func (a *Agent) createOrUpdateThermalSubsystem(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, thermal *redfish.Thermal) (err error) {
	parentNode, err := a.getDocument("thermal.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	for _, temperature := range thermal.Temperatures {
		if err = a.createOrUpdateThermalTemperature(ctx, redfishDevice, chassis, parentNode, &temperature); err != nil {
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

func (a *Agent) createOrUpdateThermalTemperature(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, temperature *redfish.Temperature) (err error) {
	temperatureLink := temperature.MemberID

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.thermal.chassis-%s.service.%s.redfish-devices.root", temperatureLink, chassis.UUID, redfishDevice.UUID())
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

	if err = a.createOrUpdateThermalTemperatureStatus(ctx, redfishDevice, chassis, temperatureLink, temperature); err != nil {
		return
	}

	return
}

func (a *Agent) createOrUpdateThermalTemperatureStatus(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, temperatureLink string, temperature *redfish.Temperature) (err error) {
	status := &bootstrap.RedfishStatus{Status: temperature.Status}

	parentNode, err := a.getDocument("%s.thermal.chassis-%s.service.%s.redfish-devices.root", temperatureLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.thermal.chassis-%s.service.%s.redfish-devices.root", temperatureLink, chassis.UUID, redfishDevice.UUID())
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

// TODO: update later, currently not available
func (a *Agent) createOrUpdateThermalFan(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, fan *redfish.Fan) (err error) {
	fanLink := fan.MemberID

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.thermal.chassis-%s.service.%s.redfish-devices.root", fanLink, chassis.UUID, redfishDevice.UUID())
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
	redfishPower, err := chassis.Power()
	if err != nil {
		return
	}

	power := &bootstrap.RedfishPower{
		ID:          redfishPower.ID,
		Name:        redfishPower.Name,
		Description: redfishPower.Description,
	}

	parentNode, err := a.getDocument("chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("power.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
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

	return a.createOrUpdatePowerSubsystem(ctx, redfishDevice, chassis, redfishPower)
}

func (a *Agent) createOrUpdatePowerSubsystem(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, power *redfish.Power) (err error) {
	parentNode, err := a.getDocument("power.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	if err = a.createOrUpdatePowerIndicatorLED(ctx, redfishDevice, chassis, parentNode, power.IndicatorLED); err != nil {
		return
	}

	for _, powerControl := range power.PowerControl {
		if err = a.createOrUpdatePowerControl(ctx, redfishDevice, chassis, parentNode, &powerControl); err != nil {
			return
		}
	}

	for _, powerSupply := range power.PowerSupplies {
		if err = a.createOrUpdatePowerSupply(ctx, redfishDevice, chassis, parentNode, &powerSupply); err != nil {
			return
		}
	}

	for _, voltage := range power.Voltages {
		if err = a.createOrUpdateVoltage(ctx, redfishDevice, chassis, parentNode, &voltage); err != nil {
			return
		}
	}

	return
}

func (a *Agent) createOrUpdatePowerIndicatorLED(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, indicatorLED common.IndicatorLED) (err error) {
	led := &bootstrap.RedfishLed{Led: indicatorLED}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("led.power.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
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
func (a *Agent) createOrUpdatePowerControl(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerControl *redfish.PowerControl) (err error) {
	powerControlLink := fmt.Sprintf("power-control-%s", powerControl.MemberID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.power.chassis-%s.service.%s.redfish-devices.root", powerControlLink, chassis.UUID, redfishDevice.UUID())
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
func (a *Agent) createOrUpdatePowerSupply(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerSupply *redfish.PowerSupply) (err error) {
	powerSupplyLink := fmt.Sprintf("power-supply-%s", powerSupply.MemberID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.power.chassis-%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
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
func (a *Agent) createOrUpdateVoltage(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, voltage *redfish.Voltage) (err error) {
	voltageLink := fmt.Sprintf("voltage-%s", voltage.MemberID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.power.chassis-%s.service.%s.redfish-devices.root", voltageLink, chassis.UUID, redfishDevice.UUID())
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
