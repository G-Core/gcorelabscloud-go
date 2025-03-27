package clusters

import (
	"fmt"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/client"
	task_client "github.com/G-Core/gcorelabscloud-go/client/tasks/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
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

func deleteClusterAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "delete")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	results, err := clusters.Delete(gpuClient, clusterID).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := task_client.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, results, false, func(task tasks.TaskID) (interface{}, error) {
		_, err := clusters.Get(gpuClient, clusterID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete GPU cluster with ID: %s. Error: %w", clusterID, err)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
}

func deleteVirtualClusterAction(c *cli.Context) error {
	return deleteClusterAction(c, client.NewGPUVirtualClientV3)
}

func deleteBaremetalClusterAction(c *cli.Context) error {
	return deleteClusterAction(c, client.NewGPUBaremetalClientV3)
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
			{
				Name:        "delete",
				Usage:       "Delete baremetal GPU cluster",
				Description: "Delete a specific baremetal GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      deleteBaremetalClusterAction,
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
			{
				Name:        "delete",
				Usage:       "Delete virtual GPU cluster",
				Description: "Delete a specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      deleteVirtualClusterAction,
			},
		},
	}
}
