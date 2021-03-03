package ports

import (
	"github.com/G-Core/gcorelabscloud-go/client/ports/v1/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/urfave/cli/v2"
)

var (
	portIDText = "port_id is mandatory argument"
)

var portSecurityEnableSubCommand = cli.Command{
	Name:      "enable",
	Usage:     "Enable port security for instance interface",
	ArgsUsage: "<port_id>",
	Category:  "portsecurity",
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := client.NewPortClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		iface, err := ports.EnablePortSecurity(client, portID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(iface, c.String("format"))
		return nil
	},
}

var portSecurityDisableSubCommand = cli.Command{
	Name:      "disable",
	Usage:     "Disable port security for instance interface",
	ArgsUsage: "<port_id>",
	Category:  "portsecurity",
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := client.NewPortClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		iface, err := ports.DisablePortSecurity(client, portID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(iface, c.String("format"))
		return nil
	},
}

var PortCommands = cli.Command{
	Name:  "port",
	Usage: "GCloud ports API",
	Subcommands: []*cli.Command{
		&portSecurityEnableSubCommand,
		&portSecurityDisableSubCommand,
	},
}
