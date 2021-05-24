package servergroups

import (
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/servergroups/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/urfave/cli/v2"
)

var serverGroupIDText = "servergroup_id is mandatory argument"

var Commands = cli.Command{
	Name:  "servergroup",
	Usage: "GCloud server groups v1 API",
	Subcommands: []*cli.Command{
		&serverGroupListCommand,
		&serverGroupGetCommand,
		&serverGroupDeleteCommand,
		&serverGroupCreateCommand,
	},
}

var serverGroupListCommand = cli.Command{
	Name:  "list",
	Usage: "List server group",
	Action: func(c *cli.Context) error {
		client, err := client.NewServerGroupClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := servergroups.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var serverGroupGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get server group information",
	ArgsUsage: "<servergroup_id>",
	Action: func(c *cli.Context) error {
		serverGroupID, err := flags.GetFirstStringArg(c, serverGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewServerGroupClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		sg, err := servergroups.Get(client, serverGroupID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(sg, c.String("format"))
		return nil
	},
}

var serverGroupDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete server group by ID",
	ArgsUsage: "<servergroup_id>",
	Action: func(c *cli.Context) error {
		serverGroupID, err := flags.GetFirstStringArg(c, serverGroupIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewServerGroupClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = servergroups.Delete(client, serverGroupID).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var serverGroupCreateCommand = cli.Command{
	Name:  "create",
	Usage: "Create server group",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Server group name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "policy",
			Aliases:  []string{"p"},
			Usage:    "Server group policy. Available value is 'affinity', 'anti-affinity'",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewServerGroupClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := servergroups.CreateOpts{
			Name:   c.String("name"),
			Policy: servergroups.ServerGroupPolicy(c.String("policy")),
		}
		result, err := servergroups.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}
