package flavors

import (
	"fmt"

	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/flavors"
	"github.com/urfave/cli/v2"
)

var listFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:     "include-prices",
		Aliases:  []string{"p"},
		Usage:    "Include prices in output",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "hide-disabled",
		Aliases:  []string{"hd"},
		Usage:    "Hide disabled flavors (by default shows all flavors)",
		Required: false,
	},
}

// listBaremetalFlavorsAction handles listing baremetal flavors
func listBaremetalFlavorsAction(c *cli.Context) error {
	cl, err := client.NewGPUBaremetalClientV3(c)
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

	includePrices := c.Bool("include-prices")
	hideDisabled := c.Bool("hide-disabled")
	opts := flavors.ListOpts{
		IncludePrices: &includePrices,
		HideDisabled:  &hideDisabled,
	}

	// Set project and region in the client
	cl.ProjectID = projectID
	cl.RegionID = regionID

	results, err := flavors.ListBaremetal(cl, opts).AllPages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	flavorList, err := flavors.ExtractBMFlavors(results)
	if err != nil {
		return cli.Exit(err, 1)
	}
	utils.ShowResults(flavorList, c.String("format"))
	return nil
}

// listVirtualFlavorsAction handles listing virtual flavors
func listVirtualFlavorsAction(c *cli.Context) error {
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

	includePrices := c.Bool("include-prices")
	hideDisabled := c.Bool("hide-disabled")
	opts := flavors.ListOpts{
		IncludePrices: &includePrices,
		HideDisabled:  &hideDisabled,
	}

	// Set project and region in the client
	cl.ProjectID = projectID
	cl.RegionID = regionID

	results, err := flavors.ListVirtual(cl, opts).AllPages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	flavorList, err := flavors.ExtractVMFlavors(results)
	if err != nil {
		return cli.Exit(err, 1)
	}
	utils.ShowResults(flavorList, c.String("format"))
	return nil
}

// BaremetalCommands returns commands for managing baremetal GPU flavors
func BaremetalCommands() *cli.Command {
	return &cli.Command{
		Name:        "flavors",
		Usage:       "Manage baremetal GPU flavors",
		Description: "Commands for managing baremetal GPU flavors",
		Subcommands: []*cli.Command{
			{
				Name:     "list",
				Usage:    "List baremetal GPU flavors",
				Category: "flavors",
				Flags:    listFlags,
				Action:   listBaremetalFlavorsAction,
			},
		},
	}
}

// VirtualCommands returns commands for managing virtual GPU flavors
func VirtualCommands() *cli.Command {
	return &cli.Command{
		Name:        "flavors",
		Usage:       "Manage virtual GPU flavors",
		Description: "Commands for managing virtual GPU flavors",
		Subcommands: []*cli.Command{
			{
				Name:     "list",
				Usage:    "List virtual GPU flavors",
				Category: "flavors",
				Flags:    listFlags,
				Action:   listVirtualFlavorsAction,
			},
		},
	}
}

// Commands returns the list of GPU flavor commands
var Commands = cli.Command{
	Name:        "gpu",
	Usage:       "Manage GPU resources",
	Description: "Parent command for GPU-related operations",
	Category:    "gpu",
	Subcommands: []*cli.Command{
		BaremetalCommands(),
		VirtualCommands(),
	},
}
