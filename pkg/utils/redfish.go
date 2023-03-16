// Copyright 2023 NJWS Inc.

package utils

import (
	"context"
	"fmt"
	"net/url"

	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/stmcginnis/gofish"
)

func ConnectToRedfish(ctx context.Context, redfishDevice device.RedfishDevice) (client *gofish.APIClient, err error) {
	u, err := url.Parse(redfishDevice.Api)
	if err != nil {
		return
	}

	config := gofish.ClientConfig{
		Endpoint: fmt.Sprintf("%s://%s", u.Scheme, u.Hostname()),
		Username: redfishDevice.Login,
		Password: redfishDevice.Password,
		Insecure: true,
	}

	client, err = gofish.ConnectContext(ctx, config)
	if err != nil {
		return
	}

	return client, nil
}
