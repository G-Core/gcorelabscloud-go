package functions

import (
	"fmt"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/faas/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/faas/v1/keys"
	"github.com/G-Core/gcorelabscloud-go/client/faas/v1/namespaces"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"
)

const functionNameText = "function_name is mandatory argument"

var Commands = cli.Command{
	Name:  "functions",
	Usage: "GCloud FaaS functions API",
	Subcommands: []*cli.Command{
		&namespaces.Commands,
		&keys.Commands,
		&functionListCommand,
		&functionShowCommand,
		&functionDeleteCommand,
		&functionCreateCommand,
		&functionUpdateCommand,
		&functionSaveCommand,
	},
}

var functionListCommand = cli.Command{
	Name:     "list",
	Usage:    "List functions",
	Category: "functions",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "namespace",
			Aliases:  []string{"ns"},
			Usage:    "show functions in namespace",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "search",
			Aliases:  []string{"s"},
			Usage:    "show functions whose names contain provided value",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "limit",
			Aliases:  []string{"l"},
			Usage:    "limit the number of returned functions. Limited by max limit value of 1000",
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
			Usage:    "order functions by transmitted fields and directions (name.asc).",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		return listFunctions(c)
	},
}

func listFunctions(c *cli.Context) error {
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

	results, err := faas.ListFunctionsALL(cl, c.String("namespace"), opts)
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(results, c.String("format"))
	return nil
}

var functionShowCommand = cli.Command{
	Name:      "show",
	Usage:     "Show function",
	Category:  "functions",
	ArgsUsage: "<function_name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "namespace",
			Aliases:  []string{"ns"},
			Usage:    "function namespace",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		return showFunction(c)
	},
}

func showFunction(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, functionNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "show")
		return err
	}

	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	function, err := faas.GetFunction(cl, c.String("namespace"), name).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(function, c.String("format"))
	return nil
}

var functionDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete function.",
	Category:  "functions",
	ArgsUsage: "<function_name>",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "namespace",
			Aliases:  []string{"ns"},
			Usage:    "function namespace",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		return deleteFunction(c)
	},
}

func deleteFunction(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, functionNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "delete")
		return err
	}

	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	results, err := faas.DeleteFunction(cl, c.String("namespace"), name).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, cl, results, false, func(task tasks.TaskID) (interface{}, error) {
		_, err := faas.GetFunction(cl, c.String("namespace"), name).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete function with Name: %s", name)
		}

		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
}

var functionCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "create function.",
	Category: "functions",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "namespace",
			Aliases:  []string{"ns"},
			Usage:    "function namespace",
			Required: true,
		},
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
		&cli.PathFlag{
			Name:     "file",
			Usage:    "file where function declared",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "disabled",
			Usage:    "if true function will not be active from start.",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "keys",
			Usage:    "list of used List of used authentication api keys.",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "enable_api_key",
			Usage:    "if true function will require API keys from 'keys' list for authorization.",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "dependencies",
			Usage:    "function dependencies to install",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "flavor",
			Aliases:  []string{"fl"},
			Usage:    "flavor name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "method",
			Usage:    "main startup method name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "runtime",
			Usage:    "function runtime",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "function description",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout",
			Usage:    "function timeout in seconds",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "scale_min",
			Usage:    "autoscale from",
			Aliases:  []string{"min"},
			Value:    0,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "scale_max",
			Usage:    "autoscale to",
			Aliases:  []string{"max"},
			Value:    1,
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		return createFunction(c)
	},
}

func createFunction(c *cli.Context) error {
	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	code, err := extractCodeFromFile(c.Path("file"))
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	disabled := c.Bool("disabled")
	enableApiKey := c.Bool("enable_api_key")
	opts := faas.CreateFunctionOpts{
		Name:         c.String("name"),
		Description:  c.String("description"),
		Envs:         extractEnvs(c.StringSlice("envs")),
		Runtime:      c.String("runtime"),
		Timeout:      c.Int("timeout"),
		Flavor:       c.String("flavor"),
		CodeText:     code,
		Disabled:     &disabled,
		Keys:         c.StringSlice("keys"),
		Dependencies: c.String("dependencies"),
		EnableApiKey: &enableApiKey,
		MainMethod:   c.String("method"),
	}

	if c.IsSet("scale_min") || c.IsSet("scale_max") {
		min := c.Int("scale_min")
		max := c.Int("scale_max")
		opts.Autoscaling = faas.FunctionAutoscaling{
			MinInstances: &min,
			MaxInstances: &max,
		}
	}

	results, err := faas.CreateFunction(cl, c.String("namespace"), opts).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, cl, results, true, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(cl, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}

		fn, err := faas.GetFunction(cl, c.String("namespace"), c.String("name")).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get function with Name: %s. Error: %w", c.String("name"), err)
		}

		return fn, nil
	})
}

func extractEnvs(list []string) map[string]string {
	envs := make(map[string]string)
	for _, item := range list {
		ss := strings.Split(item, "=")
		envs[ss[0]] = ss[1]
	}

	return envs
}

func extractCodeFromFile(filePath string) (string, error) {
	ok, err := utils.FileExists(filePath)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", fmt.Errorf("file '%s' isn't exist", filePath)
	}

	bytes, err := utils.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

var functionUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "update function.",
	ArgsUsage: "<function_name>",
	Category:  "functions",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "namespace",
			Aliases:  []string{"ns"},
			Usage:    "function namespace",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "envs",
			Aliases:  []string{"e"},
			Usage:    "environment variables. 'env_name'='value'",
			Required: false,
		},
		&cli.PathFlag{
			Name:     "file",
			Usage:    "file where function declared",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "disabled",
			Usage:    "if true function will not be active from start.",
			Required: false,
		},
		&cli.StringFlag{
			Name:  "dependencies",
			Usage: "function dependencies to install",
		},
		&cli.StringSliceFlag{
			Name:     "keys",
			Usage:    "list of used List of used authentication api keys.",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "enable_api_key",
			Usage:    "if true function will require API keys from 'keys' list for authorization.",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "flavor",
			Aliases:  []string{"fl"},
			Usage:    "flavor name",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "method",
			Usage:    "main startup method name",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "timeout",
			Usage:    "function timeout in seconds",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "scale_min",
			Usage:    "autoscale from",
			Aliases:  []string{"min"},
			Required: false,
		},
		&cli.IntFlag{
			Name:     "scale_max",
			Usage:    "autoscale to",
			Aliases:  []string{"max"},
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		return updateFunction(c)
	},
}

func updateFunction(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, functionNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "update")
		return err
	}

	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	opts := faas.UpdateFunctionOpts{}

	if c.IsSet("dependencies") {
		opts.Dependencies = c.String("dependencies")
	}

	if c.IsSet("envs") {
		opts.Envs = extractEnvs(c.StringSlice("envs"))
	}

	if c.IsSet("file") {
		opts.CodeText, err = extractCodeFromFile(c.Path("file"))
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
	}

	if c.IsSet("disabled") {
		disabled := c.Bool("disabled")
		opts.Disabled = &disabled
	}

	if c.IsSet("keys") {
		ss := c.StringSlice("keys")
		opts.Keys = &ss
	}

	if c.IsSet("enable_api_key") {
		apiKey := c.Bool("enable_api_key")
		opts.EnableApiKey = &apiKey
	}

	if c.IsSet("flavor") {
		opts.Flavor = c.String("flavor")
	}

	if c.IsSet("method") {
		opts.MainMethod = c.String("method")
	}

	if c.IsSet("timeout") {
		opts.Timeout = c.Int("timeout")
	}

	opts.Autoscaling = &faas.FunctionAutoscaling{}
	if c.IsSet("scale_min") {
		min := c.Int("scale_min")
		opts.Autoscaling.MinInstances = &min
	}

	if c.IsSet("scale_max") {
		max := c.Int("scale_max")
		opts.Autoscaling.MaxInstances = &max
	}

	results, err := faas.UpdateFunction(cl, c.String("namespace"), name, opts).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, cl, results, true, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(cl, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}

		fn, err := faas.GetFunction(cl, c.String("namespace"), name).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get function with Name: %s. Error: %w", name, err)
		}

		return fn, nil
	})
}

var functionSaveCommand = cli.Command{
	Name:      "save",
	Usage:     "saves function to file.",
	ArgsUsage: "<function_name>",
	Category:  "functions",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "namespace",
			Aliases:  []string{"ns"},
			Usage:    "function namespace",
			Required: true,
		},
		&cli.PathFlag{
			Name:     "file",
			Usage:    "file path",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		return saveFunction(c)
	},
}

func saveFunction(c *cli.Context) error {
	name, err := flags.GetFirstStringArg(c, functionNameText)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "save")
		return err
	}

	cl, err := client.NewFaaSClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	fn, err := faas.GetFunction(cl, c.String("namespace"), name).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	var file string
	if c.IsSet("file") {
		file = c.Path("file")
	} else {
		file = fmt.Sprintf("%s-%s.%s", fn.Name, c.String("namespace"), fn.Runtime[0:2])
	}

	if err := utils.WriteToFile(file, []byte(fn.CodeText)); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}
