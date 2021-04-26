package lifecyclepolicy

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/lifecyclepolicy/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	category          = "lifecycle_policy"
	argsUsagePolicyID = "<lifecycle_policy_id>"
	idErrorText       = "lifecycle_policy_id is mandatory argument"
)

var (
	policyTypes = lifecyclepolicy.ScheduleType("").StringList()
)

var Commands = cli.Command{
	Name:  "lifecycle",
	Usage: "GCloud lifecycle policy API",
	Subcommands: []*cli.Command{ // TODO: quota
		&getSubCommand,
		&listSubCommand,
		&deleteSubCommand,
		&createSubCommand,
		&updateSubCommand,
		&volumeSubCommands,
		&scheduleSubCommands,
	},
}

var volumeSubCommands = cli.Command{
	Name:  "volume",
	Usage: "GCloud lifecycle policy volume API",
	Subcommands: []*cli.Command{
		&addVolumeSubCommand,
		&removeVolumeSubCommand,
	},
}

var scheduleSubCommands = cli.Command{
	Name:  "schedule",
	Usage: "GCloud lifecycle policy schedule API",
	Subcommands: []*cli.Command{
		&addScheduleSubCommand,
		&removeScheduleSubCommand,
	},
}

var volumesFlag = cli.BoolFlag{
	Name:     "volumes",
	Aliases:  []string{"v"},
	Usage:    "Set if you need volume ids",
	Required: false,
}

var getSubCommand = cli.Command{
	Name:      "show",
	Usage:     "show lifecycle policy",
	ArgsUsage: argsUsagePolicyID,
	Category:  category,
	Flags: []cli.Flag{
		&volumesFlag,
	},
	Action: func(c *cli.Context) error {
		lifecyclePolicyID, err := flags.GetFirstIntArg(c, idErrorText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewLifecyclePolicyClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := lifecyclepolicy.GetOpts{
			NeedVolumes: c.Bool("volumes"),
		}
		result, err := lifecyclepolicy.Get(client, lifecyclePolicyID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var listSubCommand = cli.Command{
	Name:     "list",
	Usage:    "list lifecycle policies",
	Category: category,
	Flags: []cli.Flag{
		&volumesFlag,
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewLifecyclePolicyClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := lifecyclepolicy.ListOpts{
			NeedVolumes: c.Bool("volumes"),
		}
		result, err := lifecyclepolicy.List(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var deleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "delete lifecycle policy",
	ArgsUsage: argsUsagePolicyID,
	Category:  category,
	Action: func(c *cli.Context) error {
		lifecyclePolicyID, err := flags.GetFirstIntArg(c, idErrorText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewLifecyclePolicyClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = lifecyclepolicy.Delete(client, lifecyclePolicyID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var createSubCommand = cli.Command{
	Name:     "create",
	Usage:    "create lifecycle policy",
	Category: category,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Policy name",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "paused",
			Aliases:  []string{"p"},
			Usage:    "Set if you want policy paused",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewLifecyclePolicyClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		var status lifecyclepolicy.PolicyStatus
		if c.Bool("paused") {
			status = lifecyclepolicy.PolicyStatusPaused
		} else {
			status = lifecyclepolicy.PolicyStatusActive
		}
		opts := lifecyclepolicy.CreateOpts{
			Name:   c.String("name"),
			Status: status,
			Action: lifecyclepolicy.PolicyActionVolumeSnapshot,
		}
		result, err := lifecyclepolicy.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var updateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "update lifecycle policy",
	ArgsUsage: argsUsagePolicyID,
	Category:  category,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "New name for policy",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "pause",
			Aliases:  []string{"p"},
			Usage:    "Set if you want policy paused",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "unpause",
			Aliases:  []string{"u"},
			Usage:    "Set if you want policy active",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		lifecyclePolicyID, err := flags.GetFirstIntArg(c, idErrorText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := client.NewLifecyclePolicyClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		var status lifecyclepolicy.PolicyStatus
		if c.Bool("pause") {
			if c.Bool("unpause") {
				_ = cli.ShowCommandHelp(c, "update")
				return fmt.Errorf("at most one flag of 'pause' and 'unpause' should be set")
			}
			status = lifecyclepolicy.PolicyStatusPaused
		} else if c.Bool("unpause") {
			status = lifecyclepolicy.PolicyStatusActive
		}
		opts := lifecyclepolicy.UpdateOpts{
			Name:   c.String("name"),
			Status: status,
		}
		result, err := lifecyclepolicy.Update(client, lifecyclePolicyID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var addVolumeSubCommand = cli.Command{
	Name:      "add",
	Usage:     "add volume to lifecycle policy",
	ArgsUsage: argsUsagePolicyID,
	Category:  category,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "volume",
			Aliases:  []string{"v"},
			Usage:    "Volume ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		lifecyclePolicyID, err := flags.GetFirstIntArg(c, idErrorText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "add")
			return err
		}
		client, err := client.NewLifecyclePolicyClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := lifecyclepolicy.AddVolumesOpts{
			VolumeIds: []string{c.String("volume")},
		}
		result, err := lifecyclepolicy.AddVolumes(client, lifecyclePolicyID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var removeVolumeSubCommand = cli.Command{
	Name:      "remove",
	Usage:     "remove volume from lifecycle policy",
	ArgsUsage: argsUsagePolicyID,
	Category:  category,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "volume",
			Aliases:  []string{"v"},
			Usage:    "Volume ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		lifecyclePolicyID, err := flags.GetFirstIntArg(c, idErrorText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "remove")
			return err
		}
		client, err := client.NewLifecyclePolicyClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := lifecyclepolicy.RemoveVolumesOpts{
			VolumeIds: []string{c.String("volume")},
		}
		result, err := lifecyclepolicy.RemoveVolumes(client, lifecyclePolicyID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

func cronScheduleParamUsage(min, max int) string {
	return fmt.Sprintf("For '%s' type only. Comma-separated list of integers (%v-%v).", lifecyclepolicy.PolicyTypeCron, min, max)
}

func splitArg(s string) (int, rune, error) {
	suffix, _ := utf8.DecodeLastRuneInString(s)
	n, err := strconv.Atoi(strings.TrimSuffix(s, string(suffix)))
	if err == nil && n < 1 {
		err = fmt.Errorf("only positive integers accepted")
	}
	return n, suffix, err
}

func extractRetentionTime(c *cli.Context) (*lifecyclepolicy.RetentionTimer, error) {
	s := c.String("retention_time")
	if s == "" {
		return nil, nil
	}
	t := &lifecyclepolicy.RetentionTimer{}
	n, suffix, err := splitArg(s)
	if err != nil {
		return nil, fmt.Errorf("invalid value of retention_time flag (%s)", s)
	}
	switch suffix {
	case 'W':
		t.Weeks = n
	case 'D':
		t.Days = n
	case 'H':
		t.Hours = n
	case 'M':
		t.Minutes = n
	default:
		return nil, fmt.Errorf("invalid value of retention_time flag (%s)", s)
	}
	return t, nil
}

func extractCreateIntervalScheduleOpts(c *cli.Context) (*lifecyclepolicy.CreateIntervalScheduleOpts, error) {
	s := c.String("interval")
	if s == "" {
		return nil, fmt.Errorf("interval flag should be set for 'interval' type schedule")
	}
	opts := &lifecyclepolicy.CreateIntervalScheduleOpts{}
	n, suffix, err := splitArg(s)
	if err != nil {
		return nil, fmt.Errorf("invalid value of interval flag (%s)", s)
	}
	switch suffix {
	case 'W':
		opts.Weeks = n
	case 'D':
		opts.Days = n
	case 'H':
		opts.Hours = n
	case 'M':
		opts.Minutes = n
	default:
		return nil, fmt.Errorf("invalid value of interval flag (%s)", s)
	}
	return opts, nil
}

func extractCreateCronScheduleOpts(c *cli.Context) (opts *lifecyclepolicy.CreateCronScheduleOpts, err error) {
	opts = &lifecyclepolicy.CreateCronScheduleOpts{
		Timezone:  c.String("cron_timezone"),
		Week:      c.String("cron_week"),
		DayOfWeek: c.String("cron_day_of_week"),
		Month:     c.String("cron_month"),
		Day:       c.String("cron_day"),
		Hour:      c.String("cron_hour"),
		Minute:    c.String("cron_minute"),
	}
	return
}

var addScheduleSubCommand = cli.Command{
	Name:      "add",
	Usage:     "add schedule to lifecycle policy",
	ArgsUsage: argsUsagePolicyID,
	Category:  category,
	Flags: []cli.Flag{
		&cli.GenericFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Value: &utils.EnumValue{
				Enum: policyTypes,
			},
			Usage:    fmt.Sprintf("Schedule type (%s)", strings.Join(policyTypes, ", ")),
			Required: true,
		},
		&cli.StringFlag{
			Name:     "resource_name_template",
			Aliases:  []string{"rnt"},
			Usage:    "Template used for naming snapshots. All occurrences of '{volume_id}' will be replaced with ID of the volume",
			Value:    "reserve snap of the volume {volume_id}",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "max_quantity",
			Aliases:  []string{"mq"},
			Usage:    "Number of stored resources. It should be less than 10000",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "retention_time",
			Aliases:  []string{"rt"},
			Usage:    "If it is set, new snapshot will be deleted after specified time. Should be in format nX, where n is a positive integer and X is one of W, D, H, M (weeks, days, hours, minutes respectively)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "interval",
			Aliases:  []string{"i"},
			Usage:    fmt.Sprintf("Required if type set to '%s'. Specifies, how often new snapshot should be taken. Should be in format nX, where n is a positive integer and X is one of W, D, H, M (weeks, days, hours, minutes respectively)", lifecyclepolicy.PolicyTypeInterval),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cron_minute",
			Usage:    cronScheduleParamUsage(0, 59),
			Value:    "*",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cron_hour",
			Usage:    cronScheduleParamUsage(0, 23),
			Value:    "*",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cron_day",
			Usage:    cronScheduleParamUsage(1, 31),
			Value:    "*",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cron_month",
			Usage:    cronScheduleParamUsage(1, 12),
			Value:    "*",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cron_day_of_week",
			Usage:    cronScheduleParamUsage(0, 6),
			Value:    "*",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cron_week",
			Usage:    cronScheduleParamUsage(1, 53),
			Value:    "*",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "cron_timezone",
			Usage:    fmt.Sprintf("For '%s' type only. A pytz timezone", lifecyclepolicy.PolicyTypeCron),
			Value:    "UTC",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		lifecyclePolicyID, err := flags.GetFirstIntArg(c, idErrorText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "add")
			return err
		}
		client, err := client.NewLifecyclePolicyClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		retentionTime, err := extractRetentionTime(c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		common := lifecyclepolicy.CommonCreateScheduleOpts{
			Type:                 lifecyclepolicy.ScheduleType(c.String("type")),
			ResourceNameTemplate: c.String("resource_name_template"),
			MaxQuantity:          c.Int("max_quantity"),
			RetentionTime:        retentionTime,
		}
		var opts lifecyclepolicy.CreateScheduleOpts
		switch common.Type {
		case lifecyclepolicy.PolicyTypeCron:
			opts, err = extractCreateCronScheduleOpts(c)
		case lifecyclepolicy.PolicyTypeInterval:
			opts, err = extractCreateIntervalScheduleOpts(c)
		}
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		opts.SetCommonCreateScheduleOpts(common)
		addOpts := lifecyclepolicy.AddSchedulesOpts{Schedules: []lifecyclepolicy.CreateScheduleOpts{opts}}
		result, err := lifecyclepolicy.AddSchedules(client, lifecyclePolicyID, addOpts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var removeScheduleSubCommand = cli.Command{
	Name:      "remove",
	Usage:     "remove schedule from lifecycle policy",
	ArgsUsage: argsUsagePolicyID,
	Category:  category,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "schedule",
			Aliases:  []string{"s"},
			Usage:    "Schedule ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		lifecyclePolicyID, err := flags.GetFirstIntArg(c, idErrorText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "remove")
			return err
		}
		client, err := client.NewLifecyclePolicyClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := lifecyclepolicy.RemoveSchedulesOpts{
			ScheduleIDs: []string{c.String("schedule")},
		}
		result, err := lifecyclepolicy.RemoveSchedules(client, lifecyclePolicyID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}
