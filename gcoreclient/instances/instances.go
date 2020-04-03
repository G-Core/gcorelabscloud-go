package instances

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/instance/v1/instances"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"
	"github.com/urfave/cli/v2"
)

var instanceIDText = "instance_id is mandatory argument"

var instanceListCommand = cli.Command{
	Name:     "list",
	Usage:    "List instances",
	Category: "instance",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "exclude-secgroup",
			Aliases:  []string{"e"},
			Usage:    "Exclude instances with specified security group name",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "available-floating",
			Aliases:  []string{"a"},
			Usage:    "Only show instances which are able to handle floating address",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "instances", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		af := c.Bool("available-floating")
		var availableFloating string
		if af {
			availableFloating = "available-floating"
		}
		opts := instances.ListOpts{
			ExcludeSecGroup:   utils.StringToPointer(c.String("exclude-secgroup")),
			AvailableFloating: utils.StringToPointer(availableFloating),
		}
		results, err := instances.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var instanceListInterfacesCommand = cli.Command{
	Name:      "list",
	Usage:     "List instance interfaces",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}
		client, err := utils.BuildClient(c, "instances", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := instances.ListInterfacesAll(client, instanceID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var instanceListSecurityGroupsCommand = cli.Command{
	Name:      "list",
	Usage:     "List instance security groups",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}
		client, err := utils.BuildClient(c, "instances", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := instances.ListSecurityGroupsAll(client, instanceID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var instanceAssignSecurityGroupsCommand = cli.Command{
	Name:      "add",
	Usage:     "Add instance security group",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "security group name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "add")
			return err
		}
		client, err := utils.BuildClient(c, "instances", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := instances.SecurityGroupOpts{Name: c.String("name")}

		err = instances.AssignSecurityGroup(client, instanceID, opts).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var instanceUnAssignSecurityGroupsCommand = cli.Command{
	Name:      "delete",
	Usage:     "Add instance security group",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "security group name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "instances", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := instances.SecurityGroupOpts{Name: c.String("name")}

		err = instances.UnAssignSecurityGroup(client, instanceID, opts).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var instanceGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get instance information",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "instances", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(task, c.String("format"))
		return nil
	},
}

var InstanceCommands = cli.Command{
	Name:  "instance",
	Usage: "GCloud instances API",
	Subcommands: []*cli.Command{
		&instanceGetCommand,
		&instanceListCommand,
		{
			Name:  "interface",
			Usage: "Instance interfaces",
			Subcommands: []*cli.Command{
				&instanceListInterfacesCommand,
			},
		},
		{
			Name:  "securitygroup",
			Usage: "Instance security groups",
			Subcommands: []*cli.Command{
				&instanceListSecurityGroupsCommand,
				&instanceAssignSecurityGroupsCommand,
				&instanceUnAssignSecurityGroupsCommand,
			},
		},
	},
}
