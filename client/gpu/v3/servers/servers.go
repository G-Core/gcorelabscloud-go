package servers

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/servers"
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

	serversList := servers.List(gpuClient, clusterID)
	if serversList.Err != nil {
		return cli.Exit(serversList.Err, 1)
	}

	results, err := servers.ExtractServers(serversList)
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(results, c.String("format"))
	return nil
}

func listVirtualServersAction(c *cli.Context) error {
	return listServersAction(c, client.NewGPUVirtualClientV3)
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
		},
	}
}
