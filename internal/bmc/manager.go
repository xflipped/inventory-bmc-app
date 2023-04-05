// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson"
)

const managersColName = "managers"

func (b *BmcApp) inventoryManagers(ctx context.Context, redfishService db.RedfishService) (err error) {
	log.Infof("exec inventoryManagers")

	managers, err := redfishService.Managers()
	if err != nil {
		return
	}
	p := utils.NewParallel()
	for _, manager := range managers {
		manager := manager
		p.Exec(func() error { return b.inventoryManager(ctx, redfishService, manager) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryManager(ctx context.Context, redfishService db.RedfishService, manager *redfish.Manager) (err error) {
	log.Infof("exec inventoryManager")

	redfishManager := db.RedfishManager{
		ServiceId: redfishService.Id,
		Manager:   manager,
	}

	filter := bson.D{{Key: "_service_id", Value: redfishManager.ServiceId}}
	if err = b.FindOneAndReplace(ctx, managersColName, filter, &redfishManager); err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return b.inventoryManagerEthernetInterfaces(ctx, redfishManager) })
	p.Exec(func() error { return b.inventoryHostInterfaces(ctx, redfishManager) })
	p.Exec(func() error { return b.inventoryVirtualMedias(ctx, redfishManager) })

	// p.Exec(func() error { return b.inventoryManagerLogServices(ctx, redfishManager) })

	return p.Wait()
}

func (b *BmcApp) inventoryManagerEthernetInterfaces(ctx context.Context, redfishManager db.RedfishManager) (err error) {
	log.Infof("exec inventoryManagerEthernetInterfaces")

	ethernetInterfaces, err := redfishManager.EthernetInterfaces()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, ethernetInterface := range ethernetInterfaces {
		ethernetInterface := ethernetInterface
		p.Exec(func() error { return b.inventoryManagerEthernetInterface(ctx, redfishManager, ethernetInterface) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryManagerEthernetInterface(ctx context.Context, redfishManager db.RedfishManager, ethernetInterface *redfish.EthernetInterface) (err error) {
	log.Infof("exec inventoryManagerEthernetInterface")

	const colName = "managerEthernetInterfaces"

	redfishManagerEthernetInterface := db.RedfishManagerEthernetInterface{
		ManagerId:         redfishManager.Id,
		EthernetInterface: ethernetInterface,
	}

	filter := bson.D{{Key: "_manager_id", Value: redfishManagerEthernetInterface.ManagerId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishManagerEthernetInterface); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryHostInterfaces(ctx context.Context, redfishManager db.RedfishManager) (err error) {
	log.Infof("exec inventoryHostInterfaces")

	hostInterfaces, err := redfishManager.HostInterfaces()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, hostInterface := range hostInterfaces {
		hostInterface := hostInterface
		p.Exec(func() error { return b.inventoryHostInterface(ctx, redfishManager, hostInterface) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryHostInterface(ctx context.Context, redfishManager db.RedfishManager, hostInterface *redfish.HostInterface) (err error) {
	log.Infof("exec inventoryHostInterface")

	const colName = "hostInterfaces"

	redfishHostInterface := db.RedfishHostInterface{
		ManagerId:     redfishManager.Id,
		HostInterface: hostInterface,
	}

	filter := bson.D{{Key: "_manager_id", Value: redfishHostInterface.ManagerId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishHostInterface); err != nil {
		return
	}

	return
}

func (b *BmcApp) inventoryVirtualMedias(ctx context.Context, redfishManager db.RedfishManager) (err error) {
	log.Infof("exec inventoryVirtualMedias")

	virtualMedias, err := redfishManager.VirtualMedia()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, virtualMedia := range virtualMedias {
		virtualMedia := virtualMedia
		p.Exec(func() error { return b.inventoryVirtualMedia(ctx, redfishManager, virtualMedia) })
	}
	return p.Wait()
}

func (b *BmcApp) inventoryVirtualMedia(ctx context.Context, redfishManager db.RedfishManager, virtualMedia *redfish.VirtualMedia) (err error) {
	log.Infof("exec inventoryVirtualMedia")

	const colName = "virtualMedia"

	redfishVirtualMedia := db.RedfishVirtualMedia{
		ManagerId:    redfishManager.Id,
		VirtualMedia: virtualMedia,
	}

	filter := bson.D{{Key: "_manager_id", Value: redfishVirtualMedia.ManagerId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishVirtualMedia); err != nil {
		return
	}

	return
}
