// Copyright 2023 NJWS Inc.

package utils

import (
	"sync"
)

type Parallel struct {
	wg   sync.WaitGroup
	merr *Error
}

func NewParallel() *Parallel {
	return &Parallel{merr: NewMultiError()}
}

func (p *Parallel) Exec(f func() error) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if err := f(); err != nil {
			p.merr.Append(err)
			return
		}
	}()
}

func (p *Parallel) Wait() error {
	p.wg.Wait()
	return p.merr.ErrOrNil()
}
