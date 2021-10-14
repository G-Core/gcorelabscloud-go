package flavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flavors/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/flavor/v1/flavors"

	"github.com/urfave/cli/v2"
)

var flavorListCommand = cli.Command{
	Name:     "list",
	Usage:    "List flavors",
	Category: "flavor",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "include_prices",
			Aliases: []string{"p"},
			Usage:   "Include prices",
		},
		&cli.BoolFlag{
			Name:    "baremetal",
			Aliases: []string{"bm"},
			Usage:   "show only baremetal flavors",
		},
	},
	Action: func(c *cli.Context) error {
		var err error
		var cl *gcorecloud.ServiceClient
		cl, err = client.NewFlavorClientV1(c)

		if c.Bool("baremetal") {
			cl, err = client.NewBmFlavorClientV1(c)
		}
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		prices := c.Bool("include_prices")
		opts := flavors.ListOpts{
			IncludePrices: &prices,
		}
		results, err := flavors.ListAll(cl, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var Commands = cli.Command{
	Name:  "flavor",
	Usage: "GCloud flavors API",
	Subcommands: []*cli.Command{
		&flavorListCommand,
	},
}
