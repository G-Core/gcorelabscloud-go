package limits

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/client/limits/v1/client"

	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v1/types"

	"github.com/iancoleman/strcase"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/limit/v1/limits"
	"github.com/urfave/cli/v2"
)

var (
	limitIDText              = "limit_id is mandatory argument"
	limitRequestStatusesList = types.LimitRequestStatus("").StringList()
)

func commandFlags(required bool, value int) []cli.Flag {
	limit := limits.Limit{}.ToRequestMap()
	var res []cli.Flag
	var fields []string
	for f := range limit {
		fields = append(fields, f)
	}
	sort.Strings(fields)
	for _, f := range fields {
		nameParts := strings.Split(f, "_")
		usage := strings.Join(nameParts, " ")
		name := strings.Join(nameParts, "-")
		res = append(res, &cli.IntFlag{
			Name:     name,
			Usage:    usage,
			Required: required,
			Value:    value,
		})
	}
	return res
}

func sanitizeFieldName(n string) string {
	n = strings.ReplaceAll(n, "Ip", "IP")
	n = strings.ReplaceAll(n, "Cpu", "CPU")
	n = strings.ReplaceAll(n, "Vm", "VM")
	n = strings.ReplaceAll(n, "Ram", "RAM")
	return n
}

func buildLimitFromFlags(c *cli.Context) (limits.Limit, error) {
	m := limits.Limit{}.ToRequestMap()
	limit := limits.NewLimit()
	slf := reflect.ValueOf(&limit)
	s := slf.Elem()
	for f := range m {
		name := strings.ReplaceAll(f, "_", "-")
		fieldName := sanitizeFieldName(strcase.ToCamel(name))
		sf := s.FieldByName(fieldName)
		if sf.IsValid() && sf.CanSet() && sf.Kind() == reflect.Int {
			sf.SetInt(int64(c.Int(name)))
		} else {
			return limit, fmt.Errorf("cannot set field %s", fieldName)
		}
	}
	return limit, nil
}

var limitListCommand = cli.Command{
	Name:     "list",
	Usage:    "List limit requests",
	Category: "limit",
	Action: func(c *cli.Context) error {
		client, err := client.NewLimitClientV1(c)
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
		client, err := client.NewLimitClientV1(c)
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
		client, err := client.NewLimitClientV1(c)
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
	}, commandFlags(false, limits.Sentinel)...),
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

		client, err := client.NewLimitClientV1(c)
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

var limitStatusCommand = cli.Command{
	Name:      "status",
	Usage:     "Change limit request status",
	Category:  "limit",
	ArgsUsage: "<limit_id>",
	Flags: []cli.Flag{
		&cli.GenericFlag{
			Name:    "status",
			Aliases: []string{"s"},
			Value: &utils.EnumValue{
				Enum: limitRequestStatusesList,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(limitRequestStatusesList, ", ")),
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		limitID, err := flags.GetFirstIntArg(c, limitIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "status")
			return err
		}

		opts := limits.StatusOpts{
			Status: types.LimitRequestStatus(c.String("status")),
		}

		client, err := client.NewLimitClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := limits.Status(client, limitID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var limitUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update limit partially",
	Category:  "limit",
	ArgsUsage: "<limit_id>",
	Flags:     commandFlags(false, limits.Sentinel),
	Action: func(c *cli.Context) error {
		limitID, err := flags.GetFirstIntArg(c, limitIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}

		limit, err := buildLimitFromFlags(c)

		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		opts := limits.UpdateOpts{Limit: limit}

		client, err := client.NewLimitClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := limits.Update(client, limitID, opts).Extract()
		if err != nil {
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
		&limitGetCommand,
		&limitCreateCommand,
		&limitUpdateCommand,
		&limitDeleteCommand,
		&limitStatusCommand,
	},
}
