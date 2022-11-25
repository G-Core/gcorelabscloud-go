package images

import (
	"fmt"
	"strings"

	cmeta "github.com/G-Core/gcorelabscloud-go/client/utils/metadata"

	"github.com/G-Core/gcorelabscloud-go/client/images/v1/client"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images"
	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images/types"
)

var (
	imageIDText     = "image_id is mandatory argument"
	visibilityTypes = types.Visibility("").StringList()
)

func listImages(c *cli.Context) error {
	var err error
	var cl *gcorecloud.ServiceClient
	cl, err = client.NewImageClientV1(c)
	if c.Bool("baremetal") {
		cl, err = client.NewBmImageClientV1(c)
	}
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.NewExitError(err, 1)
	}

	opts := images.ListOpts{
		Private:    c.Bool("private"),
		Visibility: types.Visibility(c.String("visibility")),
	}

	results, err := images.ListAll(cl, opts)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	utils.ShowResults(results, c.String("format"))
	return nil
}

func listProjectImages(c *cli.Context) error {
	client, err := client.NewProjectImageClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.NewExitError(err, 1)
	}
	results, err := images.ListAll(client, nil)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	utils.ShowResults(results, c.String("format"))
	return nil
}

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
		&cli.BoolFlag{
			Name:     "owner",
			Aliases:  []string{"o"},
			Usage:    "only current project images",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "baremetal",
			Aliases:  []string{"b"},
			Usage:    "only baremetal images",
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
		if c.Bool("owner") {
			return listProjectImages(c)
		}
		return listImages(c)
	},
}

var imageProjectListCommand = cli.Command{
	Name:     "list",
	Usage:    "List project images",
	Category: "image",
	Action:   listProjectImages,
}

var imageCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create image",
	Category: "image",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "image name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "volume-id",
			Aliases:  []string{"v"},
			Usage:    "Required if source is volume",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "hw-firmware-type",
			Usage:    "Specifies the type of firmware with which to boot the guest. Available value is 'bios', 'uefi'",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "hw-machine-type",
			Usage:    "A virtual chipset type. Available value is 'i440', 'q35'",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "ssh-key",
			Usage:    "Permission to use a ssh key in instances. Available value is 'allow', 'deny', 'required'",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "os-type",
			Usage:    "The operating system installed on the image. Available value is 'windows', 'linux'",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "is-baremetal",
			Usage: "Set to true if the image will be used by baremetal instances. Defaults to false.",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		downloadClient, err := client.NewDownloadImageClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := images.CreateOpts{
			Name:           c.String("name"),
			HwMachineType:  types.HwMachineType(c.String("hw-machine-type")),
			SshKey:         types.SshKeyType(c.String("ssh-key")),
			OSType:         types.OSType(c.String("os-type")),
			IsBaremetal:    utils.BoolToPointer(c.Bool("is-baremetal")),
			HwFirmwareType: types.HwFirmwareType(c.String("hw-firmware-type")),
			Source:         types.ImageSourceVolume,
			VolumeID:       c.String("volume-id"),
		}

		results, err := images.Create(downloadClient, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, downloadClient, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(downloadClient, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			instanceID, err := images.ExtractImageIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve image ID from task info: %w", err)
			}
			getClient, err := client.NewImageClientV1(c)
			if err != nil {
				return nil, err
			}
			instance, err := images.Get(getClient, instanceID).Extract()
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
		imageID, err := flags.GetFirstStringArg(c, imageIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewImageClientV1(c)
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

var imageUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update image fields",
	Category:  "image",
	ArgsUsage: "<image_id>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "image name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "hw-firmware-type",
			Usage:    "Specifies the type of firmware with which to boot the guest. Available value is 'bios', 'uefi'",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "hw-machine-type",
			Usage:    "A virtual chipset type. Available value is 'i440', 'q35'",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "ssh-key",
			Usage:    "Permission to use a ssh key in instances. Available value is 'allow', 'deny', 'required'",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "os-type",
			Usage:    "The operating system installed on the image. Available value is 'windows', 'linux'",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "is-baremetal",
			Usage: "Set to true if the image will be used by baremetal instances. Defaults to false.",
		},
	},
	Action: func(c *cli.Context) error {
		imageID, err := flags.GetFirstStringArg(c, imageIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewImageClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := images.UpdateOpts{
			Name:           c.String("name"),
			HwMachineType:  types.HwMachineType(c.String("hw-machine-type")),
			SshKey:         types.SshKeyType(c.String("ssh-key")),
			OSType:         types.OSType(c.String("os-type")),
			IsBaremetal:    utils.BoolToPointer(c.Bool("is-baremetal")),
			HwFirmwareType: types.HwFirmwareType(c.String("hw-firmware-type")),
		}

		image, err := images.Update(client, imageID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(image, c.String("format"))
		return nil
	},
}

var imageUploadCommand = cli.Command{
	Name:     "upload",
	Usage:    "Upload image",
	Category: "image",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:  "os-version",
			Usage: "OS version, i.e. 19.04 (for Ubuntu) or 9.4 for Debian.",
		},
		&cli.StringFlag{
			Name:  "os-distro",
			Usage: "OS Distribution, i.e. Debian, CentOS, Ubuntu, CoreOS etc.",
		},
		&cli.StringFlag{
			Name:     "url",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "cow-format",
			Usage: "When True, image cannot be deleted unless all volumes, created from it, are deleted. Defaults to False",
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "image name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "hw-firmware-type",
			Usage:    "Specifies the type of firmware with which to boot the guest. Available value is 'bios', 'uefi'",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "hw-machine-type",
			Usage:    "A virtual chipset type. Available value is 'i440', 'q35'",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "ssh-key",
			Usage:    "Permission to use a ssh key in instances. Available value is 'allow', 'deny', 'required'",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "os-type",
			Usage:    "The operating system installed on the image. Available value is 'windows', 'linux'",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "is-baremetal",
			Usage: "Set to true if the image will be used by baremetal instances. Defaults to false.",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		downloadClient, err := client.NewDownloadImageClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := images.UploadOpts{
			OsVersion:      c.String("os-version"),
			HwMachineType:  types.HwMachineType(c.String("hw-machine-type")),
			SshKey:         types.SshKeyType(c.String("ssh-key")),
			Name:           c.String("name"),
			OsDistro:       c.String("os-distro"),
			OSType:         types.OSType(c.String("os-type")),
			URL:            c.String("url"),
			IsBaremetal:    utils.BoolToPointer(c.Bool("is-baremetal")),
			HwFirmwareType: types.HwFirmwareType(c.String("hw-firmware-type")),
			CowFormat:      c.Bool("cow-format"),
		}

		results, err := images.Upload(downloadClient, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, downloadClient, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(downloadClient, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			instanceID, err := images.ExtractImageIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve image ID from task info: %w", err)
			}
			getClient, err := client.NewImageClientV1(c)
			if err != nil {
				return nil, err
			}
			instance, err := images.Get(getClient, instanceID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get image with ID: %s. Error: %w", instanceID, err)
			}
			return instance, nil
		})
	},
}

var imageDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete image",
	Category:  "image",
	ArgsUsage: "<image_id>",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		imageID, err := flags.GetFirstStringArg(c, imageIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewImageClientV1(c)
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

var Commands = cli.Command{
	Name:  "image",
	Usage: "GCloud images API",
	Subcommands: []*cli.Command{
		&imageListCommand,
		&imageShowCommand,
		&imageDeleteCommand,
		&imageCreateCommand,
		&imageUpdateCommand,
		&imageUploadCommand,
		{
			Name:  "project",
			Usage: "GCloud project images API",
			Subcommands: []*cli.Command{
				&imageProjectListCommand,
			},
		},
		{
			Name:  "metadata",
			Usage: "Image metadata",
			Subcommands: []*cli.Command{
				cmeta.NewMetadataListCommand(
					client.NewImageClientV1,
					"Get image metadata",
					"<image_id>",
					"image_id is mandatory argument",
				),
				cmeta.NewMetadataGetCommand(
					client.NewImageClientV1,
					"Show image metadata by key",
					"<image_id>",
					"image_id is mandatory argument",
				),
				cmeta.NewMetadataDeleteCommand(
					client.NewImageClientV1,
					"Delete image_id metadata by key",
					"<image_id>",
					"image_id is mandatory argument",
				),
				cmeta.NewMetadataCreateCommand(
					client.NewImageClientV1,
					"Create instance metadata. It would update existing keys",
					"<image_id>",
					"image_id is mandatory argument",
				),
				cmeta.NewMetadataUpdateCommand(
					client.NewImageClientV1,
					"Update image_id metadata. It overriding existing records",
					"<image_id>",
					"volume_id is mandatory argument",
				),
				cmeta.NewMetadataReplaceCommand(
					client.NewImageClientV1,
					"Replace image_id metadata. It replace existing records",
					"<image_id>",
					"image_id is mandatory argument",
				),
			},
		},
	},
}
