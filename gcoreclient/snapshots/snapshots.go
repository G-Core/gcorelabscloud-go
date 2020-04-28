package snapshots

import (
	"fmt"

	"github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/snapshot/v1/snapshots"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/flags"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/utils"

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
		client, err := utils.BuildClient(c, "snapshots", "", "")
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
		client, err := utils.BuildClient(c, "snapshots", "", "")
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
		client, err := utils.BuildClient(c, "snapshots", "", "")
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
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "snapshots", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := snapshots.CreateOpts{
			VolumeID:    c.String("volume-id"),
			Name:        c.String("name"),
			Description: c.String("description"),
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

var SnapshotCommands = cli.Command{
	Name:  "snapshot",
	Usage: "GCloud snapshots API",
	Subcommands: []*cli.Command{
		&snapshotListCommand,
		&snapshotGetCommand,
		&snapshotDeleteCommand,
		&snapshotCreateCommand,
	},
}
