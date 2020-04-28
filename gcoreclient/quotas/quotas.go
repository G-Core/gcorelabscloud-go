package quotas

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/G-Core/gcorelabscloud-go/gcore/quota/v1/quotas"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/flags"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/utils"
	"github.com/urfave/cli/v2"
)

var (
	clientIDText = "client_id is mandatory argument"
)

func commandFlags(required bool, value int) []cli.Flag {
	quota := quotas.Quota{}.ToRequestMap()
	var res []cli.Flag
	var fields []string
	for f := range quota {
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

func buildQuotaFromFlags(c *cli.Context) (quotas.Quota, error) {
	m := quotas.Quota{}.ToRequestMap()
	quota := quotas.NewQuota()
	slf := reflect.ValueOf(&quota)
	s := slf.Elem()
	for f := range m {
		name := strings.ReplaceAll(f, "_", "-")
		fieldName := sanitizeFieldName(strcase.ToCamel(name))
		sf := s.FieldByName(fieldName)
		if sf.IsValid() && sf.CanSet() && sf.Kind() == reflect.Int {
			sf.SetInt(int64(c.Int(name)))
		} else {
			return quota, fmt.Errorf("cannot set field %s", fieldName)
		}
	}
	return quota, nil
}

var quotaListCommand = cli.Command{
	Name:     "list",
	Usage:    "List own quotas",
	Category: "quota",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "client_quotas", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		quota, err := quotas.OwnQuota(client).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(quota, c.String("format"))
		return nil
	},
}

var quotaGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get quota",
	ArgsUsage: "<client_id>",
	Category:  "quota",
	Action: func(c *cli.Context) error {
		quotaID, err := flags.GetFirstIntArg(c, clientIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "client_quotas", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := quotas.Get(client, quotaID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(task, c.String("format"))
		return nil
	},
}

var quotaReplaceCommand = cli.Command{
	Name:      "replace",
	Usage:     "Replace quota",
	ArgsUsage: "<client_id>",
	Category:  "quota",
	Flags:     commandFlags(true, quotas.Sentinel),
	Action: func(c *cli.Context) error {
		clientID, err := flags.GetFirstIntArg(c, clientIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}

		quota, err := buildQuotaFromFlags(c)

		if err != nil {
			_ = cli.ShowCommandHelp(c, "replace")
			return cli.NewExitError(err, 1)
		}

		opts := quotas.ReplaceOpts{Quota: quota}

		client, err := utils.BuildClient(c, "client_quotas", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := quotas.Replace(client, clientID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var quotaCreateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update quota partially",
	Category:  "quota",
	ArgsUsage: "<client_id>",
	Flags:     commandFlags(false, quotas.Sentinel),
	Action: func(c *cli.Context) error {

		clientID, err := flags.GetFirstIntArg(c, clientIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}

		quota, err := buildQuotaFromFlags(c)

		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		opts := quotas.UpdateOpts{Quota: quota}

		client, err := utils.BuildClient(c, "client_quotas", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := quotas.Update(client, clientID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var QuotaCommands = cli.Command{
	Name:  "quota",
	Usage: "GCloud quotas API",
	Subcommands: []*cli.Command{
		&quotaListCommand,
		&quotaGetCommand,
		&quotaReplaceCommand,
		&quotaCreateCommand,
	},
}
