// Copyright 2023 NJWS Inc.

package admin

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"git.fg-tech.ru/listware/go-core/pkg/executor"
	discovery "github.com/foliagecp/inventory-bmc-app/pkg/discovery/cli"
	inventoryAgent "github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent"
	inventory "github.com/foliagecp/inventory-bmc-app/pkg/inventory/cli"
	ledAgent "github.com/foliagecp/inventory-bmc-app/pkg/led/agent"
	led "github.com/foliagecp/inventory-bmc-app/pkg/led/cli"
	resetAgent "github.com/foliagecp/inventory-bmc-app/pkg/reset/agent"
	reset "github.com/foliagecp/inventory-bmc-app/pkg/reset/cli"
	subscribeAgent "github.com/foliagecp/inventory-bmc-app/pkg/subscribe/agent"
	subscribe "github.com/foliagecp/inventory-bmc-app/pkg/subscribe/cli"
	"github.com/foliagecp/inventory-bmc-app/pkg/upgrade"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
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

	executor, err := executor.New(executor.WithTimeout(time.Second * 60))
	if err != nil {
		return
	}
	defer executor.Close()

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
					Name:     "endpoint",
					Aliases:  []string{"e"},
					Required: true,
					Usage:    "BMC URL",
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
				endpoint := ctx.String("endpoint")
				login := ctx.String("login")

				prompt := promptui.Prompt{
					Label: "Enter password",
					Mask:  '*',
				}

				password, err := prompt.Run()
				if err != nil {
					return
				}

				inventoryPayload := inventoryAgent.InventoryPayload{
					ConnectionParameters: utils.ConnectionParameters{
						Endpoint: endpoint,
						Login:    login,
						Password: password,
					},
				}

				return inventory.Inventory(ctx.Context, executor, query, inventoryPayload)
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
				fileType := ctx.String("type")
				target := ctx.String("target")
				return upgrade.Upgrade(ctx.Context, query, file, fileType, target)
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
					Usage:    "qdsl query to chassis",
				},
				&cli.StringFlag{
					Name:     "endpoint",
					Aliases:  []string{"e"},
					Required: true,
					Usage:    "BMC URL",
				},
				&cli.StringFlag{
					Name:     "login",
					Aliases:  []string{"l"},
					Required: true,
					Usage:    "BMC login",
					Value:    "root",
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				query := ctx.String("query")
				endpoint := ctx.String("endpoint")
				login := ctx.String("login")

				selectPrompt := promptui.Select{
					Label: "Select IndicatorLED mode",
					Items: []common.IndicatorLED{common.BlinkingIndicatorLED, common.OffIndicatorLED, common.LitIndicatorLED},
				}
				_, mode, err := selectPrompt.Run()
				if err != nil {
					return
				}

				prompt := promptui.Prompt{
					Label: "Enter password",
					Mask:  '*',
				}
				password, err := prompt.Run()
				if err != nil {
					return
				}

				ledPayload := ledAgent.LedPayload{
					Led: common.IndicatorLED(mode),
					ConnectionParameters: utils.ConnectionParameters{
						Endpoint: endpoint,
						Login:    login,
						Password: password,
					},
				}
				return led.Led(ctx.Context, executor, query, ledPayload)
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
					Usage:    "qdsl query to computer system",
				},
				&cli.StringFlag{
					Name:     "endpoint",
					Aliases:  []string{"e"},
					Required: true,
					Usage:    "BMC URL",
				},
				&cli.StringFlag{
					Name:     "login",
					Aliases:  []string{"l"},
					Required: true,
					Usage:    "BMC login",
					Value:    "root",
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				query := ctx.String("query")
				endpoint := ctx.String("endpoint")
				login := ctx.String("login")

				selectPrompt := promptui.Select{
					Label: "Select reset type",
					Items: []redfish.ResetType{redfish.OnResetType, redfish.GracefulShutdownResetType, redfish.ForceOffResetType, redfish.ForceRestartResetType, redfish.PowerCycleResetType},
				}
				_, resetType, err := selectPrompt.Run()
				if err != nil {
					return
				}

				prompt := promptui.Prompt{
					Label: "Enter password",
					Mask:  '*',
				}
				password, err := prompt.Run()
				if err != nil {
					return
				}

				resetPayload := resetAgent.ResetPayload{
					ResetType: redfish.ResetType(resetType),
					ConnectionParameters: utils.ConnectionParameters{
						Endpoint: endpoint,
						Login:    login,
						Password: password,
					},
				}

				return reset.Reset(ctx.Context, executor, query, resetPayload)
			},
		},

		&cli.Command{
			Name:        "subscribe",
			Description: "Subscribe to BMC events",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "query",
					Aliases:  []string{"q"},
					Required: true,
					Usage:    "qdsl query to event service",
				},
				&cli.StringSliceFlag{
					Name:     "registries",
					Aliases:  []string{"rp"},
					Required: false,
					Usage:    "comma-separated list of registry prefixes to subscribe to (e.g. Base,Security)",
				},
				&cli.StringSliceFlag{
					Name:     "resources",
					Aliases:  []string{"rt"},
					Required: false,
					Usage:    "comma-separated list of resource types to subscribe to (e.g. TelemetryService,EventService)",
				},
				&cli.StringFlag{
					Name:     "endpoint",
					Aliases:  []string{"e"},
					Required: true,
					Usage:    "BMC URL",
				},
				&cli.StringFlag{
					Name:     "login",
					Aliases:  []string{"l"},
					Required: true,
					Usage:    "BMC login",
					Value:    "root",
				},
			},
			Action: func(ctx *cli.Context) (err error) {
				query := ctx.String("query")
				endpoint := ctx.String("endpoint")
				login := ctx.String("login")
				registryPrefixes := ctx.StringSlice("registries")
				resourceTypes := ctx.StringSlice("resources")

				// Destination URL
				validate := func(input string) (err error) {
					u, err := url.ParseRequestURI(input)
					if err != nil {
						return err
					}

					if !strings.HasPrefix(u.Scheme, "http") {
						return fmt.Errorf("destination should start with http")
					}

					return
				}
				subscribePrompt := promptui.Prompt{
					Label:    "Subscription URL",
					Validate: validate,
				}
				url, err := subscribePrompt.Run()
				if err != nil {
					return
				}

				// Password
				prompt := promptui.Prompt{
					Label: "Enter password",
					Mask:  '*',
				}
				password, err := prompt.Run()
				if err != nil {
					return
				}

				subscribePayload := subscribeAgent.SubscribePayload{
					DestinationUrl:   url,
					RegistryPrefixes: registryPrefixes,
					ResourceTypes:    resourceTypes,
					ConnectionParameters: utils.ConnectionParameters{
						Endpoint: endpoint,
						Login:    login,
						Password: password,
					},
				}
				return subscribe.Subscribe(ctx.Context, executor, query, subscribePayload)
			},
		},
	}
}
