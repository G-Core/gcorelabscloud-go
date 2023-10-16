package keys

import (
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/faas/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
)

const keyNameText = "key_name is mandatory argument"

var Commands = cli.Command{
	Name:  "keys",
	Usage: "GCloud FaaS keys API",
	Subcommands: []*cli.Command{
		&keyListCommand,
		&keyShowCommand,
		&keyCreateCommand,
		&keyUpdateCommand,
		&keyDeleteCommand,
	},
}

var keyListCommand = cli.Command{
	Name:     "list",
	Usage:    "List API keys.",
	Category: "api keys",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "search",
			Aliases:  []string{"s"},
			Usage:    "show keys whose names contain provided value",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "limit",
			Aliases:  []string{"l"},
			Usage:    "limit the number of returned keys. Limited by max limit value of 1000",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "offset",
			Aliases:  []string{"o"},
			Usage:    "offset value is used to exclude the first set of records from the result",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "order",
			Usage:    "order keys by transmitted fields and directions (name.asc).",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		return listKeys(c)
	},
}

func listKeys(c *cli.Context) error {
	cl, err := client.NewFaaSKeysClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
	}

	opts := faas.ListOpts{
		Limit:   c.Int("limit"),
		Offset:  c.Int("offset"),
		Search:  c.String("search"),
		OrderBy: c.String("order"),
	}

	results, err := faas.ListKeysAll(cl, opts)
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(results, c.String("format"))

	return nil
}

var keyShowCommand = cli.Command{
	Name:      "show",
	Usage:     "Show API keys.",
	Category:  "api keys",
	ArgsUsage: "<key_name>",
	Action: func(c *cli.Context) error {
		return showKey(c)
	},
}

func showKey(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, keyNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "show")
		return err
	}

	cl, err := client.NewFaaSKeysClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	key, err := faas.GetKey(cl, name).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(key, c.String("format"))
	return nil
}

var keyCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create API key.",
	Category: "api keys",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "key name",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "functions",
			Aliases:  []string{"fs"},
			Usage:    "<function namespace/name>",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "key description",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "expire",
			Usage:    "when key will expire. Format 2023-07-31T00:00:00Z",
			Required: false,
		},
		&cli.PathFlag{
			Name:     "file",
			Usage:    "where to put API key secret",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		return createKey(c)
	},
}

func createKey(c *cli.Context) error {
	cl, err := client.NewFaaSKeysClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	opts := faas.CreateKeyOpts{
		Name:        c.String("name"),
		Description: c.String("description"),
	}

	if c.IsSet("expire") {
		expireRaw := c.String("expire")
		t, err := time.Parse(gcorecloud.RFC3339ZZ, expireRaw)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.Exit("Invalid format for functions", 1)
		}
		expire := gcorecloud.JSONRFC3339ZZ{Time: t}
		opts.Expire = &expire
	}

	if c.IsSet("functions") {
		var functions []faas.KeysFunction
		items := c.StringSlice("functions")
		for _, item := range items {
			ss := strings.Split(item, "/")
			if len(ss) != 2 {
				_ = cli.ShowCommandHelp(c, "create")
				return cli.Exit("Invalid format for functions", 1)
			}
			functions = append(functions, faas.KeysFunction{
				Name:      ss[1],
				Namespace: ss[0],
			})
		}
		opts.Functions = functions
	}

	key, err := faas.CreateKey(cl, opts)
	if err != nil {
		return cli.Exit(err, 1)
	}

	if c.IsSet("file") {
		secret := key.Secret
		if err := utils.WriteToFile(c.Path("file"), []byte(secret)); err != nil {
			return cli.Exit(err, 1)
		}
	}

	utils.ShowResults(key, c.String("format"))

	return nil
}

var keyUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update API key.",
	ArgsUsage: "<key_name>",
	Category:  "api keys",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "functions",
			Aliases:  []string{"fs"},
			Usage:    "<function namespace/name>",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "key description",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "expire",
			Usage:    "when key will expire. Format 2023-07-31T00:00:00Z",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		return updateKey(c)
	},
}

func updateKey(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, keyNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "update")
		return err
	}

	cl, err := client.NewFaaSKeysClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	opts := faas.UpdateKeyOpts{}

	if c.IsSet("description") {
		opts.Description = c.String("description")
	}

	if c.IsSet("expire") {
		expireRaw := c.String("expire")
		t, err := time.Parse(gcorecloud.RFC3339ZZ, expireRaw)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.Exit("Invalid format for functions", 1)
		}
		expire := gcorecloud.JSONRFC3339ZZ{Time: t}
		opts.Expire = &expire
	}

	if c.IsSet("functions") {
		var functions []faas.KeysFunction
		items := c.StringSlice("functions")
		for _, item := range items {
			ss := strings.Split(item, "/")
			if len(ss) != 2 {
				_ = cli.ShowCommandHelp(c, "update")
				return cli.Exit("Invalid format for functions", 1)
			}
			functions = append(functions, faas.KeysFunction{
				Name:      ss[1],
				Namespace: ss[0],
			})
		}

		opts.Functions = functions
	}

	key, err := faas.UpdateKey(cl, name, opts)
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(key, c.String("format"))

	return nil
}

var keyDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete API key.",
	Category:  "api keys",
	ArgsUsage: "<key_name>",
	Action: func(c *cli.Context) error {
		return deleteKey(c)
	},
}

func deleteKey(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, keyNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "update")
		return err
	}

	cl, err := client.NewFaaSKeysClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	err = faas.DeleteKey(cl, name)
	switch err.(type) {
	case gcorecloud.ErrDefault404:
		return nil
	default:
		return cli.Exit(err, 1)
	}
}
