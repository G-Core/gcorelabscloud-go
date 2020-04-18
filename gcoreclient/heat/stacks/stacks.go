package stacks

import (
	"fmt"
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/stacks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/stacks/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var (
	stackIDText = "stack_id is mandatory argument"
	sortKeyList = types.SortKey("").StringList()
	sortDirList = types.SortDir("").StringList()
)

var stackListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Heat stacks list",
	Category: "stack",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Usage:    "stack id",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "status",
			Usage:    "stack status",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "stack name",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "tags",
			Usage:    "filter stack by tags. comma separated. AND boolean",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "tags-any",
			Usage:    "filter stack by tags. comma separated. OR boolean",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "not-tags",
			Usage:    "filter stack by tags. excluding. comma separated. AND boolean",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "not-tags-any",
			Usage:    "filter stack by tags. excluding. comma separated. OR boolean",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "show-deleted",
			Usage:    "show deleted stacks",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "show-nested",
			Usage:    "show nested stacks",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "all-tenants",
			Usage:    "show stacks for all tenants",
			Required: false,
		},
		&cli.GenericFlag{
			Name: "sort-key",
			Value: &utils.EnumValue{
				Enum: sortKeyList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(sortKeyList, ", ")),
			Required: false,
		},
		&cli.GenericFlag{
			Name: "sort-dir",
			Value: &utils.EnumValue{
				Enum: sortDirList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(sortDirList, ", ")),
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "heat", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := stacks.ListOpts{
			ID:          c.String("id"),
			Status:      c.String("status"),
			Name:        c.String("name"),
			SortKey:     types.SortKey(c.String("sort-key")),
			SortDir:     types.SortDir(c.String("sort-dir")),
			AllTenants:  c.Bool("show-tenants"),
			ShowDeleted: c.Bool("show-deleted"),
			ShowNested:  c.Bool("show-nested"),
			Tags:        c.String("tags"),
			TagsAny:     c.String("tags-any"),
			NotTags:     c.String("not-tags"),
			NotTagsAny:  c.String("not-tags-any"),
		}

		results, err := stacks.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var stackGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show heat stack",
	ArgsUsage: "<stack_id>",
	Category:  "stack",
	Action: func(c *cli.Context) error {
		stackID, err := flags.GetFirstStringArg(c, stackIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "heat", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := stacks.Get(client, stackID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var stackUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Update heat stack",
	ArgsUsage: "<stack_id>",
	Category:  "stack",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "template",
			Usage:    "stack template yaml file",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "environment",
			Usage:    "stack environment yaml file",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "parameter",
			Usage:    "stack parameters. Example: --parameter one=two --parameter three=four",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "patch",
			Usage:    "use path method. template is not mandatory",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		stackID, err := flags.GetFirstStringArg(c, stackIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := utils.BuildClient(c, "heat", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := stacks.UpdateOpts{}
		templateFile := c.String("template")
		environmentFile := c.String("environment")

		if templateFile != "" {
			content, err := utils.CheckYamlFile(templateFile)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			template := &stacks.Template{}
			template.TE = stacks.TE{
				Bin: content,
			}
			opts.TemplateOpts = template
		}

		if environmentFile != "" {
			content, err := utils.CheckYamlFile(environmentFile)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			env := &stacks.Environment{}
			env.TE = stacks.TE{
				Bin: content,
			}
			opts.EnvironmentOpts = env
		}

		params, err := utils.StringSliceToMapInterface(c.StringSlice("parameter"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		opts.Parameters = params

		if c.Bool("patch") {
			err = stacks.UpdatePatch(client, stackID, opts).ExtractErr()
		} else {
			err = stacks.Update(client, stackID, opts).ExtractErr()
		}
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var StackCommands = cli.Command{
	Name:  "stack",
	Usage: "Heat stacks commands",
	Subcommands: []*cli.Command{
		&stackGetSubCommand,
		&stackListSubCommand,
		&stackUpdateSubCommand,
	},
}
