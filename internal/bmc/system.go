// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson"
)

const systemsColName = "systems"

func (b *BmcApp) inventorySystems(ctx context.Context, redfishService db.RedfishService) (err error) {
	log.Infof("exec inventorySystems")

	systems, err := redfishService.Systems()
	if err != nil {
		return
	}
	p := utils.NewParallel()
	for _, system := range systems {
		system := system
		p.Exec(func() error { return b.inventorySystem(ctx, redfishService, system) })
	}
	err = p.Wait()
	return
}

func (b *BmcApp) inventorySystem(ctx context.Context, redfishService db.RedfishService, computerSystem *redfish.ComputerSystem) (err error) {
	log.Infof("exec inventorySystem")

	redfishSystem := db.RedfishSystem{
		ServiceId:      redfishService.Id,
		ComputerSystem: computerSystem,
	}

	filter := bson.D{{Key: "_service_id", Value: redfishSystem.ServiceId}}
	if err = b.FindOneAndReplace(ctx, systemsColName, filter, &redfishSystem); err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return b.inventoryBIOS(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventoryBootOptions(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventorySecureBoot(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventoryPCIeDevices(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventoryPCIeFunctions(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventoryProcessors(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventoryMemories(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventoryMemoryDomains(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventorySimpleStorages(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventoryStorages(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventoryNetworkInterfaces(ctx, redfishSystem) })
	p.Exec(func() error { return b.inventorySystemEthernetInterfaces(ctx, redfishSystem) })

	// p.Exec(func() error { return b.inventorySystemLogServices(ctx, redfishSystem) })

	return p.Wait()
}

func (b *BmcApp) inventoryBIOS(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryBIOS")

	const colName = "bios"

	bios, err := redfishSystem.Bios()
	if err != nil {
		return
	}

	redfishBIOS := db.RedfishBIOS{
		SystemId: redfishSystem.Id,
		Bios:     bios,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishBIOS.SystemId}}
	return b.FindOneAndReplace(ctx, colName, filter, &redfishBIOS)
}

func (b *BmcApp) inventoryBootOptions(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryBootOptions")

	bootOptions, err := redfishSystem.BootOptions()
	if err != nil {
		return
	}
	p := utils.NewParallel()
	for _, bootOption := range bootOptions {
		bootOption := bootOption
		p.Exec(func() error { return b.inventoryBootOption(ctx, redfishSystem, bootOption) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryBootOption(ctx context.Context, redfishSystem db.RedfishSystem, bootOption *redfish.BootOption) (err error) {
	log.Infof("exec inventoryBootOption")

	const colName = "bootOptions"

	redfishBootOption := db.RedfishBootOption{
		SystemId:   redfishSystem.Id,
		BootOption: bootOption,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishBootOption.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishBootOption); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventorySecureBoot(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventorySecureBoot")

	const colName = "secureBoot"

	secureBoot, err := redfishSystem.SecureBoot()
	if err != nil {
		return
	}

	redfishSecureBoot := db.RedfishSecureBoot{
		SystemId:   redfishSystem.Id,
		SecureBoot: secureBoot,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishSecureBoot.SystemId}}
	return b.FindOneAndReplace(ctx, colName, filter, &redfishSecureBoot)
}

func (b *BmcApp) inventoryPCIeDevices(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryPCIeDevices")

	pcieDevices, err := redfishSystem.PCIeDevices()
	if err != nil {
		return
	}
	p := utils.NewParallel()
	for _, pcieDevice := range pcieDevices {
		pcieDevice := pcieDevice
		p.Exec(func() error { return b.inventoryPCIeDevice(ctx, redfishSystem, pcieDevice) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryPCIeDevice(ctx context.Context, redfishSystem db.RedfishSystem, pcieDevice *redfish.PCIeDevice) (err error) {
	log.Infof("exec inventoryPCIeDevice")

	const colName = "pcieDevices"

	redfishPCIeDevice := db.RedfishPCIeDevice{
		SystemId:   redfishSystem.Id,
		PCIeDevice: pcieDevice,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishPCIeDevice.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishPCIeDevice); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryPCIeFunctions(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryPCIeFunctions")

	pcieFunctions, err := redfishSystem.PCIeFunctions()
	if err != nil {
		return
	}
	p := utils.NewParallel()
	for _, pcieFunction := range pcieFunctions {
		pcieFunction := pcieFunction
		p.Exec(func() error { return b.inventoryPCIeFunction(ctx, redfishSystem, pcieFunction) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryPCIeFunction(ctx context.Context, redfishSystem db.RedfishSystem, pcieFunction *redfish.PCIeFunction) (err error) {
	log.Infof("exec inventoryPCIeFunction")

	const colName = "pcieFunctions"

	redfishPCIeFunction := db.RedfishPCIeFunction{
		SystemId:     redfishSystem.Id,
		PCIeFunction: pcieFunction,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishPCIeFunction.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishPCIeFunction); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryProcessors(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryProcessors")

	processors, err := redfishSystem.Processors()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, processor := range processors {
		processor := processor
		p.Exec(func() error { return b.inventoryProcessor(ctx, redfishSystem, processor) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryProcessor(ctx context.Context, redfishSystem db.RedfishSystem, processor *redfish.Processor) (err error) {
	log.Infof("exec inventoryProcessor")

	const colName = "processors"

	redfishProcessor := db.RedfishProcessor{
		SystemId:  redfishSystem.Id,
		Processor: processor,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishProcessor.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishProcessor); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryMemories(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryMemories")

	memories, err := redfishSystem.Memory()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, memory := range memories {
		memory := memory
		p.Exec(func() error { return b.inventoryMemory(ctx, redfishSystem, memory) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryMemory(ctx context.Context, redfishSystem db.RedfishSystem, memory *redfish.Memory) (err error) {
	log.Infof("exec inventoryMemory")

	const colName = "memory"

	redfishMemory := db.RedfishMemory{
		SystemId: redfishSystem.Id,
		Memory:   memory,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishMemory.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishMemory); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryMemoryDomains(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryMemoryDomains")

	memoryDomains, err := redfishSystem.MemoryDomains()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, memoryDomain := range memoryDomains {
		memoryDomain := memoryDomain
		p.Exec(func() error { return b.inventoryMemoryDomain(ctx, redfishSystem, memoryDomain) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryMemoryDomain(ctx context.Context, redfishSystem db.RedfishSystem, memoryDomain *redfish.MemoryDomain) (err error) {
	log.Infof("exec inventoryMemoryDomain")

	const colName = "memoryDomains"

	redfishMemoryDomain := db.RedfishMemoryDomain{
		SystemId:     redfishSystem.Id,
		MemoryDomain: memoryDomain,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishMemoryDomain.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishMemoryDomain); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventorySimpleStorages(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventorySimpleStorages")

	simpleStorages, err := redfishSystem.SimpleStorages()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, simpleStorage := range simpleStorages {
		simpleStorage := simpleStorage
		p.Exec(func() error { return b.inventorySimpleStorage(ctx, redfishSystem, simpleStorage) })
	}
	return p.Wait()
}

func (b *BmcApp) inventorySimpleStorage(ctx context.Context, redfishSystem db.RedfishSystem, simpleStorage *redfish.SimpleStorage) (err error) {
	log.Infof("exec inventorySimpleStorage")

	const colName = "simpleStorage"

	redfishSimpleStorage := db.RedfishSimpleStorage{
		SystemId:      redfishSystem.Id,
		SimpleStorage: simpleStorage,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishSimpleStorage.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishSimpleStorage); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryStorages(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryStorages")

	storages, err := redfishSystem.Storage()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, storage := range storages {
		storage := storage
		p.Exec(func() error { return b.inventoryStorage(ctx, redfishSystem, storage) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryStorage(ctx context.Context, redfishSystem db.RedfishSystem, storage *redfish.Storage) (err error) {
	log.Infof("exec inventorySimpleStorage")

	const colName = "simpleStorage"

	redfishStorage := db.RedfishStorage{
		SystemId: redfishSystem.Id,
		Storage:  storage,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishStorage.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishStorage); err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return b.inventoryStorageDrives(ctx, redfishStorage) })
	p.Exec(func() error { return b.inventoryStorageVolumes(ctx, redfishStorage) })
	return p.Wait()
}

func (b *BmcApp) inventoryStorageDrives(ctx context.Context, redfishStorage db.RedfishStorage) (err error) {
	log.Infof("exec inventoryStorageDrives")

	drives, err := redfishStorage.Drives()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, drive := range drives {
		drive := drive
		p.Exec(func() error { return b.inventoryStorageDrive(ctx, redfishStorage, drive) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryStorageDrive(ctx context.Context, redfishStorage db.RedfishStorage, drive *redfish.Drive) (err error) {
	log.Infof("exec inventoryStorageDrive")

	const colName = "storageDrives"

	redfishStorageDrive := db.RedfishStorageDrive{
		StorageId: redfishStorage.Id,
		Drive:     drive,
	}

	filter := bson.D{{Key: "_storage_id", Value: redfishStorageDrive.StorageId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishStorageDrive); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryStorageVolumes(ctx context.Context, redfishStorage db.RedfishStorage) (err error) {
	log.Infof("exec inventoryStorageVolumes")

	volumes, err := redfishStorage.Volumes()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, volume := range volumes {
		volume := volume
		p.Exec(func() error { return b.inventoryStorageVolume(ctx, redfishStorage, volume) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryStorageVolume(ctx context.Context, redfishStorage db.RedfishStorage, volume *redfish.Volume) (err error) {
	log.Infof("exec inventoryStorageVolume")

	const colName = "storageVolumes"

	redfishStorageVolume := db.RedfishStorageVolume{
		StorageId: redfishStorage.Id,
		Volume:    volume,
	}

	filter := bson.D{{Key: "_storage_id", Value: redfishStorageVolume.StorageId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishStorageVolume); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryNetworkInterfaces(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventoryNetworkInterfaces")

	networkInterfaces, err := redfishSystem.NetworkInterfaces()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, networkInterface := range networkInterfaces {
		networkInterface := networkInterface
		p.Exec(func() error { return b.inventoryNetworkInterface(ctx, redfishSystem, networkInterface) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryNetworkInterface(ctx context.Context, redfishSystem db.RedfishSystem, networkInterface *redfish.NetworkInterface) (err error) {
	log.Infof("exec inventoryNetworkInterface")

	const colName = "networkInterfaces"

	redfishNetworkInterface := db.RedfishNetworkInterface{
		SystemId:         redfishSystem.Id,
		NetworkInterface: networkInterface,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishNetworkInterface.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishNetworkInterface); err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return b.inventoryNetworkInterfaceAdapter(ctx, redfishNetworkInterface) })
	p.Exec(func() error { return b.inventoryNetworkInterfaceDeviceFunctions(ctx, redfishNetworkInterface) })
	p.Exec(func() error { return b.inventoryNetworkInterfacePorts(ctx, redfishNetworkInterface) })
	return p.Wait()
}

func (b *BmcApp) inventoryNetworkInterfaceAdapter(ctx context.Context, redfishNetworkInterface db.RedfishNetworkInterface) (err error) {
	log.Infof("exec inventoryNetworkInterfaceAdapter")

	const colName = "networkInterfaceAdapters"

	networkInterfaceAdapter, err := redfishNetworkInterface.NetworkAdapter()
	if err != nil {
		return
	}

	redfishNetworkInterfaceAdapter := db.RedfishNetworkInterfaceAdapter{
		NetworkInterfaceId: redfishNetworkInterface.Id,
		NetworkAdapter:     networkInterfaceAdapter,
	}

	filter := bson.D{{Key: "_network_interface_id", Value: redfishNetworkInterfaceAdapter.NetworkInterfaceId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishNetworkInterfaceAdapter); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryNetworkInterfaceDeviceFunctions(ctx context.Context, redfishNetworkInterface db.RedfishNetworkInterface) (err error) {
	log.Infof("exec inventoryNetworkInterfaceDeviceFunctions")

	networkInterfaceDeviceFunctions, err := redfishNetworkInterface.NetworkDeviceFunctions()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, networkInterfaceDeviceFunction := range networkInterfaceDeviceFunctions {
		networkInterfaceDeviceFunction := networkInterfaceDeviceFunction
		p.Exec(func() error {
			return b.inventoryNetworkInterfaceDeviceFunction(ctx, redfishNetworkInterface, networkInterfaceDeviceFunction)
		})
	}
	return p.Wait()
}

func (b *BmcApp) inventoryNetworkInterfaceDeviceFunction(ctx context.Context, redfishNetworkInterface db.RedfishNetworkInterface, networkInterfaceDeviceFunction *redfish.NetworkDeviceFunction) (err error) {
	log.Infof("exec inventoryNetworkInterfaceDeviceFunction")

	const colName = "networkInterfaceDeviceFunctions"

	redfishNetworkInterfaceDeviceFunction := db.RedfishNetworkInterfaceDeviceFunction{
		NetworkInterfaceId:    redfishNetworkInterface.Id,
		NetworkDeviceFunction: networkInterfaceDeviceFunction,
	}

	filter := bson.D{{Key: "_network_interface_id", Value: redfishNetworkInterfaceDeviceFunction.NetworkInterfaceId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishNetworkInterfaceDeviceFunction); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryNetworkInterfacePorts(ctx context.Context, redfishNetworkInterface db.RedfishNetworkInterface) (err error) {
	log.Infof("exec inventoryNetworkInterfacePorts")

	networkInterfacePorts, err := redfishNetworkInterface.NetworkPorts()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, networkInterfacePort := range networkInterfacePorts {
		networkInterfacePort := networkInterfacePort
		p.Exec(func() error {
			return b.inventoryNetworkInterfacePort(ctx, redfishNetworkInterface, networkInterfacePort)
		})
	}
	return p.Wait()
}

func (b *BmcApp) inventoryNetworkInterfacePort(ctx context.Context, redfishNetworkInterface db.RedfishNetworkInterface, networkInterfacePort *redfish.NetworkPort) (err error) {
	log.Infof("exec inventoryNetworkInterfacePort")

	const colName = "networkInterfacePorts"

	redfishNetworkInterfacePort := db.RedfishNetworkInterfacePort{
		NetworkInterfaceId: redfishNetworkInterface.Id,
		NetworkPort:        networkInterfacePort,
	}

	filter := bson.D{{Key: "_network_interface_id", Value: redfishNetworkInterfacePort.NetworkInterfaceId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishNetworkInterfacePort); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventorySystemEthernetInterfaces(ctx context.Context, redfishSystem db.RedfishSystem) (err error) {
	log.Infof("exec inventorySystemEthernetInterfaces")

	ethernetInterfaces, err := redfishSystem.EthernetInterfaces()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, ethernetInterface := range ethernetInterfaces {
		ethernetInterface := ethernetInterface
		p.Exec(func() error { return b.inventorySystemEthernetInterface(ctx, redfishSystem, ethernetInterface) })
	}
	return p.Wait()
}

func (b *BmcApp) inventorySystemEthernetInterface(ctx context.Context, redfishSystem db.RedfishSystem, ethernetInterface *redfish.EthernetInterface) (err error) {
	log.Infof("exec inventorySystemEthernetInterface")

	const colName = "systemEthernetInterfaces"

	redfishSystemEthernetInterface := db.RedfishSystemEthernetInterface{
		SystemId:          redfishSystem.Id,
		EthernetInterface: ethernetInterface,
	}

	filter := bson.D{{Key: "_system_id", Value: redfishSystemEthernetInterface.SystemId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishSystemEthernetInterface); err != nil {
		return
	}

	return
}
