package ports

import (
	"errors"
	"fmt"

	"github.com/G-Core/gcorelabscloud-go/client/ports/v1/client"
	client2 "github.com/G-Core/gcorelabscloud-go/client/ports/v2/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	ports1 "github.com/G-Core/gcorelabscloud-go/gcore/port/v1/ports"
	ports2 "github.com/G-Core/gcorelabscloud-go/gcore/port/v2/ports"
	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

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
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:  "ip-address",
			Usage: "IP address of the port specified in allowed_address_pairs",
		},
		&cli.StringSliceFlag{
			Name:  "mac-address",
			Usage: "MAC address of the port specified in allowed_address_pairs",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		portID, err := flags.GetFirstStringArg(c, portIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "assign")
			return err
		}
		clientV1, err := client.NewPortClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		clientV2, err := client2.NewPortClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		addressPairs, err := getAddressPairs(c.StringSlice("ip-address"), c.StringSlice("mac-address"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "assign")
			return cli.NewExitError(err, 1)
		}

		opts := ports1.AllowAddressPairsOpts{
			AllowedAddressPairs: addressPairs,
		}

		results, err := ports2.AllowAddressPairs(clientV2, portID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, clientV1, results, true, func(task tasks.TaskID) (interface{}, error) {
			_, err := tasks.Get(clientV1, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			return nil, nil
		})
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
