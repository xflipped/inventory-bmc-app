// Copyright 2023 NJWS Inc.

package utils

import (
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// Error - is a multierror wrapper
type Error struct {
	mu  sync.Mutex
	err *multierror.Error
}

// Append errors into the list
func (e *Error) Append(errs ...error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.err = multierror.Append(e.err, errs...)
}

// ErrOrNil - return error or nil if error list is empty
func (e *Error) ErrOrNil() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.err.ErrorOrNil()
}

// NewMultiError \\
func NewMultiError(errs ...error) (merr *Error) {
	merr = &Error{
		mu: sync.Mutex{},
		err: &multierror.Error{
			Errors: errs,
			ErrorFormat: func(errs []error) string {
				s := []string{}
				for _, err := range errs {
					s = append(s, err.Error())
				}
				return strings.Join(s, "\n")
			},
		},
	}
	return
}
