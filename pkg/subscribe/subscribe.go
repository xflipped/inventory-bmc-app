// Copyright 2023 NJWS Inc.

package subscribe

import (
	"github.com/foliagecp/inventory-bmc-app/pkg/subscribe/agent"
	"github.com/foliagecp/inventory-bmc-app/pkg/subscribe/bootstrap"
	"github.com/urfave/cli/v2"
)

var (
	CLI = cli.NewApp()

	version = "v0.1.0"
)

func init() {
	CLI.Usage = "Subscribe to BMC events tool"
	CLI.Version = version
	CLI.Commands = []*cli.Command{
		&cli.Command{
			Name:        "bootstrap",
			Description: "Bootstrap subscribe to BMC events tool",
			Action: func(ctx *cli.Context) (err error) {
				return bootstrap.Run(ctx.Context)
			},
		},
		&cli.Command{
			Name:        "run",
			Description: "Run subscribe to BMC events tool",
			Action: func(ctx *cli.Context) (err error) {
				if err = bootstrap.Run(ctx.Context); err != nil {
					return
				}
				return agent.Run(ctx.Context)
			},
		},
	}
}
