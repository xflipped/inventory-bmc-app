// Copyright 2023 NJWS Inc.

package admin

import (
	"git.fg-tech.ru/listware/go-core/pkg/executor"
	discovery "github.com/foliagecp/inventory-bmc-app/pkg/discovery/cli"
	inventory "github.com/foliagecp/inventory-bmc-app/pkg/inventory/cli"
	led "github.com/foliagecp/inventory-bmc-app/pkg/led/cli"
	reset "github.com/foliagecp/inventory-bmc-app/pkg/reset/cli"
	"github.com/foliagecp/inventory-bmc-app/pkg/upgrade"
	"github.com/manifoldco/promptui"
	"github.com/stmcginnis/gofish/common"
	"github.com/stmcginnis/gofish/redfish"
	"github.com/urfave/cli/v2"
)

var (
	CLI = cli.NewApp()

	version = "v0.1.0"
)

func init() {
	CLI.Usage = "Admin redfish tool"
	CLI.Version = version

	CLI.Commands = []*cli.Command{
		&cli.Command{
			Name:        "discovery",
			Description: "Discovery bmc (create or update bmc)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "addr",
					Aliases:  []string{"a"},
					Usage:    "BMC addr",
					Required: true,
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				executor, err := executor.New()
				if err != nil {
					return
				}
				defer executor.Close()

				addr := ctx.String("addr")
				return discovery.Discovery(ctx.Context, executor, addr)
			},
		},

		&cli.Command{
			Name:        "inventory",
			Description: "Inventory bmc (update bios/bmc credentials and execute inventory)",
			Flags: []cli.Flag{
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
			},
			Action: func(ctx *cli.Context) (err error) {
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
				return inventory.Inventory(ctx.Context, query, login, password)
			},
		},

		&cli.Command{
			Name:        "fwupgrade",
			Description: "Update bios/bmc firmware",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "query",
					Aliases:  []string{"q"},
					Required: true,
					Usage:    "qdsl query to redfish device",
				},

				&cli.StringFlag{
					Name:  "type",
					Usage: "type",
					Value: "BMC",
				},

				&cli.StringFlag{
					Name:  "target",
					Usage: "target",
					Value: "/redfish/v1/UpdateService/FirmwareInventory/BMC",
				},

				&cli.StringFlag{
					Name:     "file",
					Usage:    "file",
					Required: true,
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				query := ctx.String("query")
				file := ctx.String("file")
				ftype := ctx.String("type")
				target := ctx.String("target")
				return upgrade.Upgrade(ctx.Context, query, file, ftype, target)
			},
		},

		&cli.Command{
			Name:        "led",
			Description: "Update chassis indicator LED",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "query",
					Aliases:  []string{"q"},
					Required: true,
					Usage:    "qdsl query to redfish device",
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				query := ctx.String("query")
				prompt := promptui.Select{
					Label: "Select IndicatorLED mode",
					Items: []common.IndicatorLED{common.BlinkingIndicatorLED, common.OffIndicatorLED, common.LitIndicatorLED},
				}
				_, mode, err := prompt.Run()
				if err != nil {
					return
				}

				return led.Led(ctx.Context, query, common.IndicatorLED(mode))
			},
		},

		&cli.Command{
			Name:        "reset",
			Description: "Reset system",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "query",
					Aliases:  []string{"q"},
					Required: true,
					Usage:    "qdsl query to redfish device",
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				query := ctx.String("query")
				prompt := promptui.Select{
					Label: "Select reset type",
					Items: []redfish.ResetType{redfish.OnResetType, redfish.GracefulShutdownResetType, redfish.ForceOffResetType, redfish.ForceRestartResetType},
				}
				_, resetType, err := prompt.Run()
				if err != nil {
					return
				}

				return reset.Reset(ctx.Context, query, redfish.ResetType(resetType))
			},
		},
	}
}
