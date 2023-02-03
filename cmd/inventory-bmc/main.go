// Copyright 2023 NJWS Inc.

package main

import (
	"fmt"
	"os"

	"github.com/foliagecp/inventory-bmc-app/pkg/inventory"
)

func main() {
	if err := inventory.CLI.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
