// Copyright 2023 NJWS Inc.

package bmc

import (
	"context"

	"github.com/foliagecp/inventory-bmc-app/sdk/pbbmc"
)

// device - remove if re-discovery and uuid updated
func (b *BmcApp) Health(ctx context.Context, empty *pbbmc.Empty) (*pbbmc.Empty, error) {
	return empty, nil
}
