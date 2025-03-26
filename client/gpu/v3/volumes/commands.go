package volumes

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/volumes"
)

func listVirtualVolumesAction(c *cli.Context) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "list")
		return cli.Exit("cluster ID is required", 1)
	}

	cl, err := client.NewGPUVirtualClientV3(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	// Get project ID from CLI context or service client
	projectID := c.Int("project")
	if projectID == 0 {
		projectID = cl.ProjectID
		if projectID == 0 {
			return cli.Exit(fmt.Errorf("project ID must be provided with --project flag or GCLOUD_PROJECT environment variable"), 1)
		}
	}

	// Get region ID from CLI context or service client
	regionID := c.Int("region")
	if regionID == 0 {
		regionID = cl.RegionID
		if regionID == 0 {
			return cli.Exit(fmt.Errorf("region ID must be provided with --region flag or GCLOUD_REGION environment variable"), 1)
		}
	}

	// Set project and region in the client
	cl.ProjectID = projectID
	cl.RegionID = regionID

	result := volumes.List(cl, clusterID)
	pages, err := result.AllPages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	volumeList, err := volumes.ExtractVolumes(pages)
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(volumeList, c.String("format"))
	return nil
}

// VirtualCommands returns commands for managing virtual GPU volumes
func VirtualCommands() *cli.Command {
	return &cli.Command{
		Name:        "volumes",
		Usage:       "Manage virtual GPU cluster volumes",
		Description: "Commands for managing virtual GPU cluster volumes",
		Subcommands: []*cli.Command{
			{
				Name:     "list",
				Usage:    "List virtual GPU cluster volumes",
				Category: "volumes",
				Action:   listVirtualVolumesAction,
			},
		},
	}
}
