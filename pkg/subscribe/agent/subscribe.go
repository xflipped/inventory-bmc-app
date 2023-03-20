// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"

	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
)

// subscribeFunction executes on '[device-uuid].redfish-devices.root'
func (a *Agent) subscribeFunction(ctx module.Context) (err error) {
	var redfishDevice device.RedfishDevice
	if err = json.Unmarshal(ctx.CmdbContext(), &redfishDevice); err != nil {
		return
	}
	var destinationUrl string
	if err = json.Unmarshal(ctx.Message(), &destinationUrl); err != nil {
		return
	}

	client, err := utils.ConnectToRedfish(ctx, redfishDevice)
	if err != nil {
		return err
	}
	defer client.Logout()

	return a.subscribeToBmcEvents(client.GetService(), destinationUrl)
}

func (a *Agent) subscribeToBmcEvents(service *gofish.Service, destinationUrl string) (err error) {
	eventService, err := service.EventService()
	if err != nil {
		return
	}

	// If RegistryPrefixes and ResourceTypes are empty on subscription,
	// the client is subscribing to all available Message Registries and Resource Types
	subscriptionLink, err := eventService.CreateEventSubscriptionInstance(
		destinationUrl,
		nil,
		nil,
		nil,
		redfish.RedfishEventDestinationProtocol,
		"Public",
		redfish.RetryForeverDeliveryRetryPolicy,
		nil,
	)
	if err == nil {
		log.Infof("subscription link: %s", subscriptionLink)
	}

	// TODO: subscription link re-inventory
	return
}
