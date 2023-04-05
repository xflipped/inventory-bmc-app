// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/server"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbdiscovery"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbinventory"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbled"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"github.com/stmcginnis/gofish/common"
)

// export BMC_APP_ADDR=bmc-app
// export BMC_APP_PORT=32415
type Client struct {
	c pbbmc.BmcServiceClient
}

func New() (c *Client, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	c = &Client{
		c: pbbmc.NewBmcServiceClient(conn),
	}
	return
}

func (c *Client) Discovery(ctx context.Context, url string) (device *pbredfish.Device, err error) {
	request := &pbdiscovery.Request{
		Url: url,
	}
	return c.c.Discovery(ctx, request)
}

func (c *Client) Inventory(ctx context.Context, id, username, password string) (device *pbredfish.Device, err error) {
	request := &pbinventory.Request{
		Id:       id,
		Username: username,
		Password: password,
	}
	return c.c.Inventory(ctx, request)
}

func (c *Client) ListDevices(ctx context.Context) (devices *pbredfish.Devices, err error) {
	return c.c.ListDevices(ctx, &pbbmc.Empty{})
}

func (c *Client) SwitchLed(ctx context.Context, id, username, password string, state common.IndicatorLED) (device *pbredfish.Device, err error) {
	ledRequest := pbled.Request{
		Id:       id,
		Username: username,
		Password: password,
		State:    string(state),
	}
	return c.c.SwitchLed(ctx, &ledRequest)
}
