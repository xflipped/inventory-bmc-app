// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

func (a *Agent) createOrUpdateChassisSELLogService(ctx module.Context, parentNode *documents.Node, chassis *redfish.Chassis, vendorData *VendorSpecificData) (err error) {
	return a.createOrUpdateLogService(ctx, parentNode, chassis.UUID, types.RedfishSELLogsLink, subChassisMask, chassis.Client, vendorData.ChassisSELLogs)
}

func (a *Agent) createOrUpdateSystemSELLogService(ctx module.Context, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, vendorData *VendorSpecificData) (err error) {
	return a.createOrUpdateLogService(ctx, parentNode, computerSystem.UUID, types.RedfishSELLogsLink, subSystemMask, computerSystem.Client, vendorData.SystemSELLogs)
}

func (a *Agent) createOrUpdateManagerSELLogService(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager, vendorData *VendorSpecificData) (err error) {
	return a.createOrUpdateLogService(ctx, parentNode, manager.UUID, types.RedfishSELLogsLink, subManagerMask, manager.Client, vendorData.ManagerSELLogs)
}

func (a *Agent) createOrUpdateManagerEventLogService(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager, vendorData *VendorSpecificData) (err error) {
	return a.createOrUpdateLogService(ctx, parentNode, manager.UUID, types.RedfishEventLogsLink, subManagerMask, manager.Client, vendorData.ManagerEventLogs)
}

func (a *Agent) createOrUpdateManagerAuditLogService(ctx module.Context, parentNode *documents.Node, manager *redfish.Manager, vendorData *VendorSpecificData) (err error) {
	return a.createOrUpdateLogService(ctx, parentNode, manager.UUID, types.RedfishAuditLogsLink, subManagerMask, manager.Client, vendorData.ManagerAuditLogs)
}

func (a *Agent) createOrUpdateLogService(ctx module.Context, parentNode *documents.Node, parentUUID, logServiceLinkName, logServiceMask string, client common.Client, redfishLogLink string) (err error) {
	logService, err := redfish.GetLogService(client, redfishLogLink)
	if err != nil {
		return
	}

	// logServiceMask template: "[log-service].[manager|chassis|system].service.*[?@._id == '%s'?].objects.root"
	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishLogServiceID, logServiceLinkName, logService, logServiceMask, logServiceLinkName, parentUUID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateLogEntries(ctx, document, parentUUID, logServiceLinkName, logServiceMask, logService)
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
