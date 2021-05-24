package reservedfixedips

import (
	"fmt"
	"net"

	"github.com/urfave/cli/v2"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/reservedfixedips/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

var portIDText = "port_id is mandatory argument"

var Commands = cli.Command{
	Name:  "fixed_ip",
	Usage: "GCloud reserved fixed ip API",
	Subcommands: []*cli.Command{
		&reservedFixedIPListSubCommand,
		&reservedFixedIPGetSubCommand,
		&reservedFixedIPDeleteSubCommand,
		&reservedFixedIPCreateSubCommand,
		&reservedFixedIPSwitchVIPSubCommand,
		&reservedFixedIPListInstancePortSubCommand,
		&reservedFixedIPAddPortSubCommand,
		&reservedFixedIPReplacePortSubCommand,
		&reservedFixedIPListAvailablePortSubCommand,
	},
}

var reservedFixedIPListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "List reserved fixed ip",
	Category: "fixed_ip",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "external-only",
			Usage: "Set if the response should only list public IP addresses.",
		},
		&cli.BoolFlag{
			Name:  "internal-only",
			Usage: "Set if the response should only list private IP addresses.",
		},
		&cli.BoolFlag{
			Name:  "available-only",
			Usage: "Set if the response should only list IP addresses that are not attached to any instance.",
		},
		&cli.BoolFlag{
			Name:  "vip-only",
			Usage: "Set if the response should only list VIPs.",
		},
		&cli.StringFlag{
			Name:  "device-id",
			Usage: "Filter IPs by device ID it is attached to",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewReservedFixedIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := reservedfixedips.ListOpts{
			ExternalOnly:  c.Bool("external-only"),
			InternalOnly:  c.Bool("internal-only"),
			AvailableOnly: c.Bool("available-only"),
			VipOnly:       c.Bool("vip-only"),
			DeviceID:      c.String("device-id"),
		}
		results, err := reservedfixedips.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if results != nil {
			utils.ShowResults(results, c.String("format"))
		}
		return nil
	},
}

var reservedFixedIPGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show fixed ip",
	ArgsUsage: "<port_id>",
	Category:  "fixed_ip",
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewReservedFixedIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := reservedfixedips.Get(client, portID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var reservedFixedIPDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete fixed ip",
	ArgsUsage: "<port_id>",
	Category:  "fixed_ip",
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewReservedFixedIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := reservedfixedips.Delete(client, portID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var reservedFixedIPCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create reserved fixed ip",
	Category: "fixed_ip",
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:  "is-vip",
			Usage: "Reserved fixed IP is a VIP.",
		},
		&cli.StringFlag{
			Name:     "type",
			Usage:    "Reserved fixed ip type. Available value is 'external', 'subnet', 'any_subnet', 'ip_address'.",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "network-id",
			Usage: "Reserved fixed IP will be allocated in a subnet of this network. Required if type is 'any_subnet' or 'ip_address'.",
		},
		&cli.StringFlag{
			Name:  "subnet-id",
			Usage: "Reserved fixed IP will be allocated in this subnet. Required if type is 'subnet'.",
		},
		&cli.StringFlag{
			Name:  "ip-address",
			Usage: "Reserved fixed IP will be allocated the given IP address. Required if type is 'ip_address'.",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := client.NewReservedFixedIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := reservedfixedips.CreateOpts{
			Type:      reservedfixedips.ReservedFixedIPType(c.String("type")),
			NetworkID: c.String("network-id"),
			SubnetID:  c.String("subnet-id"),
			IsVip:     c.Bool("is-vip"),
		}

		if opts.Type == reservedfixedips.IPAddress {
			opts.IPAddress = net.ParseIP(c.String("ip-address"))
		}

		results, err := reservedfixedips.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			portID, err := reservedfixedips.ExtractReservedFixedIPIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve reserved fixed ip port ID from task info: %w", err)
			}
			reservedFixedIP, err := reservedfixedips.Get(client, portID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get reserved fixed ip with port ID: %s. Error: %w", portID, err)
			}
			utils.ShowResults(reservedFixedIP, c.String("format"))
			return nil, nil
		})
	},
}

var reservedFixedIPSwitchVIPSubCommand = cli.Command{
	Name:      "switch_vip",
	Usage:     "Switch reserved fixed ip vip status",
	Category:  "fixed_ip",
	ArgsUsage: "<port_id>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "is-vip",
			Usage: "Reserved fixed IP is a VIP. If set status turn into true, if not set into false.",
		},
	},
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}

		client, err := client.NewReservedFixedIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := reservedfixedips.SwitchVIPOpts{
			IsVip: c.Bool("is-vip"),
		}
		results, err := reservedfixedips.SwitchVIP(client, portID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if results != nil {
			utils.ShowResults(results, c.String("format"))
		}
		return nil
	},
}

var reservedFixedIPListInstancePortSubCommand = cli.Command{
	Name:      "list_instance_port",
	Usage:     "List instance ports that share VIP",
	ArgsUsage: "<port_id>",
	Category:  "fixed_ip",
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewReservedFixedIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		pages, err := reservedfixedips.ListConnectedDevice(client, portID).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		result, err := reservedfixedips.ExtractReservedFixedIPs(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var reservedFixedIPListAvailablePortSubCommand = cli.Command{
	Name:      "list_available_port",
	Usage:     "List instance ports that are available for connecting to VIP",
	ArgsUsage: "<port_id>",
	Category:  "fixed_ip",
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewReservedFixedIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		results, err := reservedfixedips.ListAllAvailableDevice(client, portID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if results != nil {
			utils.ShowResults(results, c.String("format"))
		}
		return nil
	},
}

var reservedFixedIPAddPortSubCommand = cli.Command{
	Name:      "add_port",
	Usage:     "Add ports that share VIP",
	Category:  "fixed_ip",
	ArgsUsage: "<port_id>",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "port-id",
			Usage: "Port ID that will share one VIP.",
		},
	},
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}

		client, err := client.NewReservedFixedIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := reservedfixedips.PortsToShareVIPOpts{
			PortIDs: c.StringSlice("port-id"),
		}
		results, err := reservedfixedips.AddPortsToShareVIP(client, portID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if results != nil {
			utils.ShowResults(results, c.String("format"))
		}
		return nil
	},
}

var reservedFixedIPReplacePortSubCommand = cli.Command{
	Name:      "replace_port",
	Usage:     "Replace ports that share VIP",
	Category:  "fixed_ip",
	ArgsUsage: "<port_id>",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "port-id",
			Usage: "Port ID that will share one VIP.",
		},
	},
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}

		client, err := client.NewReservedFixedIPClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := reservedfixedips.PortsToShareVIPOpts{
			PortIDs: c.StringSlice("port-id"),
		}
		results, err := reservedfixedips.ReplacePortsToShareVIP(client, portID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if results != nil {
			utils.ShowResults(results, c.String("format"))
		}
		return nil
	},
}
