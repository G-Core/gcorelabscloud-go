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

	clusterDetails := clusters.Get(gpuClient, clusterID)
	if clusterDetails.Err != nil {
		return cli.Exit(clusterDetails.Err, 1)
	}

	utils.ShowResults(clusterDetails.Body, c.String("format"))
	return nil
}

func showVirtualClusterAction(c *cli.Context) error {
	return showClusterAction(c, client.NewGPUVirtualClientV3)
}

func showBaremetalClusterAction(c *cli.Context) error {
	return showClusterAction(c, client.NewGPUBaremetalClientV3)
}

// listClustersAction handles the common logic for listing both virtual and baremetal clusters
func listClustersAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}
	opts := &clusters.ListOpts{}
	pages, err := clusters.List(gpuClient, opts).AllPages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	clusterList, err := clusters.ExtractClusters(pages)
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(clusterList, c.String("format"))
	return nil
}

func listVirtualClustersAction(c *cli.Context) error {
	return listClustersAction(c, client.NewGPUVirtualClientV3)
}

func listBaremetalClustersAction(c *cli.Context) error {
	return listClustersAction(c, client.NewGPUBaremetalClientV3)
}

// BaremetalCommands returns commands for managing baremetal GPU clusters
func BaremetalCommands() *cli.Command {
	return &cli.Command{
		Name:        "clusters",
		Usage:       "Manage baremetal GPU clusters",
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
			{
				Name:        "list",
				Usage:       "List baremetal GPU clusters",
				Description: "List all baremetal GPU clusters",
				Category:    "clusters",
				ArgsUsage:   " ",
				Action:      listBaremetalClustersAction,
			},
		},
	}
}

// VirtualCommands returns commands for managing virtual GPU clusters
func VirtualCommands() *cli.Command {
	return &cli.Command{
		Name:        "clusters",
		Usage:       "Manage virtual GPU clusters",
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
			{
				Name:        "list",
				Usage:       "List virtual GPU clusters",
				Description: "List all virtual GPU clusters",
				Category:    "clusters",
				ArgsUsage:   " ",
				Action:      listVirtualClustersAction,
			},
		},
	}
}
