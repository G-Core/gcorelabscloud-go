package images

import (
	"fmt"
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/image/v1/images"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/image/v1/images/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"
	"github.com/urfave/cli/v2"
)

var (
	imageIDText     = "image_id is mandatory argument"
	visibilityTypes = types.Visibility("").StringList()
)

var imageListCommand = cli.Command{
	Name:     "list",
	Usage:    "List images",
	Category: "image",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "private",
			Aliases:  []string{"p"},
			Usage:    "only private images",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "visibility",
			Aliases: []string{"v"},
			Value: &utils.EnumValue{
				Enum: visibilityTypes,
			},
			Usage:    fmt.Sprintf("image visibility type. output in %s", strings.Join(visibilityTypes, ", ")),
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "images", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := images.ListOpts{
			Private:    c.Bool("private"),
			Visibility: types.Visibility(c.String("visibility")),
		}

		results, err := images.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var imageCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create image",
	Category: "image",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Aliases:  []string{"u"},
			Usage:    "image url",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "image name",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "cow-format",
			Aliases:  []string{"c"},
			Usage:    "image with cow format",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:        "property",
			Usage:       "Image properties. Example: --property os_distro=coreos",
			DefaultText: "nil",
			Required:    false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "downloadimage", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		properties, err := utils.StringSliceToMapNil(c.StringSlice("property"))
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := images.CreateOpts{
			URL:        c.String("url"),
			Name:       c.String("name"),
			CowFormat:  c.Bool("cow-format"),
			Properties: properties,
		}

		results, err := images.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			instanceID, err := images.ExtractImageIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve image ID from task info: %w", err)
			}
			instance, err := images.Get(client, instanceID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get image with ID: %s. Error: %w", instanceID, err)
			}
			return instance, nil
		})
	},
}

var imageShowCommand = cli.Command{
	Name:      "show",
	Usage:     "Show image details",
	Category:  "image",
	ArgsUsage: "<image_id>",
	Action: func(c *cli.Context) error {
		imageID, err := flags.GetFirstArg(c, imageIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "images", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		image, err := images.Get(client, imageID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(image, c.String("format"))
		return nil
	},
}

var imageDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete image",
	Category:  "image",
	ArgsUsage: "<image_id>",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		imageID, err := flags.GetFirstArg(c, imageIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "images", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		results, err := images.Delete(client, imageID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := images.Get(client, imageID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete image with ID: %s", imageID)
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

var ImageCommands = cli.Command{
	Name:  "image",
	Usage: "GCloud images API",
	Subcommands: []*cli.Command{
		&imageListCommand,
		&imageShowCommand,
		&imageDeleteCommand,
		&imageCreateCommand,
	},
}
