package l7rules

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/l7policies/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/l7policies"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

var ruleIDText = "rule_id is mandatory argument"

var L7RuleCommands = cli.Command{
	Name:  "l7rule",
	Usage: "GCloud l7rule API",
	Subcommands: []*cli.Command{
		&l7RuleListSubCommand,
		&l7RuleGetSubCommand,
		&l7RuleReplaceSubCommand,
		&l7RuleDeleteSubCommand,
		&l7RuleCreateSubCommand,
	},
}

var l7RuleListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "List l7rules",
	Category: "l7rule",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "policy-id",
			Aliases:  []string{"p"},
			Usage:    "Policy id",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		policyID := c.String("policy-id")
		client, err := client.NewL7RulesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := l7policies.ListAllRule(client, policyID)
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

var l7RuleGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show l7rules",
	ArgsUsage: "<l7rule_id>",
	Category:  "l7rule",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "policy-id",
			Aliases:  []string{"p"},
			Usage:    "Policy id",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		ruleID, err := flags.GetFirstStringArg(c, ruleIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}

		policyID := c.String("policy-id")
		client, err := client.NewL7RulesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := l7policies.GetRule(client, policyID, ruleID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}

		return nil
	},
}

var l7RuleDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete l7rules",
	ArgsUsage: "<l7rule_id>",
	Category:  "l7rule",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "policy-id",
			Aliases:  []string{"p"},
			Usage:    "Policy id",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		ruleID, err := flags.GetFirstStringArg(c, ruleIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}

		policyID := c.String("policy-id")
		client, err := client.NewL7RulesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := l7policies.DeleteRule(client, policyID, ruleID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}

		return nil
	},
}

var l7RuleCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create l7rules",
	Category: "l7rule",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "policy-id",
			Aliases:  []string{"p"},
			Usage:    "Policy id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "compare-type",
			Usage:    `The comparison type for the L7 rule. Available value is "CONTAINS" "ENDS_WITH" "EQUAL_TO" "REGEX" "STARTS_WITH"`,
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "invert",
			Usage: `When true the logic of the rule is inverted. For example, with invert true, 'equal to' would become 'not equal to'. Default is false.`,
		},
		&cli.StringFlag{
			Name:  "key",
			Usage: "The key to use for the comparison. For example, the name of the cookie to evaluate.",
		},
		&cli.StringFlag{
			Name:     "type",
			Usage:    `The L7 rule type. Available value is "COOKIE" "FILE_TYPE" "HEADER" "HOST_NAME" "PATH" "SSL_CONN_HAS_CERT" "SSL_VERIFY_RESULT" "SSL_DN_FIELD"`,
			Required: true,
		},
		&cli.StringFlag{
			Name:     "value",
			Usage:    "The value to use for the comparison. For example, the file type to compare",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:  "tags",
			Usage: "A list of simple strings assigned to the resource.",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		policyID := c.String("policy-id")
		client, err := client.NewL7RulesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := l7policies.CreateRuleOpts{
			CompareType: l7policies.CompareType(c.String("compare-type")),
			Invert:      c.Bool("invert"),
			Key:         c.String("key"),
			Type:        l7policies.RuleType(c.String("type")),
			Value:       c.String("value"),
			Tags:        c.StringSlice("tags"),
		}
		results, err := l7policies.CreateRule(client, policyID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			ruleID, err := l7policies.ExtractRuleIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve rule ID from task info: %w", err)
			}
			rule, err := l7policies.GetRule(client, policyID, ruleID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get rule with ID: %s. Error: %w", ruleID, err)
			}
			utils.ShowResults(rule, c.String("format"))
			return nil, nil
		})
	},
}

var l7RuleReplaceSubCommand = cli.Command{
	Name:      "replace",
	Usage:     "Replace l7rules",
	ArgsUsage: "<l7rule_id>",
	Category:  "l7rule",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "policy-id",
			Aliases:  []string{"p"},
			Usage:    "Policy id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "compare-type",
			Usage:    `The comparison type for the L7 rule. Available value is "CONTAINS" "ENDS_WITH" "EQUAL_TO" "REGEX" "STARTS_WITH"`,
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "invert",
			Usage: `When true the logic of the rule is inverted. For example, with invert true, 'equal to' would become 'not equal to'. Default is false.`,
		},
		&cli.StringFlag{
			Name:  "key",
			Usage: "The key to use for the comparison. For example, the name of the cookie to evaluate.",
		},
		&cli.StringFlag{
			Name:     "type",
			Usage:    `The L7 rule type. Available value is "COOKIE" "FILE_TYPE" "HEADER" "HOST_NAME" "PATH" "SSL_CONN_HAS_CERT" "SSL_VERIFY_RESULT" "SSL_DN_FIELD"`,
			Required: true,
		},
		&cli.StringFlag{
			Name:     "value",
			Usage:    "The value to use for the comparison. For example, the file type to compare",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:  "tags",
			Usage: "A list of simple strings assigned to the resource.",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		ruleID, err := flags.GetFirstStringArg(c, ruleIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}

		policyID := c.String("policy-id")
		client, err := client.NewL7RulesClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := l7policies.CreateRuleOpts{
			CompareType: l7policies.CompareType(c.String("compare-type")),
			Invert:      c.Bool("invert"),
			Key:         c.String("key"),
			Type:        l7policies.RuleType(c.String("type")),
			Value:       c.String("value"),
			Tags:        c.StringSlice("tags"),
		}
		results, err := l7policies.ReplaceRule(client, policyID, ruleID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			ruleID, err := l7policies.ExtractRuleIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve rule ID from task info: %w", err)
			}
			rule, err := l7policies.GetRule(client, policyID, ruleID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get rule with ID: %s. Error: %w", ruleID, err)
			}
			utils.ShowResults(rule, c.String("format"))
			return nil, nil
		})
	},
}
