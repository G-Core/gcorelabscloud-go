package securitygroups

import (
	"fmt"
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/securitygroups/securitygrouprules"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/securitygroup/v1/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/securitygroup/v1/securitygroups"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"
	"github.com/urfave/cli/v2"
)

var (
	securityGroupIDText = "securitygroup_id is mandatory argument"
	protocolTypeList    = types.Protocol("").StringList()
	directionTypeList   = types.RuleDirection("").StringList()
	etherTypeTypeList   = types.EtherType("").StringList()
)

var securityGroupListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Security groups list",
	Category: "securitygroup",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "securitygroups", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := securitygroups.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var securityGroupCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create security group",
	Category: "securitygroup",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Security group name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "Security group description",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "instance-id",
			Aliases:  []string{"i"},
			Usage:    "Security group instance to attach to",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "rule-description",
			Usage:    "Security group rule description",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "rule-remote-group-id",
			Usage:    "Security group rule remote group id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "rule-remote-ip-prefix",
			Usage:    "Security group rule remote ip prefix",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "rule-port-range-max",
			Usage:    "Security group rule port max range",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "rule-port-range-min",
			Usage:    "Security group rule port min range",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "rule-protocol",
			Aliases: []string{"p"},
			Value: &utils.EnumStringSliceValue{
				Enum: protocolTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(protocolTypeList, ", ")),
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "rule-ethertype",
			Aliases: []string{"e"},
			Value: &utils.EnumStringSliceValue{
				Enum: etherTypeTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(etherTypeTypeList, ", ")),
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "rule-direction",
			Aliases: []string{"dr"},
			Value: &utils.EnumStringSliceValue{
				Enum: directionTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(directionTypeList, ", ")),
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "securitygroups", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		protocols := utils.GetEnumStringSliceValue(c, "rule-protocol")
		directions := utils.GetEnumStringSliceValue(c, "rule-direction")
		ethertypes := utils.GetEnumStringSliceValue(c, "rule-ethertype")
		descriptions := c.StringSlice("rule-description")
		remoteGroupIDs := c.StringSlice("rule-remote-group-id")
		remoteIPPrefixes := c.StringSlice("rule-remote-ip-prefix")
		portMaxRanges := c.IntSlice("rule-port-range-max")
		portMinRanges := c.IntSlice("rule-port-range-min")

		if len(protocols) != len(directions) || len(directions) != len(ethertypes) {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(fmt.Errorf("rule-protocol, rule-direction and rule-ethertype parameters number should be same"), 1)
		}

		rulesNumber := len(protocols)

		if len(descriptions) > rulesNumber {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(
				fmt.Errorf("rule-description parameters number %d is more then number of rules: %d",
					len(descriptions),
					rulesNumber,
				), 1)
		}

		if len(remoteGroupIDs) > rulesNumber {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(
				fmt.Errorf("rule-remote-group-id parameters number %d is more then number of rules: %d",
					len(remoteGroupIDs),
					rulesNumber,
				), 1)
		}

		if len(remoteIPPrefixes) > rulesNumber {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(
				fmt.Errorf("rule-remote-ip-prefix parameters number %d is more then number of rules: %d",
					len(remoteIPPrefixes),
					rulesNumber,
				), 1)
		}

		if len(portMaxRanges) > rulesNumber {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(
				fmt.Errorf("rule-port-range-max parameters number %d is more then number of rules: %d",
					len(portMaxRanges),
					rulesNumber,
				), 1)
		}

		if len(portMinRanges) > rulesNumber {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(
				fmt.Errorf("rule-port-range-min parameters number %d is more then number of rules: %d",
					len(portMinRanges),
					rulesNumber,
				), 1)
		}

		rules := make([]securitygroups.CreateSecurityGroupRuleOpts, 0)
		instances := c.StringSlice("instance-id")

		if instances == nil {
			instances = make([]string, 0)
		}

		for idx, direction := range directions {
			var portRangeMax *int
			var portRangeMin *int
			var description *string
			var ipPrefix *string
			var groupID *string
			if idx < len(portMaxRanges) {
				portRangeMax = &portMaxRanges[idx]
			}
			if idx < len(portMinRanges) {
				portRangeMin = &portMinRanges[idx]
			}
			if idx < len(descriptions) {
				description = &descriptions[idx]
			}
			if idx < len(remoteIPPrefixes) {
				ipPrefix = &remoteIPPrefixes[idx]
			}
			if idx < len(remoteGroupIDs) {
				groupID = &remoteGroupIDs[idx]
			}
			rule := securitygroups.CreateSecurityGroupRuleOpts{
				Direction:      types.RuleDirection(direction),
				RemoteGroupID:  groupID,
				EtherType:      types.EtherType(ethertypes[idx]),
				Protocol:       types.Protocol(protocols[idx]),
				PortRangeMax:   portRangeMax,
				PortRangeMin:   portRangeMin,
				Description:    description,
				RemoteIPPrefix: ipPrefix,
			}
			rules = append(rules, rule)
		}

		opts := securitygroups.CreateOpts{
			SecurityGroup: securitygroups.CreateSecurityGroupOpts{
				Name:               c.String("name"),
				Description:        utils.StringToPointer(c.String("description")),
				SecurityGroupRules: rules,
			},
			Instances: instances,
		}

		results, err := securitygroups.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var securityGroupGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show securitygroup",
	ArgsUsage: "<securitygroup_id>",
	Category:  "securitygroup",
	Action: func(c *cli.Context) error {
		securityGroupID, err := flags.GetFirstArg(c, securityGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "securitygroups", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := securitygroups.Get(client, securityGroupID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var securityGroupListInstancesSubCommand = cli.Command{
	Name:      "list",
	Usage:     "securitygroup group instances list",
	ArgsUsage: "<securitygroup_id>",
	Category:  "securitygroup",
	Action: func(c *cli.Context) error {
		securityGroupID, err := flags.GetFirstArg(c, securityGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}
		client, err := utils.BuildClient(c, "securitygroups", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := securitygroups.ListAllInstances(client, securityGroupID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var securityGroupDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete security group",
	ArgsUsage: "<securitygroup_id>",
	Category:  "securitygroup",
	Action: func(c *cli.Context) error {
		securityGroupID, err := flags.GetFirstArg(c, securityGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "securitygroups", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = securitygroups.Delete(client, securityGroupID).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var securityGroupUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Update security group",
	ArgsUsage: "<securitygroup_id>",
	Category:  "securitygroup",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Security name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		securityGroupID, err := flags.GetFirstArg(c, securityGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := utils.BuildClient(c, "securitygroups", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := securitygroups.UpdateOpts{Name: c.String("name")}

		result, err := securitygroups.Update(client, securityGroupID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if result == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var SecurityGroupCommands = cli.Command{
	Name:  "securitygroup",
	Usage: "GCloud security groups API",
	Subcommands: []*cli.Command{
		&securityGroupListSubCommand,
		&securityGroupGetSubCommand,
		&securityGroupUpdateSubCommand,
		&securityGroupDeleteSubCommand,
		&securityGroupCreateSubCommand,
		{
			Name:  "instance",
			Usage: "Security group instances",
			Subcommands: []*cli.Command{
				&securityGroupListInstancesSubCommand,
			},
		},
		{
			Name:        "rule",
			Usage:       "Security group rules",
			Subcommands: securitygrouprules.SecurityGroupRuleCommands,
		},
	},
}
