package servers

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/client"
	tasksclient "github.com/G-Core/gcorelabscloud-go/client/tasks/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/servers"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/urfave/cli/v2"
)

func listServersAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "list")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	servers, err := servers.ListAll(gpuClient, clusterID, servers.ListOpts{})
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(servers, c.String("format"))
	return nil
}

func listBaremetalServersAction(c *cli.Context) error {
	return listServersAction(c, client.NewGPUBaremetalClientV3)
}

func listVirtualServersAction(c *cli.Context) error {
	return listServersAction(c, client.NewGPUVirtualClientV3)
}

func deleteServerAction(c *cli.Context, gpuType string, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "list")
		return cli.Exit("cluster ID is required", 1)
	}

	serverID := c.Args().Get(1)
	if serverID == "" {
		_ = cli.ShowCommandHelp(c, "list")
		return cli.Exit("server ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	opts := servers.DeleteServerOpts{
		AllFloatingIPs:      c.Bool("delete-all-floating-ips"),
		AllReservedFixedIPs: c.Bool("delete-all-reserved-fixed-ips"),
	}
	// this flag is only applicable for virtual clusters
	if gpuType == "virtual" {
		opts.AllVolumes = c.Bool("delete-all-volumes")
	}
	results, err := servers.Delete(gpuClient, clusterID, serverID, opts).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	tc, err := tasksclient.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, tc, results, c.Bool("d"), func(task tasks.TaskID) (interface{}, error) {
		servers, err := servers.ListAll(gpuClient, clusterID, servers.ListOpts{})
		if err != nil {
			return nil, err
		}
		return servers, nil
	})
}

func deleteBaremetalServerAction(c *cli.Context) error {
	return deleteServerAction(c, "baremetal", client.NewGPUBaremetalClientV3)
}

func deleteVirtualServerAction(c *cli.Context) error {
	return deleteServerAction(c, "virtual", client.NewGPUVirtualClientV3)
}

// BaremetalCommands returns commands for baremetal GPU servers
func BaremetalCommands() *cli.Command {
	return &cli.Command{
		Name:        "servers",
		Usage:       "Manage baremetal GPU cluster servers",
		Description: "Commands for managing servers in baremetal GPU clusters",
		Subcommands: []*cli.Command{
			{
				Name:        "list",
				Usage:       "List servers in a baremetal GPU cluster",
				Description: "List all servers in a specific baremetal GPU cluster",
				Category:    "servers",
				ArgsUsage:   "<cluster_id>",
				Action:      listBaremetalServersAction,
			},
			{
				Name:        "delete",
				Usage:       "Delete server from a baremetal GPU cluster",
				Description: "Delete a specific server from a baremetal GPU cluster",
				Category:    "servers",
				ArgsUsage:   "<cluster_id> <server_id>",
				Flags: append([]cli.Flag{
					&cli.BoolFlag{
						Name:     "delete-all-floating-ips",
						Usage:    "delete all server floating ips",
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "delete-all-reserved-fixed-ips",
						Usage:    "delete all server reserved fixed ips",
						Required: false,
					},
				}, flags.WaitCommandFlags...),
				Action: deleteBaremetalServerAction,
			},
		},
	}
}

// VirtualCommands returns commands for virtual GPU servers
func VirtualCommands() *cli.Command {
	return &cli.Command{
		Name:        "servers",
		Usage:       "Manage virtual GPU cluster servers",
		Description: "Commands for managing servers in virtual GPU clusters",
		Subcommands: []*cli.Command{
			{
				Name:        "list",
				Usage:       "List servers in a virtual GPU cluster",
				Description: "List all servers in a specific virtual GPU cluster",
				Category:    "servers",
				ArgsUsage:   "<cluster_id>",
				Action:      listVirtualServersAction,
			},
			{
				Name:        "delete",
				Usage:       "Delete server from a virtual GPU cluster",
				Description: "Delete a specific server from a virtual GPU cluster",
				Category:    "servers",
				ArgsUsage:   "<cluster_id> <server_id>",
				Flags: append([]cli.Flag{
					&cli.BoolFlag{
						Name:     "delete-all-floating-ips",
						Usage:    "delete all server floating ips",
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "delete-all-reserved-fixed-ips",
						Usage:    "delete all server reserved fixed ips",
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "delete-all-volumes",
						Usage:    "delete all server volumes",
						Required: false,
					},
				}, flags.WaitCommandFlags...),
				Action: deleteVirtualServerAction,
			},
		},
	}
}
