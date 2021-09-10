package schedules

import (
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/schedules/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"
	"github.com/G-Core/gcorelabscloud-go/gcore/schedule/v1/schedules"
	"github.com/urfave/cli/v2"
)

var scheduleIDText = "schedule_id is mandatory argument"

var Commands = cli.Command{
	Name:  "schedule",
	Usage: "GCloud schedule API",
	Subcommands: []*cli.Command{
		&scheduleGetSubCommand,
		&scheduleUpdateSubCommand,
	},
}

var scheduleGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show schedule",
	ArgsUsage: "<schedule_id>",
	Category:  "schedule",
	Action: func(c *cli.Context) error {
		scheduleID, err := flags.GetFirstStringArg(c, scheduleIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewScheduleClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := schedules.Get(client, scheduleID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			schedule, err := result.Cook()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			utils.ShowResults(schedule, c.String("format"))
		}
		return nil
	},
}

var scheduleUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Update schedule",
	ArgsUsage: "<schedule_id>",
	Category:  "schedule",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "hour",
			Aliases: []string{"h"},
			Usage:   "Hour (0-23, '*') or a comma-separated string listing of hours.",
		},
		&cli.StringFlag{
			Name:  "month",
			Usage: "Month (1-12, '*') or a comma-separated string listing of months.",
		},
		&cli.StringFlag{
			Name:    "timezone",
			Aliases: []string{"tz"},
			Usage:   "A pytz timezone. It's available if type is cron. As a default, it uses UTC.",
		},
		&cli.IntFlag{
			Name:    "weeks",
			Aliases: []string{"w"},
		},
		&cli.StringFlag{
			Name:  "minute",
			Usage: "Minute (0-59, '*') or a comma-separated string listing of minutes.",
		},
		&cli.StringFlag{
			Name:    "day",
			Aliases: []string{"d"},
			Usage:   "Day of the month (1-31, '*') or a comma-separated string listing of days.",
		},
		&cli.IntFlag{
			Name:  "retention-time-days",
			Usage: "Number of days to wait",
		},
		&cli.IntFlag{
			Name:  "retention-time-weeks",
			Usage: "Number of weeks to wait",
		},
		&cli.IntFlag{
			Name:  "retention-time-hours",
			Usage: "Number of hours to wait",
		},
		&cli.IntFlag{
			Name:  "retention-time-minutes",
			Usage: "Number of minutes to wait",
		},
		&cli.StringFlag{
			Name:    "time",
			Aliases: []string{"t"},
			Usage:   "Schedule time",
		},
		&cli.StringFlag{
			Name:  "resource-name-template",
			Usage: "Resource name templates.",
		},
		&cli.IntFlag{
			Name:  "days",
			Usage: "Number of days to wait",
		},
		&cli.IntFlag{
			Name:  "minutes",
			Usage: "Number of minutes to wait",
		},
		&cli.IntFlag{
			Name:  "hours",
			Usage: "Number of hours to wait",
		},
		&cli.StringFlag{
			Name:  "week",
			Usage: "ISO week (1-53, '*') or a comma-separated string listing of weeks.",
		},
		&cli.IntFlag{
			Name:     "max-quantity",
			Usage:    "Number of stored resources. It should be less than 10000.",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "day-of-week",
			Usage: "Week day or a comma-separated string listing of week days. enum[mon,tue,wed,thu,fri,sat,sun,*]",
		},
		&cli.StringFlag{
			Name:  "type",
			Usage: "Schedule time type.",
		},
	},
	Action: func(c *cli.Context) error {
		scheduleID, err := flags.GetFirstStringArg(c, scheduleIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewScheduleClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := schedules.UpdateOpts{
			Hour:     c.String("hour"),
			Month:    c.String("month"),
			Timezone: c.String("timezone"),
			Weeks:    c.Int("weeks"),
			Minute:   c.String("minute"),
			Day:      c.String("day"),
			RetentionTime: lifecyclepolicy.RetentionTimer{
				Weeks:   c.Int("retention-time-weeks"),
				Days:    c.Int("retention-time-days"),
				Hours:   c.Int("retention-time-hours"),
				Minutes: c.Int("retention-time-minutes"),
			},
			Time:                 c.String("time"),
			ResourceNameTemplate: c.String("resource-name-template"),
			Days:                 c.Int("days"),
			Minutes:              c.Int("minutes"),
			Week:                 c.String("week"),
			MaxQuantity:          c.Int("max-quantity"),
			DayOfWeek:            c.String("day-of-week"),
			Hours:                c.Int("hours"),
			Type:                 lifecyclepolicy.ScheduleType(c.String("type")),
		}

		result, err := schedules.Update(client, scheduleID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			schedule, err := result.Cook()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			utils.ShowResults(schedule, c.String("format"))
		}
		return nil
	},
}
