// Copyright 2023 NJWS Inc.

package health

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/pkg/bmc"
)

func Run(ctx context.Context) (err error) {
	client, err := bmc.New()
	if err != nil {
		return
	}
	defer client.Close()
	return client.Health(ctx)
}
