// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"

	"git.fg-tech.ru/listware/go-core/pkg/module"
)

func (a *Agent) workerFunction(ctx module.Context) (err error) {
	fmt.Println("worker", string(ctx.CmdbContext()))
	return
}
