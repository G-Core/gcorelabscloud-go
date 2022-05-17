package apitokens

import (
	"fmt"
	"strings"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/apitokens/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/apitoken/v1/apitokens"
	"github.com/G-Core/gcorelabscloud-go/gcore/apitoken/v1/types"
	"github.com/urfave/cli/v2"
)

var apiTokenIDText = "apitoken_id is mandatory argument"

var Commands = cli.Command{
	Name:  "apitokens",
	Usage: "GCloud api token API. Could be used only with platform client type",
	Subcommands: []*cli.Command{
		&apiTokenListCommand,
		&apiTokenGetCommand,
		&apiTokenDeleteCommand,
		&apiTokenCreateCommand,
	},
}

var apiTokenListCommand = cli.Command{
	Name:  "list",
	Usage: "List api tokens",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "client-id",
			Aliases:  []string{"c"},
			Usage:    "Client id",
			Required: true,
		},
		&cli.IntFlag{
			Name:  "role-id",
			Usage: fmt.Sprintf("Available value: %s", strings.Join(types.RoleIDType(0).StringList(), ", ")),
		},
		&cli.IntFlag{
			Name:  "issued-by",
			Usage: "User's ID. Use to get API tokens issued by a particular user.",
		},
		&cli.IntFlag{
			Name:  "not-issued-by",
			Usage: "User's ID. Use to get API tokens issued by anyone except a particular user.",
		},
		// todo: figure out
		// &cli.BoolFlag{
		// 	Name: "deleted",
		// 	Usage: "If set, in the response wil be included only deleted tokens",
		// },
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewAPITokenClient(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := apitokens.ListOpts{}
		if c.Int("role-id") != 0 {
			opts.RoleID = types.RoleIDType(c.Int("role-id"))
		}

		if c.Int("issued-by") != 0 {
			opts.IssuedBy = c.Int("issued-by")
		}

		if c.Int("not-issued-by") != 0 {
			opts.NotIssuedBy = c.Int("not-issued-by")
		}

		results, err := apitokens.List(client, c.Int("client-id"), opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var apiTokenGetCommand = cli.Command{
	Name:  "show",
	Usage: "Show api token",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "client-id",
			Aliases:  []string{"c"},
			Usage:    "Client id",
			Required: true,
		},
	},
	ArgsUsage: "<apitoken_id>",
	Action: func(c *cli.Context) error {
		apiTokenID, err := flags.GetFirstIntArg(c, apiTokenIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}

		client, err := client.NewAPITokenClient(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := apitokens.Get(client, c.Int("client-id"), apiTokenID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var apiTokenDeleteCommand = cli.Command{
	Name:  "delete",
	Usage: "Delete api token",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "client-id",
			Aliases:  []string{"c"},
			Usage:    "Client id",
			Required: true,
		},
	},
	ArgsUsage: "<apitoken_id>",
	Action: func(c *cli.Context) error {
		apiTokenID, err := flags.GetFirstIntArg(c, apiTokenIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}

		client, err := client.NewAPITokenClient(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		err = apitokens.Delete(client, c.Int("client-id"), apiTokenID).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var apiTokenCreateCommand = cli.Command{
	Name:  "create",
	Usage: "Create api token",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "client-id",
			Aliases:  []string{"c"},
			Usage:    "Client id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Required: true,
		},
		&cli.StringFlag{
			Name: "description",
		},
		&cli.IntFlag{
			Name:     "role-id",
			Usage:    fmt.Sprintf("Available value: %s", strings.Join(types.RoleIDType(0).StringList(), ", ")),
			Required: true,
		},
		&cli.StringFlag{
			Name:     "role-name",
			Usage:    fmt.Sprintf("Available value: %s", strings.Join(types.RoleNameType("").StringList(), ", ")),
			Required: true,
		},
		&cli.StringFlag{
			Name:  "expiration-time",
			Usage: "Date when the API token becomes expired (ISO 8086/RFC 3339 format), UTC. If null, then the API token will never expire.",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewAPITokenClient(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := apitokens.CreateOpts{
			Name:        c.String("name"),
			Description: c.String("description"),
			ClientUser: apitokens.CreateClientUser{
				Role: apitokens.ClientRole{
					ID:   types.RoleIDType(c.Int("role-id")),
					Name: types.RoleNameType(c.String("role-name")),
				},
			},
		}

		expDateRaw := c.String("expiration-date")
		if expDateRaw != "" {
			expDate, err := time.Parse(gcorecloud.RFC3339ZZ, expDateRaw)
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
			opts.ExpDate = &gcorecloud.JSONRFC3339Z{Time: expDate}
		}

		result, err := apitokens.Create(client, c.Int("client-id"), opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}
