// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson"
)

const chasseezColName = "chasseez"

func (b *BmcApp) inventoryChasseez(ctx context.Context, redfishService db.RedfishService) (err error) {
	log.Infof("exec inventoryChasseez")

	chasseez, err := redfishService.Chassis()
	if err != nil {
		return
	}
	p := utils.NewParallel()
	for _, chassis := range chasseez {
		chassis := chassis
		p.Exec(func() error { return b.inventoryChassis(ctx, redfishService, chassis) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryChassis(ctx context.Context, redfishService db.RedfishService, chassis *redfish.Chassis) (err error) {
	log.Infof("exec inventoryChassis")

	redfishChassis := db.RedfishChassis{
		ServiceId: redfishService.Id,
		Chassis:   chassis,
	}

	filter := bson.D{{Key: "_service_id", Value: redfishChassis.ServiceId}}
	if err = b.FindOneAndReplace(ctx, chasseezColName, filter, &redfishChassis); err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return b.inventoryThermal(ctx, redfishChassis) })
	p.Exec(func() error { return b.inventoryPower(ctx, redfishChassis) })
	p.Exec(func() error { return b.inventoryDrives(ctx, redfishChassis) })
	p.Exec(func() error { return b.inventoryNetworkAdapters(ctx, redfishChassis) })

	// p.Exec(func() error { return b.inventoryChassisLogServices(ctx, redfishChassis) })

	return p.Wait()
}

func (b *BmcApp) inventoryThermal(ctx context.Context, redfishChassis db.RedfishChassis) (err error) {
	log.Infof("exec inventoryThermal")

	const colName = "thermal"

	thermal, err := redfishChassis.Thermal()
	if err != nil {
		return
	}

	redfishThermal := db.RedfishThermal{
		ChassisId: redfishChassis.Id,
		Thermal:   thermal,
	}

	filter := bson.D{{Key: "_chassis_id", Value: redfishThermal.ChassisId}}
	return b.FindOneAndReplace(ctx, colName, filter, &redfishThermal)
}

func (b *BmcApp) inventoryPower(ctx context.Context, redfishChassis db.RedfishChassis) (err error) {
	log.Infof("exec inventoryPower")

	const colName = "power"

	power, err := redfishChassis.Power()
	if err != nil {
		return
	}

	redfishPower := db.RedfishPower{
		ChassisId: redfishChassis.Id,
		Power:     power,
	}

	filter := bson.D{{Key: "_chassis_id", Value: redfishPower.ChassisId}}
	return b.FindOneAndReplace(ctx, colName, filter, &redfishPower)
}

func (b *BmcApp) inventoryDrives(ctx context.Context, redfishChassis db.RedfishChassis) (err error) {
	log.Infof("exec inventoryDrives")

	drives, err := redfishChassis.Drives()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, drive := range drives {
		drive := drive
		p.Exec(func() error { return b.inventoryDrive(ctx, redfishChassis, drive) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryDrive(ctx context.Context, redfishChassis db.RedfishChassis, drive *redfish.Drive) (err error) {
	log.Infof("exec inventoryDrive")

	const colName = "drives"

	redfishDrive := db.RedfishDrive{
		ChassisId: redfishChassis.Id,
		Drive:     drive,
	}

	filter := bson.D{{Key: "_chassis_id", Value: redfishDrive.ChassisId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishDrive); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryNetworkAdapters(ctx context.Context, redfishChassis db.RedfishChassis) (err error) {
	log.Infof("exec inventoryNetworkAdapters")

	networkAdapters, err := redfishChassis.NetworkAdapters()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, networkAdapter := range networkAdapters {
		networkAdapter := networkAdapter
		p.Exec(func() error { return b.inventoryNetworkAdapter(ctx, redfishChassis, networkAdapter) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryNetworkAdapter(ctx context.Context, redfishChassis db.RedfishChassis, networkAdapter *redfish.NetworkAdapter) (err error) {
	log.Infof("exec inventoryNetworkAdapter")

	const colName = "networkAdapters"

	redfishNetworkAdapter := db.RedfishNetworkAdapter{
		ChassisId:      redfishChassis.Id,
		NetworkAdapter: networkAdapter,
	}

	filter := bson.D{{Key: "_chassis_id", Value: redfishNetworkAdapter.ChassisId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishNetworkAdapter); err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return b.inventoryNetworkAdapterDeviceFunctions(ctx, redfishNetworkAdapter) })
	p.Exec(func() error { return b.inventoryNetworkAdapterPorts(ctx, redfishNetworkAdapter) })
	return p.Wait()
}

func (b *BmcApp) inventoryNetworkAdapterDeviceFunctions(ctx context.Context, redfishNetworkAdapter db.RedfishNetworkAdapter) (err error) {
	log.Infof("exec inventoryNetworkAdapterDeviceFunctions")

	networkAdapterDeviceFunctions, err := redfishNetworkAdapter.NetworkDeviceFunctions()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, networkAdapterDeviceFunction := range networkAdapterDeviceFunctions {
		networkAdapterDeviceFunction := networkAdapterDeviceFunction
		p.Exec(func() error {
			return b.inventoryNetworkAdapterDeviceFunction(ctx, redfishNetworkAdapter, networkAdapterDeviceFunction)
		})
	}
	return p.Wait()
}

func (b *BmcApp) inventoryNetworkAdapterDeviceFunction(ctx context.Context, redfishNetworkAdapter db.RedfishNetworkAdapter, networkAdapterDeviceFunction *redfish.NetworkDeviceFunction) (err error) {
	log.Infof("exec inventoryNetworkAdapterDeviceFunction")

	const colName = "networkAdapterDeviceFunctions"

	redfishNetworkAdapterDeviceFunction := db.RedfishNetworkAdapterDeviceFunction{
		NetworkAdapterId:      redfishNetworkAdapter.Id,
		NetworkDeviceFunction: networkAdapterDeviceFunction,
	}

	filter := bson.D{{Key: "_network_adapter_id", Value: redfishNetworkAdapterDeviceFunction.NetworkAdapterId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishNetworkAdapterDeviceFunction); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryNetworkAdapterPorts(ctx context.Context, redfishNetworkAdapter db.RedfishNetworkAdapter) (err error) {
	log.Infof("exec inventoryNetworkAdapterPorts")

	networkAdapterPorts, err := redfishNetworkAdapter.NetworkPorts()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, networkAdapterPort := range networkAdapterPorts {
		networkAdapterPort := networkAdapterPort
		p.Exec(func() error {
			return b.inventoryNetworkAdapterPort(ctx, redfishNetworkAdapter, networkAdapterPort)
		})
	}
	return p.Wait()
}

func (b *BmcApp) inventoryNetworkAdapterPort(ctx context.Context, redfishNetworkAdapter db.RedfishNetworkAdapter, networkAdapterPort *redfish.NetworkPort) (err error) {
	log.Infof("exec inventoryNetworkAdapterPort")

	const colName = "networkAdapterPorts"

	redfishNetworkAdapterPort := db.RedfishNetworkAdapterPort{
		NetworkAdapterId: redfishNetworkAdapter.Id,
		NetworkPort:      networkAdapterPort,
	}

	filter := bson.D{{Key: "_network_adapter_id", Value: redfishNetworkAdapterPort.NetworkAdapterId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishNetworkAdapterPort); err != nil {
		return
	}

	return
}
