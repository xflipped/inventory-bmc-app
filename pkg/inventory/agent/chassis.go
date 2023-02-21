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
	chassisMask                              = "chassis-%s.service.*[?@._id == '%s'?].objects.root"
	ledChassisMask                           = "led.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	statusChassisMask                        = "status.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	powerStateChassisMask                    = "power-state.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	physicalSecurityChassisMask              = "physical-security.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	locationChassisMask                      = "location.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	partLocationLocationChassisMask          = "part-location.location.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	placementLocationChassisMask             = "placement.location.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	postalAddressLocationChassisMask         = "postal-address.location.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	supportedResetTypesLocationChassisMask   = "supported-reset-types.location.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	thermalMask                              = "thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	subThermalMask                           = "%s.thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	statusSubThermalMask                     = "status.%s.thermal.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	powerChassisMask                         = "power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	subPowerChassisMask                      = "%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	ledSubPowerChassisMask                   = "led.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	locationSubPowerChassisMask              = "location.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	partLocationLocationSubPowerChassisMask  = "part-location.location.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	placementLocationSubPowerChassisMask     = "placement.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	postalAddressLocationSubPowerChassisMask = "postal-address.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	statusSubPowerChassisMask                = "status.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	physicalContextSubPowerChassisMask       = "physcial-context.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	powerMetricSubPowerChassisMask           = "power-metric.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	powerLimitSubPowerChassisMask            = "power-limit.%s.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
	ledPowerChassisMask                      = "led.power.chassis-%s.service.*[?@._id == '%s'?].objects.root"
)

func (a *Agent) createOrUpdateChasseez(ctx module.Context, service *gofish.Service, parentNode *documents.Node) (err error) {
	chasseez, err := service.Chassis()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, chassee := range chasseez {
		chassee := chassee
		p.Exec(func() error { return a.createOrUpdateChassee(ctx, parentNode, chassee) })
	}
	return p.Wait()
}

// TODO: check Chassis & RedfishDevice UUID, now they are the same
func (a *Agent) createOrUpdateChassee(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	chassisLink := fmt.Sprintf("chassis-%s", chassis.UUID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishChassisID, chassisLink, chassis, chassisMask, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return a.createOrUpdateThermal(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdatePower(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisLed(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisStatus(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisPowerState(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisPhysicalSecurity(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisLocation(ctx, document, chassis) })
	p.Exec(func() error { return a.createOrUpdateChassisSupportedResetTypes(ctx, document, chassis) })
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

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishThermalID, types.RedfishThermalLink, thermal, thermalMask, chassis.UUID, ctx.Self().Id)
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
	temperatureLink := temperature.MemberID

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishTemperatureID, temperatureLink, temperature, subThermalMask, temperatureLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateThermalTemperatureStatus(ctx, document, chassis, temperatureLink, temperature)
}

func (a *Agent) createOrUpdateThermalTemperatureStatus(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, temperatureLink string, temperature redfish.Temperature) (err error) {
	status := &bootstrap.RedfishStatus{Status: temperature.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubThermalMask, temperatureLink, chassis.UUID, ctx.Self().Id)
}

// TODO: update later, currently not available
func (a *Agent) createOrUpdateThermalFan(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, fan redfish.Fan) (err error) {
	fanLink := fan.MemberID
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishFanID, fanLink, fan, subThermalMask, fanLink, chassis.UUID, ctx.Self().Id)
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
		return a.createOrUpdatePowerIndicatorLED(ctx, chassis, document, redfishPower.IndicatorLED)
	})
	for _, powerControl := range redfishPower.PowerControl {
		powerControl := powerControl
		p.Exec(func() error { return a.createOrUpdatePowerControl(ctx, chassis, document, powerControl) })
	}
	for _, powerSupply := range redfishPower.PowerSupplies {
		powerSupply := powerSupply
		p.Exec(func() error { return a.createOrUpdatePowerSupply(ctx, chassis, document, powerSupply) })
	}
	for _, voltage := range redfishPower.Voltages {
		voltage := voltage
		p.Exec(func() error { return a.createOrUpdateVoltage(ctx, chassis, document, voltage) })
	}
	return p.Wait()
}

func (a *Agent) createOrUpdatePowerIndicatorLED(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, indicatorLED common.IndicatorLED) (err error) {
	led := &bootstrap.RedfishLed{Led: indicatorLED}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLedID, types.RedfishLedLink, led, ledPowerChassisMask, chassis.UUID, ctx.Self().Id)
}

// TODO: unique link
// TODO: device with Id is present ? create : ignore
func (a *Agent) createOrUpdatePowerControl(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerControl redfish.PowerControl) (err error) {
	powerControlLink := fmt.Sprintf("power-control-%s", powerControl.MemberID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerControlID, powerControlLink, powerControl, subPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdatePowerControlStatus(ctx, chassis, document, powerControlLink, powerControl)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerControlPhysicalContext(ctx, chassis, document, powerControlLink, powerControl)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerControlMetric(ctx, chassis, document, powerControlLink, powerControl)
	})

	p.Exec(func() error {
		return a.createOrUpdatePowerControlLimit(ctx, chassis, document, powerControlLink, powerControl)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdatePowerControlStatus(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerControlLink string, powerControl redfish.PowerControl) (err error) {
	status := &bootstrap.RedfishStatus{Status: powerControl.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerControlPhysicalContext(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerControlLink string, powerControl redfish.PowerControl) (err error) {
	physicalContext := &bootstrap.RedfishPhysicalContext{PhysicalContext: &powerControl.PhysicalContext}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPhysicalContextID, types.RedfishPhysicalContextLink, physicalContext, physicalContextSubPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerControlMetric(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerControlLink string, powerControl redfish.PowerControl) (err error) {
	powerMetric := powerControl.PowerMetrics
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerMetricID, types.RedfishPowerMetricLink, powerMetric, powerMetricSubPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerControlLimit(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerControlLink string, powerControl redfish.PowerControl) (err error) {
	powerLimit := powerControl.PowerLimit
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerLimitID, types.RedfishPowerLimitLink, powerLimit, powerLimitSubPowerChassisMask, powerControlLink, chassis.UUID, ctx.Self().Id)
}

// TODO: unique link
// TODO: device with Id is present ? create : ignore
func (a *Agent) createOrUpdatePowerSupply(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerSupply redfish.PowerSupply) (err error) {
	powerSupplyLink := fmt.Sprintf("power-supply-%s", powerSupply.MemberID)

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPowerSupplyID, powerSupplyLink, powerSupply, subPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyStatus(ctx, chassis, document, powerSupplyLink, powerSupply)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyIndicatorLED(ctx, chassis, document, powerSupplyLink, powerSupply)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyLocation(ctx, chassis, document, powerSupplyLink, powerSupply)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdatePowerSupplyStatus(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, powerSupply redfish.PowerSupply) (err error) {
	status := &bootstrap.RedfishStatus{Status: powerSupply.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerSupplyIndicatorLED(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, powerSupply redfish.PowerSupply) (err error) {
	led := &bootstrap.RedfishLed{Led: powerSupply.IndicatorLED}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLedID, types.RedfishLedLink, led, ledSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerSupplyLocation(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, powerSupply redfish.PowerSupply) (err error) {
	location := powerSupply.Location

	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLocationID, types.RedfishLocationLink, location, locationSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyPartLocation(ctx, chassis, document, powerSupplyLink, location.PartLocation)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyPlacement(ctx, chassis, document, powerSupplyLink, location.Placement)
	})
	p.Exec(func() error {
		return a.createOrUpdatePowerSupplyPostalAddress(ctx, chassis, document, powerSupplyLink, location.PostalAddress)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdatePowerSupplyPartLocation(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, partLocation common.PartLocation) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPartLocationID, types.RedfishPartLocationLink, partLocation, partLocationLocationSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerSupplyPlacement(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, placement common.Placement) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPlacementID, types.RedfishPlacementLink, placement, placementLocationSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePowerSupplyPostalAddress(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, powerSupplyLink string, postalAddress common.PostalAddress) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPostalAddressID, types.RedfishPostalAddressLink, postalAddress, postalAddressLocationSubPowerChassisMask, powerSupplyLink, chassis.UUID, ctx.Self().Id)
}

// TODO: unique link
// TODO: device with Id is present ? create : ignore
func (a *Agent) createOrUpdateVoltage(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, voltage redfish.Voltage) (err error) {
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
	p.Exec(func() error { return a.createOrUpdatePartLocation(ctx, chassis, document, location.PartLocation) })
	p.Exec(func() error { return a.createOrUpdatePlacement(ctx, chassis, document, location.Placement) })
	p.Exec(func() error { return a.createOrUpdatePostalAddress(ctx, chassis, document, location.PostalAddress) })
	return p.Wait()
}

func (a *Agent) createOrUpdatePartLocation(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, partLocation common.PartLocation) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPartLocationID, types.RedfishPartLocationLink, partLocation, partLocationLocationChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePlacement(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, placement common.Placement) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPlacementID, types.RedfishPlacementLink, placement, placementLocationChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdatePostalAddress(ctx module.Context, chassis *redfish.Chassis, parentNode *documents.Node, postalAddress common.PostalAddress) (err error) {
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishPostalAddressID, types.RedfishPostalAddressLink, postalAddress, postalAddressLocationChassisMask, chassis.UUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateChassisSupportedResetTypes(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	supportedResetTypes := &bootstrap.RedfishSupportedResetTypes{SupportedResetTypes: chassis.SupportedResetTypes}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishSupportedResetTypesID, types.RedfishSupportedResetTypesLink, supportedResetTypes, supportedResetTypesLocationChassisMask, chassis.UUID, ctx.Self().Id)
}
