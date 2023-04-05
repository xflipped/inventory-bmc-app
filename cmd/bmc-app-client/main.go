// Copyright 2023 NJWS Inc.

package main

import (
	"context"
	"fmt"

	"github.com/foliagecp/inventory-bmc-app/internal/server"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbdiscovery"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbinventory"
	"github.com/foliagecp/inventory-bmc-app/sdk/pbled"
	"github.com/stmcginnis/gofish/common"
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
		Url: "https://192.168.77.102/",
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
	fmt.Println("listed", len(devices.GetItems()))

	for _, device := range devices.GetItems() {
		fmt.Println(device)
	}

	fmt.Println("inventory")

	inventoryRequest := &pbinventory.Request{
		Id:       device.GetId(),
		Username: "admin",
		Password: "P@ssw0rd",
	}

	_ = inventoryRequest
	// response, err := client.Inventory(ctx, inventoryRequest)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(response.GetId())

	fmt.Println("led")

	ledRequest := pbled.Request{
		Id:       device.GetId(),
		Username: "admin",
		Password: "P@ssw0rd",
		State:    string(common.BlinkingIndicatorLED),
	}
	device, err = client.SwitchLed(ctx, &ledRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(device)
}
