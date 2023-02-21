// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"
	"time"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types"
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

	message, err := system.Register(types.App, registerTypes, registerObjects, nil)
	if err != nil {
		return
	}

	return exec.ExecSync(ctx, message)
}

func Run(ctx context.Context) (err error) {
	exec, err := executor.New(executor.WithTimeout(time.Minute))
	if err != nil {
		return
	}
	defer exec.Close()
	return register(ctx, exec)
}
