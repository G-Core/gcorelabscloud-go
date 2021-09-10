package apptemplates

import (
	"github.com/G-Core/gcorelabscloud-go/client/apptemplates/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/apptemplate/v1/apptemplates"
	"github.com/urfave/cli/v2"
)

var Commands = cli.Command{
	Name:  "apptemplates",
	Usage: "GCloud apptemplates API",
	Subcommands: []*cli.Command{
		&appTemplateListSubCommand,
		&appTemplateGetSubCommand,
	},
}

var appTemplateIDText = "apptemplate_id is mandatory argument"

var appTemplateListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "List apptemplates",
	Category: "apptemplates",
	Action: func(c *cli.Context) error {
		client, err := client.NewAppTemplateClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := apptemplates.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if results != nil {
			utils.ShowResults(results, c.String("format"))
		}
		return nil
	},
}

var appTemplateGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show apptemplate",
	ArgsUsage: "<apptemplate_id>",
	Category:  "router",
	Action: func(c *cli.Context) error {
		appTemplateID, err := flags.GetFirstStringArg(c, appTemplateIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewAppTemplateClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := apptemplates.Get(client, appTemplateID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}
