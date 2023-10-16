package namespaces

import (
	"fmt"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/faas/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"
)

const (
	namespaceNameText = "namespace_name is mandatory argument"
)

var Commands = cli.Command{
	Name:  "namespaces",
	Usage: "GCloud FaaS namespaces API",
	Subcommands: []*cli.Command{
		&namespaceListCommand,
		&namespaceShowCommand,
		&namespaceCreateCommand,
		&namespaceUpdateCommand,
		&namespaceDeleteCommand,
	},
}

var namespaceListCommand = cli.Command{
	Name:     "list",
	Usage:    "Display list of namespaces",
	Category: "namespaces",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "search",
			Aliases:  []string{"s"},
			Usage:    "show namespaces whose names contain provided value",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "limit",
			Aliases:  []string{"l"},
			Usage:    "limit the number of returned namespaces. Limited by max limit value of 1000",
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
			Usage:    "order namespaces by transmitted fields and directions (name.asc).",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		return listNamespaces(c)
	},
}

func listNamespaces(c *cli.Context) error {
	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	opts := faas.ListOpts{
		Limit:   c.Int("limit"),
		Offset:  c.Int("offset"),
		Search:  c.String("search"),
		OrderBy: c.String("order"),
	}

	results, err := faas.ListNamespaceALL(cl, opts)
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(results, c.String("format"))

	return nil
}

var namespaceShowCommand = cli.Command{
	Name:      "show",
	Usage:     "Show namespace",
	Category:  "namespaces",
	ArgsUsage: "<namespace>",
	Action: func(c *cli.Context) error {
		return showNamespace(c)
	},
}

func showNamespace(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, namespaceNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "show")
		return err
	}

	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	namespace, err := faas.GetNamespace(cl, name).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(namespace, c.String("format"))

	return nil
}

var namespaceCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "create namespace.",
	Category: "namespaces",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "function name",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "envs",
			Aliases:  []string{"e"},
			Usage:    "environment variables. 'env_name'='value'",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "function description",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		return createNamespace(c)
	},
}

func extractEnvs(list []string) map[string]string {
	envs := make(map[string]string)
	for _, item := range list {
		ss := strings.Split(item, "=")
		envs[ss[0]] = ss[1]
	}

	return envs
}

func createNamespace(c *cli.Context) error {
	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	opts := faas.CreateNamespaceOpts{
		Name:        c.String("name"),
		Description: c.String("description"),
		Envs:        extractEnvs(c.StringSlice("envs")),
	}

	results, err := faas.CreateNamespace(cl, opts).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, cl, results, true, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(cl, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}

		ns, err := faas.GetNamespace(cl, c.String("name")).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get namespace with Name: %s. Error: %w", c.String("name"), err)
		}

		return ns, nil
	})
}

var namespaceUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "update namespace.",
	Category:  "namespaces",
	ArgsUsage: "<namespace>",
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:     "envs",
			Aliases:  []string{"e"},
			Usage:    "environment variables. 'env_name'='value'",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "namespace description",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		return updateNamespace(c)
	},
}

func updateNamespace(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, namespaceNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "update")
		return err
	}

	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	var opts faas.UpdateNamespaceOpts
	if c.IsSet("description") {
		opts.Description = c.String("description")
	}

	if c.IsSet("envs") {
		opts.Envs = extractEnvs(c.StringSlice("envs"))
	}

	results, err := faas.UpdateNamespace(cl, name, opts).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, cl, results, true, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(cl, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}

		ns, err := faas.GetNamespace(cl, name).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get namespace with Name: %s. Error: %w", name, err)
		}

		return ns, nil
	})
}

var namespaceDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete namespace",
	Category:  "namespaces",
	ArgsUsage: "<namespace>",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		return deleteNamespace(c)
	},
}

func deleteNamespace(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, namespaceNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "show")
		return err
	}

	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	results, err := faas.DeleteNamespace(cl, name).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, cl, results, false, func(task tasks.TaskID) (interface{}, error) {
		_, err := faas.GetNamespace(cl, name).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot delete namespace with Name: %s", name)
		}

		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
}
