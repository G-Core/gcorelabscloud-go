package routers

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/routers/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/subnets/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

var routerIDText = "router_id is mandatory argument"

var RouterCommands = cli.Command{
	Name:  "router",
	Usage: "GCloud router API",
	Subcommands: []*cli.Command{
		&routerListSubCommand,
		&routerGetSubCommand,
		&routerUpdateSubCommand,
		&routerDeleteSubCommand,
		&routerCreateSubCommand,
		&routerAttachSubCommand,
		&routerDetachSubCommand,
	},
}

var routerListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "List routers",
	Category: "router",
	Action: func(c *cli.Context) error {
		client, err := client.NewRouterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := routers.ListOpts{}
		results, err := routers.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if results != nil {
			utils.ShowResults(results, c.String("format"))
		}
		return nil
	},
}

var routerGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show router",
	ArgsUsage: "<router_id>",
	Category:  "router",
	Action: func(c *cli.Context) error {
		routerID, err := flags.GetFirstStringArg(c, routerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewRouterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := routers.Get(client, routerID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var routerDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete router",
	ArgsUsage: "<router_id>",
	Category:  "router",
	Action: func(c *cli.Context) error {
		routerID, err := flags.GetFirstStringArg(c, routerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}

		client, err := client.NewRouterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := routers.Delete(client, routerID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var routerCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create router",
	Category: "router",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Router name",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:    "route-destination",
			Aliases: []string{"rd"},
			Usage:   "CIDR of destination IPv4 subnet.",
		},
		&cli.StringSliceFlag{
			Name:    "route-nexthop",
			Aliases: []string{"rh"},
			Usage:   "IPv4 address to forward traffic to if it's destination IP matches 'destination' CIDR.",
		},
		&cli.BoolFlag{
			Name:  "enable-snat",
			Usage: "Is SNAT enabled. Defaults to true.",
		},
		&cli.StringFlag{
			Name:  "network-id",
			Usage: "ID of the external network. Required if gateway-type set to 'manual'.",
		},
		&cli.StringFlag{
			Name:  "gateway-type",
			Usage: "Router gateway type. Available value is 'manual', 'default'.",
		},
		&cli.StringSliceFlag{
			Name:  "subnet-id",
			Usage: "ID of the subnet to attach to.",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := client.NewRouterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		routes, err := subnets.GetHostRoutes(c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		enableSNAT := c.Bool("enable-snat")
		opts := routers.CreateOpts{
			Name: c.String("name"),
			ExternalGatewayInfo: routers.GatewayInfo{
				Type:       types.GatewayType(c.String("gateway-type")),
				EnableSNat: &enableSNAT,
				NetworkID:  c.String("network-id"),
			},
		}

		if len(routes) > 0 {
			opts.Routes = routes
		}

		subnetIDs := c.StringSlice("subnet-id")
		if len(subnetIDs) > 0 {
			ifs := make([]routers.Interface, 0, len(subnetIDs))
			for i, subnetID := range subnetIDs {
				ifs[i] = routers.Interface{
					Type:     types.SubnetInterfaceType,
					SubnetID: subnetID,
				}
			}
			opts.Interfaces = ifs
		}

		results, err := routers.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			routerID, err := routers.ExtractRouterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve router ID from task info: %w", err)
			}
			router, err := routers.Get(client, routerID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get router with ID: %s. Error: %w", routerID, err)
			}
			utils.ShowResults(router, c.String("format"))
			return nil, nil
		})
	},
}

var routerUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Update router",
	Category:  "router",
	ArgsUsage: "<router_id>",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Router name",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:    "route-destination",
			Aliases: []string{"rd"},
			Usage:   "CIDR of destination IPv4 subnet.",
		},
		&cli.StringSliceFlag{
			Name:    "route-nexthop",
			Aliases: []string{"rh"},
			Usage:   "IPv4 address to forward traffic to if it's destination IP matches 'destination' CIDR.",
		},
		&cli.BoolFlag{
			Name:  "enable-snat",
			Usage: "Is SNAT enabled. Defaults to true.",
		},
		&cli.StringFlag{
			Name:     "network-id",
			Usage:    "ID of the external network.",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		routerID, err := flags.GetFirstStringArg(c, routerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}

		client, err := client.NewRouterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		routes, err := subnets.GetHostRoutes(c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		enableSNAT := c.Bool("enable-snat")
		opts := routers.UpdateOpts{
			Name: c.String("name"),
			ExternalGatewayInfo: routers.GatewayInfo{
				Type:       types.ManualGateway,
				EnableSNat: &enableSNAT,
				NetworkID:  c.String("network-id"),
			},
		}

		if len(routes) > 0 {
			opts.Routes = routes
		}

		result, err := routers.Update(client, routerID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var routerAttachSubCommand = cli.Command{
	Name:      "attach",
	Usage:     "Attach router",
	Category:  "router",
	ArgsUsage: "<router_id>",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "subnet-id",
			Usage:    "ID of the subnet to attach to.",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		routerID, err := flags.GetFirstStringArg(c, routerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "attach")
			return err
		}

		client, err := client.NewRouterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := routers.Attach(client, routerID, c.String("subnet-id")).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var routerDetachSubCommand = cli.Command{
	Name:      "detach",
	Usage:     "Detach router",
	Category:  "router",
	ArgsUsage: "<router_id>",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "subnet-id",
			Usage:    "ID of the subnet to detach from.",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		routerID, err := flags.GetFirstStringArg(c, routerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "detach")
			return err
		}

		client, err := client.NewRouterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := routers.Detach(client, routerID, c.String("subnet-id")).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		utils.ShowResults(result, c.String("format"))
		return nil
	},
}
