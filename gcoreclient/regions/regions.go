package regions

import (
	"fmt"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/regions"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/flags"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/utils"
	"github.com/urfave/cli/v2"
)

var (
	regionIDText           = "region_id is mandatory argument"
	regionStatesList       = types.RegionState("").StringList()
	regionEndpointTypeList = types.EndpointType("").StringList()
)

var regionListCommand = cli.Command{
	Name:     "list",
	Usage:    "List regions",
	Category: "region",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "regions", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := regions.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var regionGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get region",
	ArgsUsage: "<region_id>",
	Category:  "region",
	Action: func(c *cli.Context) error {
		regionID, err := flags.GetFirstIntArg(c, regionIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "regions", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := regions.Get(client, regionID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(task, c.String("format"))
		return nil
	},
}

var regionUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update region",
	ArgsUsage: "<region_id>",
	Category:  "region",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "display-name",
			Aliases:  []string{"n"},
			Usage:    "region name",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "state",
			Aliases: []string{"s"},
			Value: &utils.EnumValue{
				Enum: regionStatesList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(regionStatesList, ", ")),
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "endpoint-type",
			Aliases: []string{"e"},
			Value: &utils.EnumValue{
				Enum: regionEndpointTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(regionEndpointTypeList, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "network-id",
			Usage:    "external network id",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "spice-url",
			Usage:    "spice proxy url",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		regionID, err := flags.GetFirstIntArg(c, regionIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}

		url, err := gcorecloud.ParseURLNonMandatory(c.String("spice-url"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		opts := regions.UpdateOpts{
			DisplayName:       c.String("display-name"),
			State:             types.RegionState(c.String("state")),
			EndpointType:      types.EndpointType(c.String("endpoint-type")),
			ExternalNetworkID: c.String("network-id"),
			SpiceProxyURL:     url,
		}

		err = gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		client, err := utils.BuildClient(c, "regions", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := regions.Update(client, regionID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var regionCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create region",
	Category: "region",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "display-name",
			Aliases:  []string{"n"},
			Usage:    "region name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "keystone-name",
			Aliases:  []string{"k"},
			Usage:    "keystone name",
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "state",
			Aliases: []string{"s"},
			Value: &utils.EnumValue{
				Enum: regionStatesList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(regionStatesList, ", ")),
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "endpoint-type",
			Aliases: []string{"e"},
			Value: &utils.EnumValue{
				Enum: regionEndpointTypeList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(regionEndpointTypeList, ", ")),
			Required: true,
		},
		&cli.StringFlag{
			Name:     "network-id",
			Usage:    "external network id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "spice-url",
			Usage:    "spice proxy url",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "keystone-id",
			Usage:    "keystone id",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {

		url, err := gcorecloud.ParseURLNonMandatory(c.String("spice-url"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		opts := regions.CreateOpts{
			DisplayName:       c.String("display-name"),
			KeystoneName:      c.String("keystone-name"),
			State:             types.RegionState(c.String("state")),
			EndpointType:      types.EndpointType(c.String("endpoint-type")),
			ExternalNetworkID: c.String("network-id"),
			SpiceProxyURL:     url,
			KeystoneID:        c.Int("keystone-id"),
		}

		err = gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		client, err := utils.BuildClient(c, "regions", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := regions.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var RegionCommands = cli.Command{
	Name:  "region",
	Usage: "GCloud regions API",
	Subcommands: []*cli.Command{
		&regionListCommand,
		&regionGetCommand,
		&regionUpdateCommand,
		&regionCreateCommand,
	},
}
