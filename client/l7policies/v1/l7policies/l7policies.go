package l7policies

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/l7policies/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/l7policies/v1/l7rules"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/l7policies"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

var policyIDText = "policy_id is mandatory argument"

var L7PolicyCommands = cli.Command{
	Name:  "l7policy",
	Usage: "GCloud l7policy API",
	Subcommands: []*cli.Command{
		&l7policyListSubCommand,
		&l7policyGetSubCommand,
		&l7policyReplaceSubCommand,
		&l7policyDeleteSubCommand,
		&l7policyCreateSubCommand,
		&l7rules.L7RuleCommands,
	},
}

var l7policyListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "List l7policies",
	Category: "l7policy",
	Action: func(c *cli.Context) error {
		client, err := client.NewL7PoliciesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := l7policies.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var l7policyGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show l7policies",
	ArgsUsage: "<l7policy_id>",
	Category:  "l7policy",
	Action: func(c *cli.Context) error {
		policyID, err := flags.GetFirstStringArg(c, policyIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewL7PoliciesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := l7policies.Get(client, policyID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var l7policyDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete l7policies",
	ArgsUsage: "<l7policy_id>",
	Category:  "l7policy",
	Action: func(c *cli.Context) error {
		policyID, err := flags.GetFirstStringArg(c, policyIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}

		client, err := client.NewL7PoliciesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := l7policies.Delete(client, policyID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var l7policyCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create l7policies",
	Category: "l7policy",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "L7policy name",
		},
		&cli.StringFlag{
			Name:     "listener-id",
			Usage:    "L7policy listener id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "action",
			Usage:    `L7policy action. Available value is "REDIRECT_PREFIX" "REDIRECT_TO_POOL" "REDIRECT_TO_URL" "REJECT"`,
			Required: true,
		},
		&cli.IntFlag{
			Name:  "position",
			Usage: "The position of this policy on the listener. Positions start at 1.",
		},
		&cli.IntFlag{
			Name: "redirect-http-code",
			Usage: `Requests matching this policy will be redirected to the specified URL or Prefix URL with the HTTP response code.
			Valid if action is REDIRECT_TO_URL or REDIRECT_PREFIX. Valid options are 301, 302, 303, 307, or 308. Default is 302.`,
		},
		&cli.StringFlag{
			Name:  "redirect-pool-id",
			Usage: `Requests matching this policy will be redirected to the pool with this ID. Only valid if action is REDIRECT_TO_POOL.`,
		},
		&cli.StringFlag{
			Name:  "redirect-prefix",
			Usage: "Requests matching this policy will be redirected to this Prefix URL. Only valid if action is REDIRECT_PREFIX.",
		},
		&cli.StringFlag{
			Name:  "redirect-url",
			Usage: "Requests matching this policy will be redirected to this URL. Only valid if action is REDIRECT_TO_URL.",
		},
		&cli.StringSliceFlag{
			Name:  "tags",
			Usage: "A list of simple strings assigned to the resource.",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := client.NewL7PoliciesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := l7policies.CreateOpts{
			Name:             c.String("name"),
			ListenerID:       c.String("listener-id"),
			Action:           l7policies.Action(c.String("action")),
			Position:         int32(c.Int("position")),
			RedirectHTTPCode: c.Int("redirect-http-code"),
			RedirectPoolID:   c.String("redirect-pool-id"),
			RedirectPrefix:   c.String("redirect-prefix"),
			RedirectURL:      c.String("redirect-url"),
			Tags:             c.StringSlice("tags"),
		}
		results, err := l7policies.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			policyID, err := l7policies.ExtractL7PolicyIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve policy ID from task info: %w", err)
			}
			router, err := l7policies.Get(client, policyID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get policy with ID: %s. Error: %w", policyID, err)
			}
			utils.ShowResults(router, c.String("format"))
			return nil, nil
		})
	},
}

var l7policyReplaceSubCommand = cli.Command{
	Name:      "replace",
	Usage:     "Replace l7policies",
	Category:  "l7policy",
	ArgsUsage: "<l7policy_id>",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "L7policy name",
		},
		&cli.StringFlag{
			Name:     "action",
			Usage:    `L7policy action. Available value is "REDIRECT_PREFIX" "REDIRECT_TO_POOL" "REDIRECT_TO_URL" "REJECT"`,
			Required: true,
		},
		&cli.IntFlag{
			Name:  "position",
			Usage: "The position of this policy on the listener. Positions start at 1.",
		},
		&cli.IntFlag{
			Name: "redirect-http-code",
			Usage: `Requests matching this policy will be redirected to the specified URL or Prefix URL with the HTTP response code.
			Valid if action is REDIRECT_TO_URL or REDIRECT_PREFIX. Valid options are 301, 302, 303, 307, or 308. Default is 302.`,
		},
		&cli.StringFlag{
			Name:  "redirect-pool-id",
			Usage: `Requests matching this policy will be redirected to the pool with this ID. Only valid if action is REDIRECT_TO_POOL.`,
		},
		&cli.StringFlag{
			Name:  "redirect-prefix",
			Usage: "Requests matching this policy will be redirected to this Prefix URL. Only valid if action is REDIRECT_PREFIX.",
		},
		&cli.StringFlag{
			Name:  "redirect-url",
			Usage: "Requests matching this policy will be redirected to this URL. Only valid if action is REDIRECT_TO_URL.",
		},
		&cli.StringSliceFlag{
			Name:  "tags",
			Usage: "A list of simple strings assigned to the resource.",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		policyID, err := flags.GetFirstStringArg(c, policyIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}

		client, err := client.NewL7PoliciesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := l7policies.ReplaceOpts{
			Name:             c.String("name"),
			Action:           l7policies.Action(c.String("action")),
			Position:         int32(c.Int("position")),
			RedirectHTTPCode: c.Int("redirect-http-code"),
			RedirectPoolID:   c.String("redirect-pool-id"),
			RedirectPrefix:   c.String("redirect-prefix"),
			RedirectURL:      c.String("redirect-url"),
			Tags:             c.StringSlice("tags"),
		}
		results, err := l7policies.Replace(client, policyID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			policyID, err := l7policies.ExtractL7PolicyIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve policy ID from task info: %w", err)
			}
			router, err := l7policies.Get(client, policyID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get policy with ID: %s. Error: %w", policyID, err)
			}
			utils.ShowResults(router, c.String("format"))
			return nil, nil
		})
	},
}
