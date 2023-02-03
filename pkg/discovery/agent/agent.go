// Copyright 2023 NJWS Inc.

package agent

import (
	"context"
	"strings"

	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

type Agent struct {
	ctx    context.Context
	cancel context.CancelFunc

	executor executor.Executor
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

	return a.entrypoint()
}
