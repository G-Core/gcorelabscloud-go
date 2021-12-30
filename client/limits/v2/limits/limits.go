package limits

import (
	"github.com/G-Core/gcorelabscloud-go/client/limits/v2/client"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v2/limits"
	"github.com/urfave/cli/v2"
)

var (
	limitIDText = "limit_id is mandatory argument"
)

var limitDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete limit request",
	ArgsUsage: "<limit_id>",
	Category:  "limit",
	Action: func(c *cli.Context) error {
		limitID, err := flags.GetFirstIntArg(c, limitIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewLimitClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = limits.Delete(client, limitID).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var Commands = cli.Command{
	Name:  "limit",
	Usage: "GCloud limits API",
	Subcommands: []*cli.Command{
		&limitDeleteCommand,
	},
}
