// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/internal/server"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbdiscovery"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbinventory"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbled"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbpower"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbredfish"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
	"google.golang.org/grpc"
)

// export BMC_APP_ADDR=bmc-app
// export BMC_APP_PORT=32415
type Client struct {
	conn *grpc.ClientConn

	c pbbmc.BmcServiceClient
}

// New create new instance of connection
// do not forget to close it!
func New() (c *Client, err error) {
	conn, err := server.Client()
	if err != nil {
		return
	}
	c = &Client{
		conn: conn,
		c:    pbbmc.NewBmcServiceClient(conn),
	}
	return
}

func (c *Client) Health(ctx context.Context) (err error) {
	_, err = c.c.Health(ctx, &pbbmc.Empty{})
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
	request := pbled.Request{
		Id:       id,
		Username: username,
		Password: password,
		State:    string(state),
	}
	return c.c.SwitchLed(ctx, &request)
}

func (c *Client) SwitchPower(ctx context.Context, id, username, password string, resetType redfish.ResetType) (device *pbredfish.Device, err error) {
	request := pbpower.Request{
		Id:       id,
		Username: username,
		Password: password,
		Type:     string(resetType),
	}
	return c.c.SwitchPower(ctx, &request)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
