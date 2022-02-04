package limits

import (
	"github.com/G-Core/gcorelabscloud-go/client/limits/v2/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"reflect"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v2/limits"
	"github.com/urfave/cli/v2"
)

var (
	limitIDText = "limit_id is mandatory argument"
)

func buildLimitFromFlags(c *cli.Context) (limits.Limit, error) {
	var err error
	limit := limits.NewLimit()

	if c.IsSet("global") {
		globalFlagData := c.Generic("global")
		err = limit.GlobalLimits.Update(globalFlagData)
		if err != nil {
			return limit, err
		}
	}

	if c.IsSet("regions") {
		regionalData := c.Generic("regions")
		regionalEl := reflect.ValueOf(regionalData).Elem()
		for _, key := range regionalEl.MapKeys() {
			regionData := regionalEl.MapIndex(key).Interface()
			newRegionalLimits := limits.RegionalLimits{}
			err = newRegionalLimits.Update(regionData)
			if err != nil {
				return limit, err
			}
			limit.RegionalLimits = append(limit.RegionalLimits, newRegionalLimits)
		}
	}
	return limit, nil
}

var limitListCommand = cli.Command{
	Name:     "list",
	Usage:    "List limit requests",
	Category: "limit",
	Action: func(c *cli.Context) error {
		client, err := client.NewLimitClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		limit, err := limits.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(limit, c.String("format"))
		return nil
	},
}

var limitGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get limit request",
	ArgsUsage: "<limit_id>",
	Category:  "limit",
	Action: func(c *cli.Context) error {
		limitID, err := flags.GetFirstIntArg(c, limitIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewLimitClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		limitRequest, err := limits.Get(client, limitID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(limitRequest, c.String("format"))
		return nil
	},
}

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

type GlobalLimitsFlag limits.GlobalLimits

func (g *GlobalLimitsFlag) Set(value string) error {
	return utils.UpdateStructFromString(g, value)
}

func (g *GlobalLimitsFlag) String() string {
	return utils.StructToString(g)
}

type RegionLimitsFlag map[int]limits.RegionalLimits

func (g RegionLimitsFlag) Set(value string) error {
	regionItems := strings.Split(value, ",")
	for _, regionParams := range regionItems {
		newRegion := limits.RegionalLimits{}
		err := utils.UpdateStructFromString(&newRegion, regionParams)
		if err != nil {
			return err
		}

		g[newRegion.RegionID] = newRegion
	}
	return nil
}

func (g RegionLimitsFlag) String() string {
	return utils.StructToString(&limits.RegionalLimits{})
}

var limitCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create limit request",
	Category: "limit",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "description",
			Usage:    "limit request description",
			Required: true,
		},
		&cli.GenericFlag{
			Name:  "global",
			Value: &GlobalLimitsFlag{},
			Usage: "Create global limits, example [--global=key=value;key2=value2]",
		},
		&cli.GenericFlag{
			Name:  "regions",
			Value: &RegionLimitsFlag{},
			Usage: "Create regional limits. Use & for separate few regions with key/value content.\n" +
				"Example [--regions=region_id=1;key=value&region_id=2;key=value]",
		},
	}),
	Action: func(c *cli.Context) error {

		requestedLimits, err := buildLimitFromFlags(c)

		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		opts := limits.CreateOpts{
			Description:     c.String("description"),
			RequestedQuotas: requestedLimits,
		}

		client, err := client.NewLimitClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := limits.Create(client, opts).Extract()
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var Commands = cli.Command{
	Name:  "limit",
	Usage: "GCloud limits API",
	Subcommands: []*cli.Command{
		&limitListCommand,
		&limitDeleteCommand,
		&limitCreateCommand,
		&limitGetCommand,
	},
}
