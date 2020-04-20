package floatingips

import (
	"fmt"
	"net"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/floatingips/availablefloatingips"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/volume/v1/volumes"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/floatingip/v1/floatingips"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"
	"github.com/urfave/cli/v2"
)

var (
	floatingIPIDText = "floatingip_id is mandatory argument"
)

var floatingIPListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Floating ips list",
	Category: "floatingip",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "floatingips", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := floatingips.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var floatingIPCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create floating ip",
	Category: "floatingip",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "port-id",
			Aliases:  []string{"p"},
			Usage:    "port id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "fixed-ip-address",
			Aliases:  []string{"a"},
			Usage:    "fixed ip address",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "floatingips", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		ip := net.ParseIP(c.String("fixed-ip-address"))

		if ip == nil {
			_ = cli.ShowCommandHelp(c, "show")
			return fmt.Errorf("malformer ip address: %s", c.String("fixed-ip-address"))
		}

		opts := floatingips.CreateOpts{
			PortID:         c.String("port-id"),
			FixedIPAddress: ip,
		}

		results, err := floatingips.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			floatingIPID, err := floatingips.ExtractFloatingIPIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve floating IP ID from task info: %w", err)
			}
			volume, err := floatingips.Get(client, floatingIPID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get floating IP ID: %s. Error: %w", floatingIPID, err)
			}
			utils.ShowResults(volume, c.String("format"))
			return nil, nil
		})

	},
}

var floatingIPGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show floatingip",
	ArgsUsage: "<floatingip_id>",
	Category:  "floatingip",
	Action: func(c *cli.Context) error {
		floatingIPID, err := flags.GetFirstStringArg(c, floatingIPIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "floatingips", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := floatingips.Get(client, floatingIPID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var floatingIPDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete floating ip",
	ArgsUsage: "<floatingip_id>",
	Category:  "floatingip",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		floatingIPID, err := flags.GetFirstStringArg(c, floatingIPIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "floatingips", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := floatingips.Delete(client, floatingIPID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := volumes.Get(client, floatingIPID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete floating IP with ID: %s", floatingIPID)
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

var floatingIPAssignSubCommand = cli.Command{
	Name:      "assign",
	Usage:     "Update floating ip",
	ArgsUsage: "<floatingip_id>",
	Category:  "floatingip",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "port-id",
			Aliases:  []string{"p"},
			Usage:    "port id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "fixed-ip-address",
			Aliases:  []string{"a"},
			Usage:    "fixed ip address",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		floatingIPID, err := flags.GetFirstStringArg(c, floatingIPIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := utils.BuildClient(c, "floatingips", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		ip := net.ParseIP(c.String("fixed-ip-address"))

		if ip == nil {
			_ = cli.ShowCommandHelp(c, "show")
			return fmt.Errorf("malformer ip address: %s", c.String("fixed-ip-address"))
		}

		opts := floatingips.CreateOpts{
			PortID:         c.String("port-id"),
			FixedIPAddress: ip,
		}

		floatingIP, err := floatingips.Assign(client, floatingIPID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(floatingIP, c.String("format"))
		return nil
	},
}

var floatingIPUnAssignSubCommand = cli.Command{
	Name:      "unassign",
	Usage:     "Update floating ip",
	ArgsUsage: "<floatingip_id>",
	Category:  "floatingip",
	Action: func(c *cli.Context) error {
		floatingIPID, err := flags.GetFirstStringArg(c, floatingIPIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := utils.BuildClient(c, "floatingips", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		floatingIP, err := floatingips.UnAssign(client, floatingIPID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(floatingIP, c.String("format"))
		return nil
	},
}

var FloatingIPCommands = cli.Command{
	Name:  "floatingip",
	Usage: "GCloud floating ips API",
	Subcommands: []*cli.Command{
		&floatingIPListSubCommand,
		&floatingIPGetSubCommand,
		&floatingIPAssignSubCommand,
		&floatingIPUnAssignSubCommand,
		&floatingIPDeleteSubCommand,
		&floatingIPCreateSubCommand,
		&availablefloatingips.AvailableFloatingIPCommands,
	},
}
