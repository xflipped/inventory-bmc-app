// Copyright 2023 NJWS Inc.

package agent

import (
	"time"

	"github.com/foliagecp/inventory-bmc-app/pkg/led/agent/types"
)

func (a *Agent) entrypoint() (err error) {
	// wait router, need to register port
	time.Sleep(time.Millisecond * 50)

	return a.createOrUpdateFunctionLink(types.FunctionContainerPath, types.LedFunctionPath, types.LedFunctionLink)
}
