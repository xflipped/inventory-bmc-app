// Copyright 2023 NJWS Inc.

package agent

import (
	"context"

	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/reset/agent/types"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

type Agent struct {
	ctx    context.Context
	cancel context.CancelFunc

	executor executor.Executor

	m module.Module
}

// Run agent
func Run(ctx context.Context) (err error) {
	a := &Agent{}
	a.ctx, a.cancel = context.WithCancel(ctx)

	if a.executor, err = executor.New(); err != nil {
		return
	}

	return a.run()
}

func (a *Agent) run() (err error) {
	defer a.executor.Close()

	log.Infof("run reset agent")

	a.osSignalCtrl()

	a.m = module.New(types.Namespace, module.WithPort(31004))

	log.Infof("use port (%d)", a.m.Port())

	if err = a.m.Bind(types.ResetFunctionType, a.resetFunction); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(a.ctx)

	go func() {
		defer cancel()

		if err = a.m.RegisterAndListen(ctx); err != nil {
			log.Error(err)
		}

	}()

	if err = a.entrypoint(); err != nil {
		return
	}

	<-ctx.Done()

	return
}
