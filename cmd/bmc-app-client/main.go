// Copyright 2023 NJWS Inc.

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/stmcginnis/gofish/redfish"

	"github.com/foliagecp/inventory-bmc-app/pkg/bmc"
	"github.com/stmcginnis/gofish/common"
)

func main() {
	ctx := context.Background()

	client, err := bmc.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	fmt.Println("discovery")

	ctx2, _ := context.WithTimeout(ctx, time.Second*5)

	device, err := client.Discovery(ctx2, "https://192.168.77.102/")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("discovery", device)

	fmt.Println("list")

	devices, err := client.ListDevices(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, device := range devices.GetItems() {
		fmt.Println("list", device)
	}

	fmt.Println("inventory")

	device, err = client.Inventory(ctx, device.GetId(), "admin", "P@ssw0rd")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("inventory", device)

	fmt.Println("led")

	device, err = client.SwitchLed(ctx, device.GetId(), "admin", "P@ssw0rd", common.BlinkingIndicatorLED)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("led", device)

	fmt.Println("power")

	device, err = client.SwitchPower(ctx, device.GetId(), "admin", "P@ssw0rd", redfish.OnResetType)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("power", device)
}
