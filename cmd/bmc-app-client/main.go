// Copyright 2023 NJWS Inc.

package main

import (
	"context"
	"fmt"

	"github.com/foliagecp/inventory-bmc-app/internal/server"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbdiscovery"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbinventory"
)

func main() {

	ctx := context.Background()

	conn, err := server.Client()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := pbbmc.NewBmcServiceClient(conn)

	discoveryRequest := &pbdiscovery.Request{
		Url: "https://ip/",
	}

	fmt.Println("discovery")

	device, err := client.Discovery(ctx, discoveryRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("discovery ", device.GetUrl(), " id ", device.GetId())

	fmt.Println("list")

	devices, err := client.ListDevices(ctx, &pbbmc.Empty{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, device := range devices.GetItems() {
		fmt.Println(device.GetUrl())
	}

	fmt.Println("inventory")

	inventoryRequest := &pbinventory.Request{
		Id:       device.GetId(),
		Username: "admin",
		Password: "admin",
	}

	response, err := client.Inventory(ctx, inventoryRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(response.GetId())
}
