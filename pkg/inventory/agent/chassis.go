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

	return a.createOrUpdateThermalTemperatureStatus(ctx, redfishDevice, chassis, temperatureLink, temperature)
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

	return a.createOrUpdateChassisLed(ctx, redfishDevice, chassis)
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

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createPowerControlDetails(ctx, redfishDevice, chassis, powerControlLink, powerControl)
}

func (a *Agent) createPowerControlDetails(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, powerControlLink string, powerControl *redfish.PowerControl) (err error) {
	parentNode, err := a.getDocument("%s.power.chassis-%s.service.%s.redfish-devices.root", powerControlLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	if err = a.createOrUpdatePowerControlStatus(ctx, redfishDevice, chassis, parentNode, powerControlLink, powerControl); err != nil {
		return
	}

	if err = a.createOrUpdatePowerControlPhysicalContext(ctx, redfishDevice, chassis, parentNode, powerControlLink, powerControl); err != nil {
		return
	}

	if err = a.createOrUpdatePowerControlMetric(ctx, redfishDevice, chassis, parentNode, powerControlLink, powerControl); err != nil {
		return
	}

	if err = a.createOrUpdatePowerControlLimit(ctx, redfishDevice, chassis, parentNode, powerControlLink, powerControl); err != nil {
		return
	}

	return
}

func (a *Agent) createOrUpdatePowerControlStatus(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerControlLink string, powerControl *redfish.PowerControl) (err error) {
	status := &bootstrap.RedfishStatus{Status: powerControl.Status}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.power.chassis-%s.service.%s.redfish-devices.root", powerControlLink, chassis.UUID, redfishDevice.UUID())
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

func (a *Agent) createOrUpdatePowerControlPhysicalContext(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerControlLink string, powerControl *redfish.PowerControl) (err error) {
	physicalContext := &bootstrap.RedfishPhysicalContext{PhysicalContext: &powerControl.PhysicalContext}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("physcial-context.%s.power.chassis-%s.service.%s.redfish-devices.root", powerControlLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-physical-context", "physcial-context", physicalContext)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), physicalContext)
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

func (a *Agent) createOrUpdatePowerControlMetric(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerControlLink string, powerControl *redfish.PowerControl) (err error) {
	powerMetric := powerControl.PowerMetrics

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("power-metric.%s.power.chassis-%s.service.%s.redfish-devices.root", powerControlLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-power-metric", "power-metric", powerMetric)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), powerMetric)
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

func (a *Agent) createOrUpdatePowerControlLimit(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerControlLink string, powerControl *redfish.PowerControl) (err error) {
	powerLimit := powerControl.PowerLimit

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("power-limit.%s.power.chassis-%s.service.%s.redfish-devices.root", powerControlLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-power-limit", "power-limit", powerLimit)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), powerLimit)
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

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createPowerSupplyDetails(ctx, redfishDevice, chassis, powerSupplyLink, powerSupply)
}

func (a *Agent) createPowerSupplyDetails(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, powerSupplyLink string, powerSupply *redfish.PowerSupply) (err error) {
	parentNode, err := a.getDocument("%s.power.chassis-%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	if err = a.createOrUpdatePowerSupplyStatus(ctx, redfishDevice, chassis, parentNode, powerSupplyLink, powerSupply); err != nil {
		return
	}

	if err = a.createOrUpdatePowerSupplyIndicatorLED(ctx, redfishDevice, chassis, parentNode, powerSupplyLink, powerSupply); err != nil {
		return
	}

	if err = a.createOrUpdatePowerSupplyLocation(ctx, redfishDevice, chassis, parentNode, powerSupplyLink, powerSupply); err != nil {
		return
	}

	return
}

func (a *Agent) createOrUpdatePowerSupplyStatus(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, powerSupply *redfish.PowerSupply) (err error) {
	status := &bootstrap.RedfishStatus{Status: powerSupply.Status}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.power.chassis-%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
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

func (a *Agent) createOrUpdatePowerSupplyIndicatorLED(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, powerSupply *redfish.PowerSupply) (err error) {
	led := &bootstrap.RedfishLed{Led: powerSupply.IndicatorLED}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("led.%s.power.chassis-%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
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

func (a *Agent) createOrUpdatePowerSupplyLocation(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, powerSupply *redfish.PowerSupply) (err error) {
	location := powerSupply.Location

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("location.%s.power.chassis-%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-location", "location", location)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), location)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdatePowerSupplyLocationDetails(ctx, redfishDevice, chassis, powerSupplyLink, &location)
}

func (a *Agent) createOrUpdatePowerSupplyLocationDetails(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, powerSupplyLink string, location *common.Location) (err error) {
	parentNode, err := a.getDocument("location.%s.power.chassis-%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	if err = a.createOrUpdatePowerSupplyPartLocation(ctx, redfishDevice, chassis, parentNode, powerSupplyLink, &location.PartLocation); err != nil {
		return
	}

	if err = a.createOrUpdatePowerSupplyPlacement(ctx, redfishDevice, chassis, parentNode, powerSupplyLink, &location.Placement); err != nil {
		return
	}

	if err = a.createOrUpdatePowerSupplyPostalAddress(ctx, redfishDevice, chassis, parentNode, powerSupplyLink, &location.PostalAddress); err != nil {
		return
	}

	return a.createOrUpdateChassisSupportedResetTypes(ctx, redfishDevice, chassis)
}

func (a *Agent) createOrUpdatePowerSupplyPartLocation(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, partLocation *common.PartLocation) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("part-location.location.%s.power.chassis-%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-part-location", "part-location", partLocation)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), partLocation)
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

func (a *Agent) createOrUpdatePowerSupplyPlacement(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, placement *common.Placement) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("placement.location.%s.power.chassis-%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-placement", "placement", placement)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), placement)
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

func (a *Agent) createOrUpdatePowerSupplyPostalAddress(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, postalAddress *common.PostalAddress) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("postal-address.location.%s.power.chassis-%s.service.%s.redfish-devices.root", powerSupplyLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-postal-address", "postal-address", postalAddress)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), postalAddress)
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

	// FIXME: context canceled
	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return

	// if err = a.executor.ExecSync(ctx, functionContext); err != nil {
	// 	return
	// }

	// return a.createOrUpdateVoltageStatus(ctx, redfishDevice, chassis, voltageLink, voltage)
}

func (a *Agent) createOrUpdateVoltageStatus(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, voltageLink string, voltage *redfish.Voltage) (err error) {
	status := &bootstrap.RedfishStatus{Status: voltage.Status}

	parentNode, err := a.getDocument("%s.power.chassis-%s.service.%s.redfish-devices.root", voltageLink, chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.power.chassis-%s.service.%s.redfish-devices.root", voltageLink, chassis.UUID, redfishDevice.UUID())
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

func (a *Agent) createOrUpdateChassisLed(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis) (err error) {
	led := &bootstrap.RedfishLed{Led: chassis.IndicatorLED}

	parentNode, err := a.getDocument("chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("led.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
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

	return a.createOrUpdateChassisStatus(ctx, redfishDevice, chassis)
}

func (a *Agent) createOrUpdateChassisStatus(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis) (err error) {
	status := &bootstrap.RedfishStatus{Status: chassis.Status}

	parentNode, err := a.getDocument("chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
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

	return a.createOrUpdateChassisPowerState(ctx, redfishDevice, chassis)
}

func (a *Agent) createOrUpdateChassisPowerState(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis) (err error) {
	powerState := &bootstrap.RedfishPowerState{PowerState: chassis.PowerState}

	parentNode, err := a.getDocument("chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("power-state.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
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

	return a.createOrUpdateChassisPhysicalSecurity(ctx, redfishDevice, chassis)
}

func (a *Agent) createOrUpdateChassisPhysicalSecurity(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis) (err error) {
	physicalSecurity := chassis.PhysicalSecurity

	parentNode, err := a.getDocument("chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("physical-security.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-physical-security", "physical-security", physicalSecurity)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), physicalSecurity)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateChassisLocation(ctx, redfishDevice, chassis)
}

func (a *Agent) createOrUpdateChassisLocation(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis) (err error) {
	location := chassis.Location

	parentNode, err := a.getDocument("chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("location.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-location", "location", location)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), location)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateChassisLocationDetails(ctx, redfishDevice, chassis, &location)
}

func (a *Agent) createOrUpdateChassisLocationDetails(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, location *common.Location) (err error) {
	parentNode, err := a.getDocument("location.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	if err = a.createOrUpdatePartLocation(ctx, redfishDevice, chassis, parentNode, &location.PartLocation); err != nil {
		return
	}

	if err = a.createOrUpdatePlacement(ctx, redfishDevice, chassis, parentNode, &location.Placement); err != nil {
		return
	}

	if err = a.createOrUpdatePostalAddress(ctx, redfishDevice, chassis, parentNode, &location.PostalAddress); err != nil {
		return
	}

	return a.createOrUpdateChassisSupportedResetTypes(ctx, redfishDevice, chassis)
}

func (a *Agent) createOrUpdatePartLocation(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, partLocation *common.PartLocation) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("part-location.location.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-part-location", "part-location", partLocation)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), partLocation)
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

func (a *Agent) createOrUpdatePlacement(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, placement *common.Placement) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("placement.location.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-placement", "placement", placement)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), placement)
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

func (a *Agent) createOrUpdatePostalAddress(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis, parentNode *documents.Node, postalAddress *common.PostalAddress) (err error) {
	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("postal-address.location.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-postal-address", "postal-address", postalAddress)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), postalAddress)
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

func (a *Agent) createOrUpdateChassisSupportedResetTypes(ctx module.Context, redfishDevice device.RedfishDevice, chassis *redfish.Chassis) (err error) {
	supportedResetTypes := &bootstrap.RedfishSupportedResetTypes{SupportedResetTypes: chassis.SupportedResetTypes}

	parentNode, err := a.getDocument("chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("supported-reset-types.chassis-%s.service.%s.redfish-devices.root", chassis.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-supported-reset-types", "supported-reset-types", supportedResetTypes)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), supportedResetTypes)
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
