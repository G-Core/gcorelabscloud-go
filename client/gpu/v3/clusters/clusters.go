package clusters

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/clusters"
	"github.com/urfave/cli/v2"
)

func showClusterAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "show")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	imageDetails := clusters.Get(gpuClient, clusterID)
	if imageDetails.Err != nil {
		return cli.Exit(imageDetails.Err, 1)
	}

	utils.ShowResults(imageDetails.Body, c.String("format"))
	return nil
}

func showVirtualClusterAction(c *cli.Context) error {
	return showClusterAction(c, client.NewGPUVirtualClientV3)
}

func showBaremetalClusterAction(c *cli.Context) error {
	return showClusterAction(c, client.NewGPUBaremetalClientV3)
}

// BaremetalCommands returns commands for managing baremetal GPU clusters
func BaremetalCommands() *cli.Command {
	return &cli.Command{
		Name:        "clusters",
		Usage:       "Manage baremetal GPU images",
		Description: "Commands for managing baremetal GPU clusters",
		Subcommands: []*cli.Command{
			{
				Name:        "show",
				Usage:       "Show baremetal GPU cluster details",
				Description: "Show details of a specific baremetal GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      showBaremetalClusterAction,
			},
		},
	}
}

// VirtualCommands returns commands for managing virtual GPU clusters
func VirtualCommands() *cli.Command {
	return &cli.Command{
		Name:        "clusters",
		Usage:       "Manage virtual GPU images",
		Description: "Commands for managing virtual GPU clusters",
		Subcommands: []*cli.Command{
			{
				Name:        "show",
				Usage:       "Show virtual GPU cluster details",
				Description: "Show details of a specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      showVirtualClusterAction,
			},
		},
	}
}
