// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/bootstrap"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
)

func (a *Agent) createOrUpdateChassisLogServices(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis) (err error) {
	logServices, err := chassis.LogServices()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, logService := range logServices {
		logService := logService
		p.Exec(func() error {
			return a.createOrUpdateLogService(ctx, parentNode, chassis.UUID, logService.ID, subChassisMask, logService)
		})
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateSystemLogServices(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem) (err error) {
	logServices, err := computerSystem.LogServices()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, logService := range logServices {
		logService := logService
		p.Exec(func() error {
			return a.createOrUpdateLogService(ctx, parentNode, computerSystem.UUID, logService.ID, subSystemMask, logService)
		})
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateManagerLogServices(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager) (err error) {
	logServices, err := manager.LogServices()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, logService := range logServices {
		logService := logService
		p.Exec(func() error {
			return a.createOrUpdateLogService(ctx, parentNode, manager.UUID, logService.ID, subManagerMask, logService)
		})
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateLogService(ctx module.Context, parentNode *documents.Node, parentUUID, logServiceLinkName, logServiceMask string, logService *redfish.LogService) (err error) {
	// logServiceMask template: "[log-service].[manager|chassis|system].service.*[?@._id == '%s'?].objects.root"
	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLogServiceID, logServiceLinkName, logService, logServiceMask, logServiceLinkName, parentUUID, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error {
		return a.createOrUpdateLogServiceStatus(ctx, document, parentUUID, logServiceLinkName, logServiceMask, logService)
	})
	p.Exec(func() error {
		return a.createOrUpdateLogEntries(ctx, document, parentUUID, logServiceLinkName, logServiceMask, logService)
	})
	return p.Wait()
}

func (a *Agent) createOrUpdateLogServiceStatus(ctx module.Context, parentNode *documents.Node, subServiceUUID, logServiceLinkName, logServiceMask string, logService *redfish.LogService) (err error) {
	status := &bootstrap.RedfishStatus{Status: logService.Status}
	// statusLogServiceMask template: "status.[log-service].[manager|chassis|system].service.*[?@._id == '%s'?].objects.root"
	statusLogServiceMask := "status." + logServiceMask
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusLogServiceMask, logServiceLinkName, subServiceUUID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateLogEntries(ctx module.Context, parentNode *documents.Node, subServiceUUID, logServiceLinkName, logServiceMask string, logService *redfish.LogService) (err error) {
	// TODO: check time issue
	logEntries, err := logService.Entries()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, logEntry := range logEntries {
		logEntry := logEntry
		p.Exec(func() error {
			return a.createOrUpdateLogEntry(ctx, parentNode, subServiceUUID, logServiceLinkName, logServiceMask, logEntry)
		})
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateLogEntry(ctx module.Context, parentNode *documents.Node, subServiceUUID, logServiceLinkName, logServiceMask string, logEntry *redfish.LogEntry) (err error) {
	logEntryLink := fmt.Sprintf("log-entry-%s", logEntry.ID)
	// logEntryMask template: "[log-entry].[log-service].[manager|chassis|system].service.*[?@._id == '%s'?].objects.root"
	logEntryMask := "%s." + logServiceMask

	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLogEntryID, logEntryLink, logEntry, logEntryMask, logEntryLink, logServiceLinkName, subServiceUUID, ctx.Self().Id)
}
