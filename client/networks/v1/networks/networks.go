package networks

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/networks/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	cmeta "github.com/G-Core/gcorelabscloud-go/client/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/availablenetworks"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"
)

var networkIDText = "network_id is mandatory argument"

var networkListCommand = cli.Command{
	Name:     "list",
	Usage:    "List networks",
	Category: "network",
	Action: func(c *cli.Context) error {
		client, err := client.NewNetworkClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		pages, err := networks.List(client, nil).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := networks.ExtractNetworks(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var availableNetworkListCommand = cli.Command{
	Name:     "list-available",
	Usage:    "List available networks",
	Category: "network",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "network-id",
			Aliases:  []string{"i"},
			Usage:    "show subnets of the specific network",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "network-type",
			Aliases:  []string{"t"},
			Usage:    "filter network by network type (vlan or vxlan)",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewAvailableNetworkClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := availablenetworks.ListOpts{
			NetworkID:   c.String("network-id"),
			NetworkType: c.String("network-type"),
		}

		result, err := availablenetworks.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var networkGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get network information",
	ArgsUsage: "<network_id>",
	Category:  "network",
	Action: func(c *cli.Context) error {
		networkID, err := flags.GetFirstStringArg(c, networkIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewNetworkClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		network, err := networks.Get(client, networkID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if network == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(network, c.String("format"))
		return nil
	},
}

var networkDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete network by ID",
	ArgsUsage: "<network_id>",
	Category:  "network",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		networkID, err := flags.GetFirstStringArg(c, networkIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewNetworkClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := networks.Delete(client, networkID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := networks.Get(client, networkID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete network with ID: %s", networkID)
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

var networkUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update network",
	ArgsUsage: "<network_id>",
	Category:  "network",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Network name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		networkID, err := flags.GetFirstStringArg(c, networkIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewNetworkClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := networks.UpdateOpts{
			Name: c.String("name"),
		}

		network, err := networks.Update(client, networkID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if network == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(network, c.String("format"))
		return nil

	},
}

var networkCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create network",
	Category: "network",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Network name",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "create-router",
			Usage:    "Create network router",
			Required: false,
		},
		&cli.StringFlag{
			Name:        "type",
			Usage:       "Network type `vlan` or `vxlan`. Default to `vxlan`.",
			Required:    false,
			DefaultText: "vxlan",
			Value:       "vxlan",
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := client.NewNetworkClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := networks.CreateOpts{
			Name:         c.String("name"),
			CreateRouter: c.Bool("create-router"),
			Type:         c.String("type"),
		}
		results, err := networks.Create(client, opts).Extract()
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
			networkID, err := networks.ExtractNetworkIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve network ID from task info: %w", err)
			}
			network, err := networks.Get(client, networkID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get network with ID: %s. Error: %w", networkID, err)
			}
			utils.ShowResults(network, c.String("format"))
			return nil, nil
		})
	},
}

var networkInstancePortCommand = cli.Command{
	Name:      "instance_port",
	Usage:     "List of instance ports by ID",
	ArgsUsage: "<network_id>",
	Category:  "network",
	Action: func(c *cli.Context) error {
		networkID, err := flags.GetFirstStringArg(c, networkIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "instance_port")
			return err
		}
		client, err := client.NewNetworkClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := networks.ListAllInstancePort(client, networkID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var Commands = cli.Command{
	Name:  "network",
	Usage: "GCloud networks API",
	Subcommands: []*cli.Command{
		&networkListCommand,
		&availableNetworkListCommand,
		&networkGetCommand,
		&networkDeleteCommand,
		&networkCreateCommand,
		&networkUpdateCommand,
		&extensionCommands,
		&networkInstancePortCommand,
		{
			Name:  "metadata",
			Usage: "Network metadata",
			Subcommands: []*cli.Command{
				cmeta.NewMetadataListCommand(
					client.NewNetworkClientV1,
					"Get networks metadata",
					"<network_id>",
					"network_id is mandatory argument",
				),
				cmeta.NewMetadataGetCommand(
					client.NewNetworkClientV1,
					"Show network metadata by key",
					"<network_id>",
					"network_id is mandatory argument",
				),
				cmeta.NewMetadataDeleteCommand(
					client.NewNetworkClientV1,
					"Delete network metadata by key",
					"<network_id>",
					"network_id is mandatory argument",
				),
				cmeta.NewMetadataCreateCommand(
					client.NewNetworkClientV1,
					"Create network metadata. It would update existing keys",
					"<network_id>",
					"network_id is mandatory argument",
				),
				cmeta.NewMetadataUpdateCommand(
					client.NewNetworkClientV1,
					"Update network metadata. It overriding existing records",
					"<network_id>",
					"network_id is mandatory argument",
				),
				cmeta.NewMetadataReplaceCommand(
					client.NewNetworkClientV1,
					"Replace network metadata. It replace existing records",
					"<network_id>",
					"network_id is mandatory argument",
				),
			},
		},
	},
}
