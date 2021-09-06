package quotas

import (
	"github.com/G-Core/gcorelabscloud-go/client/quotas/v2/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/quota/v2/quotas"
	"github.com/urfave/cli/v2"
)

var Commands = cli.Command{
	Name:  "quotas",
	Usage: "GCloud quotas API",
	Subcommands: []*cli.Command{
		{
			Name:  "list",
			Usage: "List quotas",
			Subcommands: []*cli.Command{
				&quotasListCombinedSubCommands,
				&quotasListGlobalSubCommands,
				&quotasListRegionalSubCommands,
			},
		},
	},
}

var quotasListCombinedSubCommands = cli.Command{
	Name:     "combined",
	Usage:    "Get combined client quotas, regional and global",
	Category: "list",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "client-id",
			Aliases: []string{"c"},
			Usage:   "Id of the client",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewQuotaClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := quotas.ListCombinedOpts{}
		clientID := c.Int("client-id")
		if clientID != 0 {
			opts.ClientID = clientID
		}
		result, err := quotas.ListCombined(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var quotasListGlobalSubCommands = cli.Command{
	Name:     "global",
	Usage:    "Get global quota",
	Category: "list",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "client-id",
			Aliases:  []string{"c"},
			Usage:    "Id of the client",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewQuotaClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		clientID := c.Int("client-id")
		result, err := quotas.ListGlobal(client, clientID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var quotasListRegionalSubCommands = cli.Command{
	Name:     "regional",
	Usage:    "Get regional quota",
	Category: "list",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "client-id",
			Aliases:  []string{"c"},
			Usage:    "Id of the client",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "region-id",
			Aliases:  []string{"r"},
			Usage:    "Id of the region",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewQuotaClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		clientID := c.Int("client-id")
		regionID := c.Int("region-id")
		result, err := quotas.ListRegional(client, clientID, regionID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}
