package projects

import (
	"fmt"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/projects/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/project/v1/projects"
	"github.com/G-Core/gcorelabscloud-go/gcore/project/v1/types"
	"github.com/urfave/cli/v2"
)

var (
	projectIDText     = "project_id is mandatory argument"
	projectStatesList = types.ProjectState("").StringList()
)

var projectListCommand = cli.Command{
	Name:     "list",
	Usage:    "List projects",
	Category: "project",
	Action: func(c *cli.Context) error {
		client, err := client.NewProjectClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := projects.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var projectGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get project",
	ArgsUsage: "<project_id>",
	Category:  "project",
	Action: func(c *cli.Context) error {
		projectID, err := flags.GetFirstIntArg(c, projectIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewProjectClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := projects.Get(client, projectID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(task, c.String("format"))
		return nil
	},
}

var projectUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update project",
	ArgsUsage: "<project_id>",
	Category:  "project",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "project name",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "project description",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		projectID, err := flags.GetFirstIntArg(c, projectIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}

		opts := projects.UpdateOpts{
			Name:        c.String("name"),
			Description: c.String("description"),
		}

		err = gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		client, err := client.NewProjectClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := projects.Update(client, projectID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var projectDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete project",
	ArgsUsage: "<project_id>",
	Category:  "project",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		projectID, err := flags.GetFirstIntArg(c, projectIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}

		client, err := client.NewProjectClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		results, err := projects.Delete(client, projectID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := projects.Get(client, projectID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete project with ID: %d", projectID)
			}
			switch err.(type) {
			case gcorecloud.ErrDefault404:
				return nil, nil
			default:
				return nil, err
			}
		})

	},
}

var projectCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create project",
	Category: "project",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "project name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "project description",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "state",
			Aliases: []string{"s"},
			Value: &utils.EnumValue{
				Enum: projectStatesList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(projectStatesList, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "client-id",
			Usage:    "client id",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {

		opts := projects.CreateOpts{
			ClientID:    c.Int("client-id"),
			State:       types.ProjectState(c.String("state")),
			Name:        c.String("name"),
			Description: c.String("description"),
		}

		err := gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		client, err := client.NewProjectClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := projects.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var ProjectCommands = cli.Command{
	Name:  "project",
	Usage: "GCloud projects API",
	Subcommands: []*cli.Command{
		&projectListCommand,
		&projectGetCommand,
		&projectDeleteCommand,
		&projectUpdateCommand,
		&projectCreateCommand,
	},
}
