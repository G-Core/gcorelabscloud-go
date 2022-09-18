package subnets

import (
	"fmt"
	cmeta "github.com/G-Core/gcorelabscloud-go/client/utils/metadata"
	"net"

	"github.com/G-Core/gcorelabscloud-go/client/subnets/v1/client"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"
)

var subnetIDText = "subnet_id is mandatory argument"

func getDNSNameservers(c *cli.Context) ([]net.IP, error) {
	dns := c.StringSlice("dns-nameserver")
	var result []net.IP
	for _, server := range dns {
		ip := net.ParseIP(server)
		if ip == nil {
			return result, fmt.Errorf("cannot parse dns nameserver ip: %s", server)
		}
		result = append(result, ip)
	}
	return result, nil
}

func GetHostRoutes(c *cli.Context) ([]subnets.HostRoute, error) {
	destinations := c.StringSlice("route-destination")
	hops := c.StringSlice("route-nexthop")
	if len(destinations) > 0 && len(destinations) != len(hops) {
		return nil, fmt.Errorf("should be equal number of route-destination and route-nexthop arguments")
	}
	var result []subnets.HostRoute
	for idx, desc := range destinations {
		dst, err := gcorecloud.ParseCIDRString(desc)
		if err != nil {
			return result, fmt.Errorf("cannot parse route destination: %s: %w", desc, err)
		}
		hop := net.ParseIP(hops[idx])
		if hop == nil {
			return result, fmt.Errorf("cannot parse route nexthop: %s", hops[idx])
		}
		result = append(result, subnets.HostRoute{
			Destination: *dst,
			NextHop:     hop,
		})
	}
	return result, nil
}

var subnetListCommand = cli.Command{
	Name:     "list",
	Usage:    "List subnets",
	Category: "subnet",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "network_id",
			Aliases:  []string{"n"},
			Usage:    "Network ID",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewSubnetClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := subnets.ListOpts{
			NetworkID: c.String("network_id"),
		}

		pages, err := subnets.List(client, opts).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := subnets.ExtractSubnets(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var subnetGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get subnet information",
	ArgsUsage: "<subnet_id>",
	Category:  "subnet",
	Action: func(c *cli.Context) error {
		subnetID, err := flags.GetFirstStringArg(c, subnetIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewSubnetClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := subnets.Get(client, subnetID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if task == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(task, c.String("format"))
		return nil
	},
}

var subnetDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete subnet by ID",
	ArgsUsage: "<subnet_id>",
	Category:  "subnet",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		subnetID, err := flags.GetFirstStringArg(c, subnetIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewSubnetClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := subnets.Delete(client, subnetID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := subnets.Get(client, subnetID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete subnet with ID: %s", subnetID)
			}
			switch err.(type) {
			case gcorecloud.ErrDefault404:
				return nil, nil
			default:
				return nil, err
			}
		})

	},
}

var subnetUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update subnet",
	ArgsUsage: "<subnet_id>",
	Category:  "subnet",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Subnet name",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "route-destination",
			Aliases:  []string{"rd"},
			Usage:    "Subnet route destination",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "route-nexthop",
			Aliases:  []string{"rh"},
			Usage:    "Subnet route next hop",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		subnetID, err := flags.GetFirstStringArg(c, subnetIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewSubnetClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		dns, err := getDNSNameservers(c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		hostRoutes, err := GetHostRoutes(c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		opts := subnets.UpdateOpts{
			Name:           c.String("name"),
			DNSNameservers: dns,
			HostRoutes:     hostRoutes,
		}

		subnet, err := subnets.Update(client, subnetID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if subnet == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(subnet, c.String("format"))
		return nil

	},
}

var subnetCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create subnet",
	Category: "subnet",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Subnet name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "cidr",
			Aliases:  []string{"c"},
			Usage:    "Subnet CIDR",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "network-id",
			Aliases:  []string{"i"},
			Usage:    "Subnet network ID",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "enable-dhcp",
			Usage:    "Enable DHCP",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "connect-to-router",
			Usage:    "Connect subnet to router",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "dns-nameserver",
			Aliases:  []string{"dns"},
			Usage:    "Subnet dns nameserver",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "route-destination",
			Aliases:  []string{"rd"},
			Usage:    "Subnet route destination",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "route-nexthop",
			Aliases:  []string{"rh"},
			Usage:    "Subnet route next hop",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "gateway-ip",
			Aliases:  []string{"gw"},
			Usage:    "Gateway ip",
			Required: false,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := client.NewSubnetClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		cidr, err := gcorecloud.ParseCIDRString(c.String("cidr"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		dns, err := getDNSNameservers(c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		hostRoutes, err := GetHostRoutes(c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		opts := subnets.CreateOpts{
			Name:                   c.String("name"),
			EnableDHCP:             c.Bool("enable-dhcp"),
			CIDR:                   *cidr,
			NetworkID:              c.String("network-id"),
			ConnectToNetworkRouter: c.Bool("connect-to-router"),
			DNSNameservers:         dns,
			HostRoutes:             hostRoutes,
		}

		rawGateway := c.String("gateway-ip")
		if rawGateway != "" {
			gatewayIP := net.ParseIP(rawGateway)
			opts.GatewayIP = &gatewayIP
		}

		results, err := subnets.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			subnetID, err := subnets.ExtractSubnetIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve subnet ID from task info: %w", err)
			}
			subnet, err := subnets.Get(client, subnetID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get subnet with ID: %s. Error: %w", subnetID, err)
			}
			utils.ShowResults(subnet, c.String("format"))
			return nil, nil
		})
	},
}

var Commands = cli.Command{
	Name:  "subnet",
	Usage: "GCloud subnets API",
	Subcommands: []*cli.Command{
		&subnetListCommand,
		&subnetGetCommand,
		&subnetDeleteCommand,
		&subnetCreateCommand,
		&subnetUpdateCommand,
		{
			Name:  "metadata",
			Usage: "Network metadata",
			Subcommands: []*cli.Command{
				cmeta.NewMetadataListCommand(
					client.NewSubnetClientV1,
					"Get subnet metadata",
					"<subnet_id>",
					"subnet_id is mandatory argument",
				),
				cmeta.NewMetadataGetCommand(
					client.NewSubnetClientV1,
					"Show subnet metadata by key",
					"<subnet_id>",
					"subnet_id is mandatory argument",
				),
				cmeta.NewMetadataDeleteCommand(
					client.NewSubnetClientV1,
					"Delete subnet metadata by key",
					"<subnet_id>",
					"subnet_id is mandatory argument",
				),
				cmeta.NewMetadataCreateCommand(
					client.NewSubnetClientV1,
					"Create subnet metadata. It would update existing keys",
					"<subnet_id>",
					"subnet_id is mandatory argument",
				),
				cmeta.NewMetadataUpdateCommand(
					client.NewSubnetClientV1,
					"Update subnet metadata. It overriding existing records",
					"<subnet_id>",
					"subnet_id is mandatory argument",
				),
				cmeta.NewMetadataReplaceCommand(
					client.NewSubnetClientV1,
					"Replace subnet metadata. It replace existing records",
					"<subnet_id>",
					"subnet_id is mandatory argument",
				),
			},
		},
	},
}
