// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"
	"net/url"
	"path"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish/redfish"
)

type SubscribePayload struct {
	DestinationUrl       string
	RegistryPrefixes     []string
	ResourceTypes        []string
	ConnectionParameters utils.ConnectionParameters
}

const subscriptionMask = "%s.*[?@._id == '%s'?].objects.root"

// subscribeFunction executes on '[event-service].[device-uuid].redfish-devices.root'
func (a *Agent) subscribeFunction(ctx module.Context) (err error) {
	eventService := &redfish.EventService{}
	if err = json.Unmarshal(ctx.CmdbContext(), eventService); err != nil {
		return
	}
	var payload SubscribePayload
	if err = json.Unmarshal(ctx.Message(), &payload); err != nil {
		return
	}

	client, err := utils.Connect(ctx, payload.ConnectionParameters)
	if err != nil {
		return err
	}
	defer client.Logout()

	eventService, err = redfish.GetEventService(client, eventService.ODataID)
	if err != nil {
		return
	}

	return a.subcribeToBmcEvents(ctx, eventService, payload)
}

func (a *Agent) subcribeToBmcEvents(ctx module.Context, eventService *redfish.EventService, payload SubscribePayload) (err error) {
	// If RegistryPrefixes and ResourceTypes are empty on subscription,
	// the client is subscribing to all available Message Registries and Resource Types
	subscriptionLink, err := eventService.CreateEventSubscriptionInstance(
		payload.DestinationUrl,
		payload.RegistryPrefixes,
		payload.ResourceTypes,
		nil,
		redfish.RedfishEventDestinationProtocol,
		"Public",
		redfish.RetryForeverDeliveryRetryPolicy,
		nil,
	)
	if err != nil {
		return
	}

	log.Infof("subscription link: %s", subscriptionLink)

	subscriptionPayload, err := eventService.GetEventSubscription(subscriptionLink)
	if err != nil {
		return
	}

	subscriptionUrl, err := url.Parse(subscriptionLink)
	if err != nil {
		return
	}
	subscriptionLink = path.Base(subscriptionUrl.Path)

	return a.asyncCreateChild(ctx, types.RedfishEventDestinationID, subscriptionLink, subscriptionPayload, subscriptionMask, subscriptionLink, ctx.Self().Id)
}

func (a *Agent) createSyncCreateOrUpdateChild(from, moType, name string, payload any, format string, args ...any) (functionContext *pbtypes.FunctionContext, err error) {
	document, err := a.getDocument(format, args...)
	if err != nil {
		return system.CreateChild(from, moType, name, payload)
	}

	return system.UpdateObject(document.Id.String(), payload)
}

func (a *Agent) asyncCreateChild(ctx module.Context, moType, name string, payload any, format string, args ...any) (err error) {
	functionContext, err := a.createSyncCreateOrUpdateChild(ctx.Self().Id, moType, name, payload, format, args...)
	if err != nil {
		return
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)
	return
}
