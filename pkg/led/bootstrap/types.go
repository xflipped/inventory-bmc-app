// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"

	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
)

var (
	registerTypes = []*pbcmdb.RegisterTypeMessage{}
)

func createTypes(ctx context.Context) (err error) {
	// new types are not required
	return
}
