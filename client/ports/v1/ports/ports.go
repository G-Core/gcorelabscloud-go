package ports

import (
	"errors"
	"github.com/G-Core/gcorelabscloud-go/client/ports/v1/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"

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

var assignAllowedAddressPairsSubCommand = cli.Command{
	Name:      "assign",
	Usage:     "Assign allowed address pairs for instance port",
	ArgsUsage: "<port_id>",
	Category:  "portsecurity",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "ip-address",
			Usage: "IP address of the port specified in allowed_address_pairs",
		},
		&cli.StringSliceFlag{
			Name:  "mac-address",
			Usage: "MAC address of the port specified in allowed_address_pairs",
		},
	},
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "assign")
			return err
		}
		client, err := client.NewPortClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		addressPairs, err := getAddressPairs(c.StringSlice("ip-address"), c.StringSlice("mac-address"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "assign")
			return cli.NewExitError(err, 1)
		}

		opts := ports.AllowAddressPairsOpts{
			AllowedAddressPairs: addressPairs,
		}

		result, err := ports.AllowAddressPairs(client, portID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

func getAddressPairs(ips, macs []string) ([]reservedfixedips.AllowedAddressPairs, error) {
	if len(ips) != len(macs) {
		return nil, errors.New("length of ip-address and mac-address should be equal")
	}

	result := make([]reservedfixedips.AllowedAddressPairs, len(ips))
	for i, ipRaw := range ips {
		result[i] = reservedfixedips.AllowedAddressPairs{
			IPAddress:  ipRaw,
			MacAddress: macs[i],
		}
	}
	return result, nil
}

var Commands = cli.Command{
	Name:  "port",
	Usage: "GCloud ports API",
	Subcommands: []*cli.Command{
		&portSecurityEnableSubCommand,
		&portSecurityDisableSubCommand,
		&assignAllowedAddressPairsSubCommand,
	},
}
