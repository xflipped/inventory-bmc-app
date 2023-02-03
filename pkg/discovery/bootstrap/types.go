// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex/types"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
)

var (
	registerTypes = []*pbcmdb.RegisterTypeMessage{}
)

func createType(ctx context.Context, pt *types.Type) (err error) {
	query := fmt.Sprintf("%s.types.root", pt.Schema.Title)
	elements, err := qdsl.Qdsl(ctx, query)
	if err != nil {
		return
	}

	// TODO already exists
	if len(elements) > 0 {
		return
	}

	message, err := system.RegisterType(pt, true)
	if err != nil {
		return
	}

	registerTypes = append(registerTypes, message)
	return
}

func createRedfishDeviceContainerType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishDeviceContainer{})
	return createType(ctx, pt)
}

func createRedfishDeviceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&device.RedfishDevice{})
	return createType(ctx, pt)
}

func createTypes(ctx context.Context) (err error) {
	if err = createRedfishDeviceContainerType(ctx); err != nil {
		return
	}

	if err = createRedfishDeviceType(ctx); err != nil {
		return
	}
	return
}
