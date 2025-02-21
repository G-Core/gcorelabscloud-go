package flavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
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
		Name:     "show-disabled",
		Aliases:  []string{"sd"},
		Usage:    "Show disabled flavors (by default disabled flavors are not shown)",
		Required: false,
	},
}

// listFlavorsAction handles the common logic for listing both virtual and baremetal flavors
func listFlavorsAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	cl, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	includePrices := c.Bool("include-prices")
	disabled := c.Bool("show-disabled")
	opts := flavors.ListOpts{
		IncludePrices: &includePrices,
		Disabled:      &disabled,
	}

	results, err := flavors.List(cl, opts).AllPages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	flavorList, err := flavors.ExtractFlavors(results)
	if err != nil {
		return cli.Exit(err, 1)
	}
	utils.ShowResults(flavorList, c.String("format"))
	return nil
}

func listBaremetalFlavorsAction(c *cli.Context) error {
	return listFlavorsAction(c, client.NewGPUBaremetalClientV3)
}

func listVirtualFlavorsAction(c *cli.Context) error {
	return listFlavorsAction(c, client.NewGPUVirtualClientV3)
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
