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
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

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
