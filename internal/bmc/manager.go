// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/db"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
	"go.mongodb.org/mongo-driver/bson"
)

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

	const colName = "managers"

	redfishManager := db.RedfishManager{
		ServiceId: redfishService.Id,
		Manager:   manager,
	}

	filter := bson.D{{"_service_id", redfishManager.ServiceId}}
	if err = b.FindOneAndReplace(ctx, colName, filter, &redfishManager); err != nil {
		return
	}

	return
}
