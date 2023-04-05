// Copyright 2023 NJWS Inc.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/foliagecp/inventory-bmc-app/internal/server"
)

func main() {

	if err := server.Run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
