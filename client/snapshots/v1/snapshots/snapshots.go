package snapshots

import (
	"errors"
	"fmt"

	"github.com/G-Core/gcorelabscloud-go/client/snapshots/v1/client"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/snapshot/v1/snapshots"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"
)

var snapshotIDText = "snapshot_id is mandatory argument"

var snapshotListCommand = cli.Command{
	Name:     "list",
	Usage:    "List snapshots",
	Category: "snapshot",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "volume-id",
			Aliases:  []string{"v"},
			Usage:    "shapshot volume ID",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "instance-id",
			Aliases:  []string{"i"},
			Usage:    "shapshot instance ID",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewSnapshotClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := snapshots.ListOpts{
			VolumeID:   c.String("volume-id"),
			InstanceID: c.String("instance-id"),
		}
		results, err := snapshots.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var snapshotGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get snapshot information",
	ArgsUsage: "<snapshot_id>",
	Category:  "snapshot",
	Action: func(c *cli.Context) error {
		snapshotID, err := flags.GetFirstStringArg(c, snapshotIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewSnapshotClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		snapshot, err := snapshots.Get(client, snapshotID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if snapshot == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(snapshot, c.String("format"))
		return nil
	},
}

var snapshotDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete snapshot by ID",
	ArgsUsage: "<snapshot_id>",
	Category:  "snapshot",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		snapshotID, err := flags.GetFirstStringArg(c, snapshotIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewSnapshotClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := snapshots.Delete(client, snapshotID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := snapshots.Get(client, snapshotID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete snapshot with ID: %s", snapshotID)
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

var snapshotCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create snapshot",
	Category: "snapshot",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "snapshot name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "volume-id",
			Aliases:  []string{"v"},
			Usage:    "volume ID",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "snapshot description",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "meta-key",
			Usage:    "metadata key, use with meta-value",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "meta-value",
			Usage:    "metadata key, use with meta-value",
			Required: false,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := client.NewSnapshotClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		meta, err := parseMetadata(c.StringSlice("meta-key"), c.StringSlice("meta-value"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		opts := snapshots.CreateOpts{
			VolumeID:    c.String("volume-id"),
			Name:        c.String("name"),
			Description: c.String("description"),
			Metadata:    meta,
		}

		results, err := snapshots.Create(client, opts).Extract()
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
			snapshotID, err := snapshots.ExtractSnapshotIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve snapshot ID from task info: %w", err)
			}
			snapshot, err := snapshots.Get(client, snapshotID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get snapshot with ID: %s. Error: %w", snapshotID, err)
			}
			utils.ShowResults(snapshot, c.String("format"))
			return nil, nil
		})
	},
}

func parseMetadata(keys []string, values []string) (map[string]string, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	if len(keys) != len(values) {
		return nil, errors.New("meta-key args count should equal meta-value args count")
	}

	meta := map[string]string{}
	for i, _ := range keys {
		meta[keys[i]] = values[i]
	}
	return meta, nil
}

var metadataReplaceCommand = cli.Command{
	Name:      "replace",
	Usage:     "Replace snapshot metadata by key",
	ArgsUsage: "<snapshot_id>",
	Category:  "metadata",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "meta-key",
			Usage: "metadata key, use with meta-value",
		},
		&cli.StringSliceFlag{
			Name:  "meta-value",
			Usage: "metadata key, use with meta-value",
		},
	},
	Action: func(c *cli.Context) error {
		snapshotID, err := flags.GetFirstStringArg(c, snapshotIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "replace")
			return err
		}

		client, err := client.NewSnapshotClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		meta, err := parseMetadata(c.StringSlice("meta-key"), c.StringSlice("meta-value"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		metaData := make([]snapshots.MetadataOpts, 0, len(meta))
		for k, v := range meta {
			metaData = append(metaData, snapshots.MetadataOpts{Key: k, Value: v})
		}
		opts := snapshots.MetadataSetOpts{Metadata: metaData}

		snapshot, err := snapshots.MetadataReplace(client, snapshotID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(snapshot, c.String("format"))
		return nil
	},
}

var Commands = cli.Command{
	Name:  "snapshot",
	Usage: "GCloud snapshots API",
	Subcommands: []*cli.Command{
		&snapshotListCommand,
		&snapshotGetCommand,
		&snapshotDeleteCommand,
		&snapshotCreateCommand,
		{
			Name:  "metadata",
			Usage: "Snapshot metadata",
			Subcommands: []*cli.Command{
				&metadataReplaceCommand,
			},
		},
	},
}
