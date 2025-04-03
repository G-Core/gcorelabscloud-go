package regionsaccess

import (
	"github.com/G-Core/gcorelabscloud-go/client/regions/v1/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/regionaccess/v1/regionsaccess"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/urfave/cli/v2"
)

var (
	resellerIDText = "reseller_id is mandatory argument"
)

var regionAccessListCommand = cli.Command{
	Name:     "list",
	Usage:    "List regions",
	Category: "regionaccess",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "reseller-id",
			Aliases: []string{"r"},
			Usage:   "Reseller ID",
		},
		&cli.IntFlag{
			Name:    "client-id",
			Aliases: []string{"c"},
			Usage:   "Client ID. It is assumed to be null when unspecified.",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewRegionClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := regionsaccess.ListOpts{
			ResellerID: c.Int("reseller-id"),
			ClientID:   c.Int("client-id"),
		}
		results, err := regionsaccess.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var regionAccessDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete reseller's region access records",
	ArgsUsage: "<reseller_id>",
	Category:  "regionaccess",
	Action: func(c *cli.Context) error {
		resellerID, err := flags.GetFirstIntArg(c, resellerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewRegionClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		if err = regionsaccess.Delete(client, resellerID).ExtractErr(); err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var regionAccessCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create or update the limit of regions for clients of the reseller",
	Category: "regionaccess",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "reseller-id",
			Aliases: []string{"r"},
			Usage:   "\n\nReseller ID. It's null when a client doesn't have a reseller.",
		},
		&cli.IntFlag{
			Name:    "client-id",
			Aliases: []string{"c"},
			Usage:   "Client ID. If you want to set limits to all reseller clients, skip this field. The client_id has priority over the reseller_id. If both client_id and reseller_id are specified, reseller_id must be the correct reseller_id of the client.",
		},
		&cli.BoolFlag{
			Name:  "all-edge-region",
			Usage: "If true, allow access to all edge regions, regardless of the content of region_ids array.",
		},
		&cli.IntSliceFlag{
			Name:  "regions-id",
			Usage: "List of available region ids.",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewRegionClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		clientID := c.Int("client-id")
		resellerID := c.Int("reseller-id")

		opts := regionsaccess.CreateOpts{
			AccessAllEdgeRegions: c.Bool("all-edge-region"),
			RegionIDs:            c.IntSlice("regions-id"),
			ClientID:             &clientID,
			ResellerID:           &resellerID,
		}

		result, err := regionsaccess.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var Commands = cli.Command{
	Name:  "regionaccess",
	Usage: "GCloud regions access API",
	Subcommands: []*cli.Command{
		&regionAccessListCommand,
		&regionAccessDeleteCommand,
		&regionAccessCreateCommand,
	},
}
