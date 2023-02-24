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
	// Chassis
	chassisMask                      = "chassis-%s.service.*[?@._id == '%s'?].objects.root"
	statusChassisMask                = "status.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	ledChassisMask                   = "led.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	locationChassisMask              = "location.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	partLocationLocationChassisMask  = "part-location.location.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	placementLocationChassisMask     = "placement.location.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	postalAddressLocationChassisMask = "postal-address.location.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	powerStateChassisMask            = "power-state.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	physicalSecurityChassisMask      = "physical-security.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	supportedResetTypesChassisMask   = "supported-reset-types.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	subChassisMask                   = "%s.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	statusSubChassisMask             = "status.%s.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	subSubChassisMask                = "%s.%s.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	statusSubSubChassisMask          = "status.%s.%s.chassis-%s.service.*[?@._id == '%s'?].objects.root"

	// Chassis -> Thermal Subsystem
	thermalChassisMask                         = "thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	subThermalChassisMask                      = "%s.thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	statusSubThermalChassisMask                = "status.%s.thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	ledSubThermalChassisMask                   = "led.%s.thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	locationSubThermalChassisMask              = "location.%s.thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	partLocationLocationSubThermalChassisMask  = "part-location.location.%s.thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	placementLocationSubThermalChassisMask     = "placement.location.%s.thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	postalAddressLocationSubThermalChassisMask = "postal-address.location.%s.thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"

	// Chassis -> Power Subsystem
	powerChassisMask                         = "power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	ledPowerChassisMask                      = "led.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	subPowerChassisMask                      = "%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	statusSubPowerChassisMask                = "status.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	ledSubPowerChassisMask                   = "led.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	locationSubPowerChassisMask              = "location.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	partLocationLocationSubPowerChassisMask  = "part-location.location.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	placementLocationSubPowerChassisMask     = "placement.location.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	postalAddressLocationSubPowerChassisMask = "postal-address.location.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	physicalContextSubPowerChassisMask       = "physcial-context.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	powerMetricSubPowerChassisMask           = "power-metric.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	powerLimitSubPowerChassisMask            = "power-limit.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
)

func (a *Agent) createOrUpdateChasseez(ctx module.Context, service *gofish.Service, parentNode *documents.Node, vendorData *VendorSpecificData) (err error) {
	chasseez, err := service.Chassis()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, chassee := range chasseez {
		chassee := chassee
		p.Exec(func() error { return a.createOrUpdateChassee(ctx, parentNode, chassee, vendorData) })
	}
	return p.Wait()
}

// TODO: check Chassis & RedfishDevice UUID, now they are the same
func (a *Agent) createOrUpdateChassee(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, vendorData *VendorSpecificData) (err error) {
	chassisLink := fmt.Sprintf("chassis-%s", chassis.UUID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishChassisID, chassisLink, chassis, chassisMask, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return a.createOrUpdateThermal(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdatePower(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateNetworkAdapters(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisLed(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisStatus(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisPowerState(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisPhysicalSecurity(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisLocation(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisSupportedResetTypes(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisSELLogService(ctx, document, chassis, vendorData) })
	return p.Wait()
}

func (a *Agent) createOrUpdateThermal(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	redfishThermal, err := chassis.Thermal()
	if err != nil {
		return
	}

	thermal := &bootstrap.RedfishThermal{
		ID:          redfishThermal.ID,
		Name:        redfishThermal.Name,
		Description: redfishThermal.Description,
	}

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishThermalID, types.RedfishThermalLink, thermal, thermalChassisMask, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, temperature := range redfishThermal.Temperatures {
		temperature := temperature
		p.Exec(func() error { return a.createOrUpdateThermalTemperature(ctx, document, chassis, temperature) })
	}
	for _, fan := range redfishThermal.Fans {
		fan := fan
		p.Exec(func() error { return a.createOrUpdateThermalFan(ctx, document, chassis, fan) })
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateThermalTemperature(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, temperature redfish.Temperature) (err error) {
	// FIXME: added to avoid link name conflicts with fans
	temperatureLink := fmt.Sprintf("temperature-%s", temperature.MemberID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishTemperatureID, temperatureLink, temperature, subThermalChassisMask, temperatureLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateThermalTemperatureStatus(ctx, document, chassis, temperatureLink, temperature)
}

func (a *Agent) createOrUpdateThermalTemperatureStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, temperatureLink string, temperature redfish.Temperature) (err error) {
	status := &bootstrap.RedfishStatus{Status: temperature.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubThermalChassisMask, temperatureLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateThermalFan(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, fan redfish.Fan) (err error) {
	fanLink := fmt.Sprintf("fan-%s", fan.MemberID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishFanID, fanLink, fan, subThermalChassisMask, fanLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdateFanStatus(ctx, document, chassis, fanLink, fan)
	})
	p.Exec(func() error {
		return a.createOrUpdateFanIndicatorLED(ctx, document, chassis, fanLink, fan)
	})
	p.Exec(func() error {
		return a.createOrUpdateFanLocation(ctx, document, chassis, fanLink, fan)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdateFanStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, fanLink string, fan redfish.Fan) (err error) {
	status := &bootstrap.RedfishStatus{Status: fan.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubThermalChassisMask, fanLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateFanIndicatorLED(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, fanLink string, fan redfish.Fan) (err error) {
	led := &bootstrap.RedfishLed{Led: fan.IndicatorLED}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLedID, types.RedfishLedLink, led, ledSubThermalChassisMask, fanLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateFanLocation(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, fanLink string, fan redfish.Fan) (err error) {
	location := fan.Location

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLocationID, types.RedfishLocationLink, location, locationSubThermalChassisMask, fanLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdateFanPartLocation(ctx, document, chassis, fanLink, location.PartLocation)
	})
	p.Exec(func() error {
		return a.createOrUpdateFanPlacement(ctx, document, chassis, fanLink, location.Placement)
	})
	p.Exec(func() error {
		return a.createOrUpdateFanPostalAddress(ctx, document, chassis, fanLink, location.PostalAddress)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdateFanPartLocation(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, fanLink string, partLocation common.PartLocation) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPartLocationID, types.RedfishPartLocationLink, partLocation, partLocationLocationSubThermalChassisMask, fanLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateFanPlacement(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, fanLink string, placement common.Placement) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPlacementID, types.RedfishPlacementLink, placement, placementLocationSubThermalChassisMask, fanLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateFanPostalAddress(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, fanLink string, postalAddress common.PostalAddress) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPostalAddressID, types.RedfishPostalAddressLink, postalAddress, postalAddressLocationSubThermalChassisMask, fanLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePower(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	redfishPower, err := chassis.Power()
	if err != nil {
		return
	}

	power := &bootstrap.RedfishPower{
		ID:          redfishPower.ID,
		Name:        redfishPower.Name,
		Description: redfishPower.Description,
	}

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerID, types.RedfishPowerLink, power, powerChassisMask, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdatePowerIndicatorLED(ctx, document, chassis, redfishPower.IndicatorLED)
	})
	for _, powerControl := range redfishPower.PowerControl {
		powerControl := powerControl
		p.Exec(func() error { return a.createOrUpdatePowerControl(ctx, document, chassis, powerControl) })
	}
	for _, powerSupply := range redfishPower.PowerSupplies {
		powerSupply := powerSupply
		p.Exec(func() error { return a.createOrUpdatePowerSupply(ctx, document, chassis, powerSupply) })
	}
	for _, voltage := range redfishPower.Voltages {
		voltage := voltage
		p.Exec(func() error { return a.createOrUpdateVoltage(ctx, document, chassis, voltage) })
	}
	return p.Wait()
}

func (a *Agent) createOrUpdatePowerIndicatorLED(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, indicatorLED common.IndicatorLED) (err error) {
	led := &bootstrap.RedfishLed{Led: indicatorLED}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLedID, types.RedfishLedLink, led, ledPowerChassisMask, chassis.UUID, ctx.Self().Id)
}

// TODO: unique link
// TODO: device with Id is present ? create : ignore
func (a *Agent) createOrUpdatePowerControl(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerControl redfish.PowerControl) (err error) {
	powerControlLink := fmt.Sprintf("power-control-%s", powerControl.MemberID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerControlID, powerControlLink, powerControl, subPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdatePowerControlStatus(ctx, document, chassis, powerControlLink, powerControl)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerControlPhysicalContext(ctx, document, chassis, powerControlLink, powerControl)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerControlMetric(ctx, document, chassis, powerControlLink, powerControl)
	})

	p.Exec(func() error {
		return a.createOrUpdatePowerControlLimit(ctx, document, chassis, powerControlLink, powerControl)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdatePowerControlStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerControlLink string, powerControl redfish.PowerControl) (err error) {
	status := &bootstrap.RedfishStatus{Status: powerControl.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerControlPhysicalContext(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerControlLink string, powerControl redfish.PowerControl) (err error) {
	physicalContext := &bootstrap.RedfishPhysicalContext{PhysicalContext: &powerControl.PhysicalContext}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPhysicalContextID, types.RedfishPhysicalContextLink, physicalContext, physicalContextSubPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerControlMetric(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerControlLink string, powerControl redfish.PowerControl) (err error) {
	powerMetric := powerControl.PowerMetrics
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerMetricID, types.RedfishPowerMetricLink, powerMetric, powerMetricSubPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerControlLimit(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerControlLink string, powerControl redfish.PowerControl) (err error) {
	powerLimit := powerControl.PowerLimit
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerLimitID, types.RedfishPowerLimitLink, powerLimit, powerLimitSubPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
}

// TODO: unique link
// TODO: device with Id is present ? create : ignore
func (a *Agent) createOrUpdatePowerSupply(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerSupply redfish.PowerSupply) (err error) {
	powerSupplyLink := fmt.Sprintf("power-supply-%s", powerSupply.MemberID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerSupplyID, powerSupplyLink, powerSupply, subPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyStatus(ctx, document, chassis, powerSupplyLink, powerSupply)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyIndicatorLED(ctx, document, chassis, powerSupplyLink, powerSupply)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyLocation(ctx, document, chassis, powerSupplyLink, powerSupply)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdatePowerSupplyStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerSupplyLink string, powerSupply redfish.PowerSupply) (err error) {
	status := &bootstrap.RedfishStatus{Status: powerSupply.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerSupplyIndicatorLED(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerSupplyLink string, powerSupply redfish.PowerSupply) (err error) {
	led := &bootstrap.RedfishLed{Led: powerSupply.IndicatorLED}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLedID, types.RedfishLedLink, led, ledSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerSupplyLocation(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerSupplyLink string, powerSupply redfish.PowerSupply) (err error) {
	location := powerSupply.Location

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLocationID, types.RedfishLocationLink, location, locationSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyPartLocation(ctx, document, chassis, powerSupplyLink, location.PartLocation)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyPlacement(ctx, document, chassis, powerSupplyLink, location.Placement)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyPostalAddress(ctx, document, chassis, powerSupplyLink, location.PostalAddress)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdatePowerSupplyPartLocation(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerSupplyLink string, partLocation common.PartLocation) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPartLocationID, types.RedfishPartLocationLink, partLocation, partLocationLocationSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerSupplyPlacement(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerSupplyLink string, placement common.Placement) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPlacementID, types.RedfishPlacementLink, placement, placementLocationSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerSupplyPostalAddress(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, powerSupplyLink string, postalAddress common.PostalAddress) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPostalAddressID, types.RedfishPostalAddressLink, postalAddress, postalAddressLocationSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

// TODO: unique link
// TODO: device with Id is present ? create : ignore
func (a *Agent) createOrUpdateVoltage(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, voltage redfish.Voltage) (err error) {
	voltageLink := fmt.Sprintf("voltage-%s", voltage.MemberID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishVoltageID, voltageLink, voltage, subPowerChassisMask, voltageLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateVoltageStatus(ctx, document, chassis, voltageLink, voltage)
}

func (a *Agent) createOrUpdateVoltageStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, voltageLink string, voltage redfish.Voltage) (err error) {
	status := &bootstrap.RedfishStatus{Status: voltage.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubPowerChassisMask, voltageLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateNetworkAdapters(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	networkAdapters, err := chassis.NetworkAdapters()
	if err != nil {
		return
	}
	p := utils.NewParallel()
	for _, networkAdapter := range networkAdapters {
		networkAdapter := networkAdapter
		p.Exec(func() error { return a.createOrUpdateNetworkAdapter(ctx, parentNode, chassis, networkAdapter) })
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateNetworkAdapter(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, networkAdapter *redfish.NetworkAdapter) (err error) {
	networkAdapterLink := networkAdapter.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishNetworkAdapterID, networkAdapterLink, networkAdapter, subChassisMask, networkAdapterLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdateNetworkAdapterStatus(ctx, document, chassis, networkAdapterLink, networkAdapter)
	})
	p.Exec(func() error {
		return a.createOrUpdateNetworkAdapterDeviceFunctions(ctx, document, chassis, networkAdapterLink, networkAdapter)
	})
	p.Exec(func() error {
		return a.createOrUpdateNetworkAdapterPorts(ctx, document, chassis, networkAdapterLink, networkAdapter)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdateNetworkAdapterStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, networkAdapterLink string, networkAdapter *redfish.NetworkAdapter) (err error) {
	status := &bootstrap.RedfishStatus{Status: networkAdapter.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubChassisMask, networkAdapterLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateNetworkAdapterDeviceFunctions(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, networkAdapterLink string, networkAdapter *redfish.NetworkAdapter) (err error) {
	deviceFunctions, err := networkAdapter.NetworkDeviceFunctions()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, deviceFunction := range deviceFunctions {
		deviceFunction := deviceFunction
		p.Exec(func() error {
			return a.createOrUpdateNetworkAdapterDeviceFunction(ctx, parentNode, chassis, networkAdapterLink, deviceFunction)
		})
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateNetworkAdapterDeviceFunction(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, networkAdapterLink string, deviceFunction *redfish.NetworkDeviceFunction) (err error) {
	deviceFunctionLink := deviceFunction.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishNetworkDeviceFunctionID, deviceFunctionLink, deviceFunction, subSubChassisMask, deviceFunctionLink, networkAdapterLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateNetworkAdapterDeviceFunctionStatus(ctx, document, chassis, networkAdapterLink, deviceFunctionLink, deviceFunction)
}

func (a *Agent) createOrUpdateNetworkAdapterDeviceFunctionStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, networkAdapterLink, deviceFunctionLink string, deviceFunction *redfish.NetworkDeviceFunction) (err error) {
	status := &bootstrap.RedfishStatus{Status: deviceFunction.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSubChassisMask, deviceFunctionLink, networkAdapterLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateNetworkAdapterPorts(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, networkAdapterLink string, networkAdapter *redfish.NetworkAdapter) (err error) {
	// TODO: check [] networkadapter.networkPorts
	networkPorts, err := networkAdapter.NetworkPorts()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, networkPort := range networkPorts {
		networkPort := networkPort
		p.Exec(func() error {
			return a.createOrUpdateNetworkAdapterPort(ctx, parentNode, chassis, networkAdapterLink, networkPort)
		})
	}

	return p.Wait()
}

func (a *Agent) createOrUpdateNetworkAdapterPort(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, networkAdapterLink string, networkPort *redfish.NetworkPort) (err error) {
	networkPortLink := networkPort.ID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishNetworkPortID, networkPortLink, networkPort, subSubChassisMask, networkPortLink, networkAdapterLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateNetworkAdapterPortStatus(ctx, document, chassis, networkAdapterLink, networkPortLink, networkPort)
}

func (a *Agent) createOrUpdateNetworkAdapterPortStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, networkAdapterLink, networkPortLink string, networkPort *redfish.NetworkPort) (err error) {
	status := &bootstrap.RedfishStatus{Status: networkPort.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSubChassisMask, networkPortLink, networkAdapterLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateChassisLed(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	led := &bootstrap.RedfishLed{Led: chassis.IndicatorLED}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLedID, types.RedfishLedLink, led, ledChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateChassisStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	status := &bootstrap.RedfishStatus{Status: chassis.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateChassisPowerState(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	powerState := &bootstrap.RedfishPowerState{PowerState: chassis.PowerState}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerStateID, types.RedfishPowerStateLink, powerState, powerStateChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateChassisPhysicalSecurity(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	physicalSecurity := chassis.PhysicalSecurity
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPhysicalSecurityID, types.RedfishPhysicalSecurityLink, physicalSecurity, physicalSecurityChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateChassisLocation(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	location := chassis.Location

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLocationID, types.RedfishLocationLink, location, locationChassisMask, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return a.createOrUpdatePartLocation(ctx, document, chassis, location.PartLocation) })
	p.Exec(func() error { return a.createOrUpdatePlacement(ctx, document, chassis, location.Placement) })
	p.Exec(func() error { return a.createOrUpdatePostalAddress(ctx, document, chassis, location.PostalAddress) })
	return p.Wait()
}

func (a *Agent) createOrUpdatePartLocation(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, partLocation common.PartLocation) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPartLocationID, types.RedfishPartLocationLink, partLocation, partLocationLocationChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePlacement(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, placement common.Placement) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPlacementID, types.RedfishPlacementLink, placement, placementLocationChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePostalAddress(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, postalAddress common.PostalAddress) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPostalAddressID, types.RedfishPostalAddressLink, postalAddress, postalAddressLocationChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateChassisSupportedResetTypes(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	supportedResetTypes := &bootstrap.RedfishSupportedResetTypes{SupportedResetTypes: chassis.SupportedResetTypes}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishSupportedResetTypesID, types.RedfishSupportedResetTypesLink, supportedResetTypes, supportedResetTypesChassisMask, chassis.UUID, ctx.Self().Id)
}
