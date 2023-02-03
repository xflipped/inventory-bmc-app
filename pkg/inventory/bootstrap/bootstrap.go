// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
)

func register(ctx context.Context, exec executor.Executor) (err error) {
	// create types
	if err = createTypes(ctx); err != nil {
		return
	}

	// create objects
	if err = createObjects(ctx); err != nil {
		return
	}

	// example of trigger between types, not used in inventory
	// create links
	// if err = createLinks(ctx); err != nil {
	// 	return
	// }

	message, err := system.Register(types.App, registerTypes, registerObjects, nil)
	if err != nil {
		return
	}

	return exec.ExecSync(ctx, message)
}

func Run() (err error) {
	ctx := context.Background()

	exec, err := executor.New()
	if err != nil {
		return
	}
	defer exec.Close()

	if err = register(ctx, exec); err != nil {
		return
	}

	return
}
