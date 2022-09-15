package volumes

import (
	"fmt"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	cmeta "github.com/G-Core/gcorelabscloud-go/client/utils/metadata"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/client/volumes/v1/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/urfave/cli/v2"
)

var (
	volumeIDText      = "volume_id is mandatory argument"
	volumeSourceNames = volumes.VolumeSource("").StringList()
	volumeTypeNames   = volumes.VolumeType("").StringList()
)

var volumeListCommand = cli.Command{
	Name:     "list",
	Usage:    "List volumes",
	Category: "volume",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "instance-id",
			Aliases:     []string{"i"},
			Usage:       "Instance ID",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "cluster-id",
			Aliases:     []string{"c"},
			Usage:       "Cluster ID, use to get pvc",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "id-part",
			Usage:       "Filter the volume list result by the ID part of the volume",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "name-part",
			Usage:       "Filter out volumes by name_part inclusion in volume name",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.BoolFlag{
			Name:     "bootable",
			Usage:    "Filter by a bootable field",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "has-attachments",
			Usage:    "Filter by the presence of attachments",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewVolumeClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := volumes.ListOpts{
			InstanceID:     utils.StringToPointer(c.String("instance-id")),
			ClusterID:      utils.StringToPointer(c.String("cluster-id")),
			IDPart:         utils.StringToPointer(c.String("id-part")),
			NamePart:       utils.StringToPointer(c.String("name-part")),
			Bootable:       utils.BoolToPointer(c.Bool("bootable")),
			HasAttachments: utils.BoolToPointer(c.Bool("has-attachments")),
		}
		results, err := volumes.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var volumeGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get volume information",
	ArgsUsage: "<volume_id>",
	Category:  "volume",
	Action: func(c *cli.Context) error {
		volumeID, err := flags.GetFirstStringArg(c, volumeIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewVolumeClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := volumes.Get(client, volumeID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if task == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(task, c.String("format"))
		return nil
	},
}

var volumeDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete volume by ID",
	ArgsUsage: "<volume_id>",
	Category:  "volume",
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:        "snapshot",
			Aliases:     []string{"s"},
			Usage:       "Shapshots to delete",
			DefaultText: "nil",
			Required:    false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		volumeID, err := flags.GetFirstStringArg(c, volumeIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewVolumeClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := volumes.DeleteOpts{
			Snapshots: c.StringSlice("snapshot"),
		}
		results, err := volumes.Delete(client, volumeID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := volumes.Get(client, volumeID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete volume with ID: %s", volumeID)
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

var volumeCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create volume",
	Category: "volume",
	Flags: append([]cli.Flag{
		&cli.GenericFlag{
			Name:    "source",
			Aliases: []string{"s"},
			Value: &utils.EnumValue{
				Enum:    volumeSourceNames,
				Default: volumeSourceNames[0],
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(volumeSourceNames, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Volume name",
			Required: true,
		},
		&cli.IntFlag{
			Name:        "size",
			Usage:       "Volume size",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.GenericFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Value: &utils.EnumValue{
				Enum:    volumeTypeNames,
				Default: volumeTypeNames[0],
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(volumeTypeNames, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:        "image-id",
			Aliases:     []string{"i"},
			Usage:       "Image ID",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "snapshot-id",
			Usage:       "Snapshot ID",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "instance-id",
			Usage:       "Instance ID to attach",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.IntSliceFlag{
			Name:     "lifecycle-policy-ids",
			Usage:    "Lifecycle policy ids list",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := client.NewVolumeClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := volumes.CreateOpts{
			Source:               volumes.VolumeSource(c.String("source")),
			Name:                 c.String("name"),
			Size:                 c.Int("size"),
			TypeName:             volumes.VolumeType(c.String("type")),
			ImageID:              c.String("image-id"),
			SnapshotID:           c.String("snapshot-id"),
			InstanceIDToAttachTo: c.String("instance-id"),
		}
		lfPid := c.IntSlice("lifecycle-policy-ids")
		if len(lfPid) != 0 {
			opts.LifeCyclePolicyIDs = lfPid
		}

		results, err := volumes.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			volumeID, err := volumes.ExtractVolumeIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve volume ID from task info: %w", err)
			}
			volume, err := volumes.Get(client, volumeID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get volume with ID: %s. Error: %w", volumeID, err)
			}
			return volume, nil
		})
	},
}

var volumeAttachCommand = cli.Command{
	Name:      "attach",
	Usage:     "Attach volume to instance",
	ArgsUsage: "<volume_id>",
	Category:  "volume",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "instance-id",
			Aliases:  []string{"i"},
			Usage:    "Instance ID to attach",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		volumeID, err := flags.GetFirstStringArg(c, volumeIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "attach")
			return err
		}
		client, err := client.NewVolumeClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := volumes.InstanceOperationOpts{
			InstanceID: c.String("instance-id"),
		}
		volume, err := volumes.Attach(client, volumeID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(volume, c.String("format"))
		return nil
	},
}

var volumeDetachCommand = cli.Command{
	Name:      "detach",
	Usage:     "Detach volume to instance",
	ArgsUsage: "<volume_id>",
	Category:  "volume",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "instance-id",
			Aliases:  []string{"i"},
			Usage:    "Instance ID to attach",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		volumeID, err := flags.GetFirstStringArg(c, volumeIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "detach")
			return err
		}
		client, err := client.NewVolumeClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := volumes.InstanceOperationOpts{
			InstanceID: c.String("instance-id"),
		}
		volume, err := volumes.Detach(client, volumeID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(volume, c.String("format"))
		return nil
	},
}

var volumeRetypeCommand = cli.Command{
	Name:      "retype",
	Usage:     "Change volume type",
	ArgsUsage: "<volume_id>",
	Category:  "volume",
	Flags: []cli.Flag{
		&cli.GenericFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Value: &utils.EnumValue{
				Enum:    volumeTypeNames,
				Default: volumeTypeNames[0],
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(volumeTypeNames, ", ")),
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		volumeID, err := flags.GetFirstStringArg(c, volumeIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "retype")
			return err
		}
		client, err := client.NewVolumeClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := volumes.VolumeTypePropertyOperationOpts{
			VolumeType: volumes.VolumeType(c.String("type")),
		}
		volume, err := volumes.Retype(client, volumeID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(volume, c.String("format"))
		return nil
	},
}

var volumeExtendCommand = cli.Command{
	Name:      "extend",
	Usage:     "Change volume size",
	ArgsUsage: "<volume_id>",
	Category:  "volume",
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:     "size",
			Aliases:  []string{"s"},
			Usage:    "Volume size",
			Required: true,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		volumeID, err := flags.GetFirstStringArg(c, volumeIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "extend")
			return err
		}
		client, err := client.NewVolumeClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		size := c.Int("size")
		opts := volumes.SizePropertyOperationOpts{
			Size: size,
		}
		results, err := volumes.Extend(client, volumeID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			volume, err := volumes.Get(client, volumeID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get volume with ID: %s. Error: %w", volumeID, err)
			}
			utils.ShowResults(volume, c.String("format"))
			return nil, nil
		})
	},
}

var volumeRevertCommand = cli.Command{
	Name:      "revert",
	Usage:     "Revert volume to it's last snapshot",
	ArgsUsage: "<volume_id>",
	Category:  "volume",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		volumeID, err := flags.GetFirstStringArg(c, volumeIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "revert")
			return err
		}
		client, err := client.NewVolumeClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := volumes.Revert(client, volumeID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			volume, err := volumes.Get(client, volumeID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get volume with ID: %s. Error: %w", volumeID, err)
			}
			utils.ShowResults(volume, c.String("format"))
			return nil, nil
		})
	},
}

var Commands = cli.Command{
	Name:  "volume",
	Usage: "GCloud volumes API",
	Subcommands: []*cli.Command{
		&volumeListCommand,
		&volumeGetCommand,
		&volumeDeleteCommand,
		&volumeCreateCommand,
		&volumeAttachCommand,
		&volumeDetachCommand,
		&volumeRetypeCommand,
		&volumeExtendCommand,
		&volumeRevertCommand,
		{
			Name:  "metadata",
			Usage: "Volume metadata",
			Subcommands: []*cli.Command{
				cmeta.NewMetadataListCommand(
					client.NewVolumeClientV1,
					"Get volume metadata",
					"<volume_id>",
					"volume_id is mandatory argument",
				),
				cmeta.NewMetadataGetCommand(
					client.NewVolumeClientV1,
					"Show volume metadata by key",
					"<volume_id>",
					"volume_id is mandatory argument",
				),
				cmeta.NewMetadataDeleteCommand(
					client.NewVolumeClientV1,
					"Delete volume metadata by key",
					"<volume_id>",
					"volume_id is mandatory argument",
				),
				cmeta.NewMetadataCreateCommand(
					client.NewVolumeClientV1,
					"Create instance metadata. It would update existing keys",
					"<volume_id>",
					"volume_id is mandatory argument",
				),
				cmeta.NewMetadataUpdateCommand(
					client.NewVolumeClientV1,
					"Update volume metadata. It overriding existing records",
					"<volume_id>",
					"volume_id is mandatory argument",
				),
				cmeta.NewMetadataReplaceCommand(
					client.NewVolumeClientV1,
					"Replace volume metadata. It replace existing records",
					"<volume_id>",
					"volume_id is mandatory argument",
				),
			},
		},
	},
}
