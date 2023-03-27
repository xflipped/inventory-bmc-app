// Copyright 2023 NJWS Inc.

package agent

import (
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/bootstrap"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
)

const (
	subServiceMask          = "%s.service.*[?@._id == '%s'?].objects.root"
	statusSubServiceMask    = "status.%s.service.*[?@._id == '%s'?].objects.root"
	subSubServiceMask       = "%s.%s.service.*[?@._id == '%s'?].objects.root"
	statusSubSubServiceMask = "status.%s.%s.service.*[?@._id == '%s'?].objects.root"
)

func (a *Agent) createOrUpdateEventService(ctx module.Context, service *gofish.Service, parentNode *documents.Node) (err error) {
	eventService, err := service.EventService()
	if err != nil {
		return
	}

	eventServiceLink := eventService.ID
	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishEventServiceID, eventServiceLink, eventService, subServiceMask, eventServiceLink, ctx.Self().Id)
	if err != nil {
		return
	}

	p := utils.NewParallel()
	p.Exec(func() error { return a.createOrUpdateEventServiceStatus(ctx, document, eventService) })
	p.Exec(func() error { return a.createOrUpdateEventDestinations(ctx, document, eventService) })
	return p.Wait()
}

func (a *Agent) createOrUpdateEventServiceStatus(ctx module.Context, parentNode *documents.Node, eventService *redfish.EventService) (err error) {
	status := &bootstrap.RedfishStatus{Status: eventService.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubServiceMask, eventService.ID, ctx.Self().Id)
}

func (a *Agent) createOrUpdateEventDestinations(ctx module.Context, parentNode *documents.Node, eventService *redfish.EventService) (err error) {
	eventDestinations, err := eventService.GetEventSubscriptions()
	if err != nil {
		return
	}

	p := utils.NewParallel()
	for _, eventDestination := range eventDestinations {
		eventDestination := eventDestination
		p.Exec(func() error { return a.createOrUpdateEventDestination(ctx, parentNode, eventService, eventDestination) })
	}
	return p.Wait()
}

func (a *Agent) createOrUpdateEventDestination(ctx module.Context, parentNode *documents.Node, eventService *redfish.EventService, eventDestination *redfish.EventDestination) (err error) {
	eventDestinationLink := eventDestination.ID
	document, err := a.syncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishEventDestinationID, eventDestinationLink, eventDestination, subSubServiceMask, eventDestinationLink, eventService.ID, ctx.Self().Id)
	if err != nil {
		return
	}

	return a.createOrUpdateEventDestinationStatus(ctx, document, eventService, eventDestinationLink, eventDestination)
}

func (a *Agent) createOrUpdateEventDestinationStatus(ctx module.Context, parentNode *documents.Node, eventService *redfish.EventService, eventDestinationLink string, eventDestination *redfish.EventDestination) (err error) {
	status := &bootstrap.RedfishStatus{Status: eventDestination.Status}
	return a.asyncCreateOrUpdateChild(ctx, parentNode.Id.String(), types.RedfishStatusID, types.RedfishStatusLink, status, statusSubSubServiceMask, eventDestinationLink, eventService.ID, ctx.Self().Id)
}
