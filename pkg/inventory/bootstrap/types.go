// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/vertex/types"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
)

var (
	registerTypes = []*pbcmdb.RegisterTypeMessage{}
)

// FIXME will be custom structs?
type RedfishService struct {
	*gofish.Service
}

type RedfishSystem struct {
	*redfish.ComputerSystem
}

type RedfishBios struct {
	*redfish.Bios
}

type RedfishLed struct {
	Led common.IndicatorLED `json:"led"`
}

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

func createRedfishServiceType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishService{})
	return createType(ctx, pt)
}

func createRedfishSystemType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishSystem{})
	return createType(ctx, pt)
}

func createRedfishBiosType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishBios{})
	return createType(ctx, pt)
}

func createRedfishLedType(ctx context.Context) (err error) {
	pt := types.ReflectType(&RedfishLed{})
	return createType(ctx, pt)
}

func createTypes(ctx context.Context) (err error) {
	if err = createRedfishServiceType(ctx); err != nil {
		return
	}
	if err = createRedfishSystemType(ctx); err != nil {
		return
	}
	if err = createRedfishBiosType(ctx); err != nil {
		return
	}
	if err = createRedfishLedType(ctx); err != nil {
		return
	}
	return
}
