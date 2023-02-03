// Copyright 2023 NJWS Inc.

package admin

import (
	"github.com/foliagecp/inventory-bmc-app/pkg/admin/agent"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

var (
	CLI = cli.NewApp()

	version = "v0.1.0"
)

func init() {
	CLI.Usage = "Admin redfish tool"
	CLI.Version = version
	CLI.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "query",
			Aliases:  []string{"q"},
			Required: true,
			Usage:    "qdsl query to redfish device",
		},

		&cli.StringFlag{
			Name:     "login",
			Aliases:  []string{"l"},
			Required: true,
			Usage:    "Login to set",
			Value:    "root",
		},
	}
	CLI.Action = func(ctx *cli.Context) (err error) {
		query := ctx.String("query")
		login := ctx.String("login")

		prompt := promptui.Prompt{
			Label: "Enter password",
			Mask:  '*',
		}

		password, err := prompt.Run()
		if err != nil {
			return
		}

		return agent.ChangeCredentials(ctx.Context, query, login, password)
	}
}
