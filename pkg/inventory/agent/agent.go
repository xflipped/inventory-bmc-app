// Copyright 2023 NJWS Inc.

package agent

import (
	"context"
	"strings"

	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
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
func Run() (err error) {
	a := &Agent{}
	a.ctx, a.cancel = context.WithCancel(context.Background())

	if a.executor, err = executor.New(); err != nil {
		return
	}

	return a.run()
}

func appendPath(paths ...string) string {
	return strings.Join(paths, ".")
}

func (a *Agent) run() (err error) {
	defer a.executor.Close()

	log.Infof("run system agent")

	a.osSignalCtrl()

	a.m = module.New(types.Namespace, module.WithPort(31001))

	log.Infof("use port (%d)", a.m.Port())

	if err = a.m.Bind(types.FunctionType, a.workerFunction); err != nil {
		return
	}

	// // TODO move to another app
	// if err = a.m.Bind("monit", monit.Monit); err != nil {
	// 	return
	// }

	ctx, cancel := context.WithCancel(a.ctx)

	go func() {
		defer cancel()

		if err = a.m.RegisterAndListen(ctx); err != nil {
			log.Error(err)
			return
		}

	}()

	if err = a.entrypoint(); err != nil {
		return
	}

	<-ctx.Done()

	return
}
