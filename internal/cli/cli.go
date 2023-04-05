// Copyright 2023 NJWS Inc.

package cli

import (
	"github.com/foliagecp/inventory-bmc-app/internal/health"
	"github.com/foliagecp/inventory-bmc-app/internal/server"
	"github.com/urfave/cli/v2"
)

var (
	CLI = cli.NewApp()

	version = "v0.1.0"
)

func init() {
	CLI.Usage = "Redfish app tool"
	CLI.Version = version

	CLI.Commands = []*cli.Command{
		&cli.Command{
			Name:        "health",
			Description: "App health check",
			Action: func(ctx *cli.Context) (err error) {
				return health.Run(ctx.Context)
			},
		},
		&cli.Command{
			Name:        "run",
			Description: "Run app",
			Action: func(ctx *cli.Context) (err error) {
				return server.Run(ctx.Context)
			},
		},
	}
}
