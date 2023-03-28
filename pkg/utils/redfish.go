// Copyright 2023 NJWS Inc.

package utils

import (
	"context"
	"fmt"
	"net/url"

	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/stmcginnis/gofish"
)

type ConnectionParameters struct {
	Endpoint string
	Login    string
	Password string
}

func ConnectToRedfish(ctx context.Context, redfishDevice device.RedfishDevice) (client *gofish.APIClient, err error) {
	u, err := url.Parse(redfishDevice.Api)
	if err != nil {
		return
	}

	connectionParameters := ConnectionParameters{
		Endpoint: fmt.Sprintf("%s://%s", u.Scheme, u.Hostname()),
		Login:    redfishDevice.Login,
		Password: redfishDevice.Password,
	}

	return Connect(ctx, connectionParameters)
}

func Connect(ctx context.Context, connectionParameters ConnectionParameters) (client *gofish.APIClient, err error) {
	config := gofish.ClientConfig{
		Endpoint: connectionParameters.Endpoint,
		Username: connectionParameters.Login,
		Password: connectionParameters.Password,
		Insecure: true,
	}

	client, err = gofish.ConnectContext(ctx, config)
	if err != nil {
		return
	}

	return
}
