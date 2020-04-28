package templates

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/gcore/magnum/v1/clustertemplates"
	"github.com/G-Core/gcorelabscloud-go/gcore/magnum/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/flags"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var (
	clusterTemplateIDText = "clustertemplate_id is mandatory argument"
	clusterUpdateTypes    = types.ClusterUpdateOperation("").StringList()
)

var clusterTemplateCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Magnum create cluster template",
	Category: "template",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "image",
			Aliases:  []string{"i"},
			Usage:    "Base image in Glance",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "keypair",
			Aliases:  []string{"k"},
			Usage:    "The name of the SSH keypair",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Cluster template name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "docker-volume-size",
			Usage:    "The size in GB for the local storage on each server for the Docker daemon to cache the images and host the containers",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "external-network-id",
			Usage:    "External network ID for cluster",
			Required: false,
		},
		&cli.StringFlag{
			Name:        "fixed-subnet",
			Usage:       "Fixed subnet that are using to allocate network address for nodes in cluster.",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "master-flavor",
			Usage:       "The flavor of the master node for this cluster template",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:     "flavor",
			Usage:    "The flavor of the node for this cluster template",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:        "label",
			Usage:       "Arbitrary labels. The accepted keys and valid values are defined in the cluster drivers. --label one=two --label three=four ",
			DefaultText: "nil",
			Required:    false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		labels, err := utils.StringSliceToMap(c.StringSlice("label"))
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := clustertemplates.CreateOpts{
			ImageID:           c.String("image"),
			KeyPairID:         c.String("keypair"),
			Name:              c.String("name"),
			DockerVolumeSize:  c.Int("docker-volume-size"),
			Labels:            &labels,
			ExternalNetworkID: c.String("external-network-id"),
			FixedSubnet:       c.String("fixed-subnet"),
			MasterFlavorID:    c.String("master-flavor"),
			FlavorID:          c.String("flavor"),
		}
		result, err := clustertemplates.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var clusterTemplateListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Magnum list cluster templates",
	Category: "template",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clustertemplates.ListAll(client, clustertemplates.ListOpts{})
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var clusterTemplateDeleteDubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Magnum delete cluster template",
	ArgsUsage: "<template_id>",
	Category:  "template",
	Action: func(c *cli.Context) error {
		clusterTemplateID, err := flags.GetFirstStringArg(c, clusterTemplateIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = clustertemplates.Delete(client, clusterTemplateID).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var clusterTemplateGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Magnum get cluster template",
	ArgsUsage: "<template_id>",
	Category:  "template",
	Action: func(c *cli.Context) error {
		clusterTemplateID, err := flags.GetFirstStringArg(c, clusterTemplateIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := clustertemplates.Get(client, clusterTemplateID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var clusterTemplateUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Magnum update cluster template",
	ArgsUsage: "<template_id>",
	Category:  "template",
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:     "path",
			Aliases:  []string{"p"},
			Usage:    "Update json path. Example /node/count",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "value",
			Aliases:  []string{"v"},
			Usage:    "Update json value. For path /labels in format: label_one=value_one,label_two=value_two",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "op",
			Aliases: []string{"o"},
			Value: &utils.EnumStringSliceValue{
				Enum: clusterUpdateTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(clusterUpdateTypes, ", ")),
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterTemplateID, err := flags.GetFirstStringArg(c, clusterTemplateIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		paths := c.StringSlice("path")
		values := c.StringSlice("value")
		ops := utils.GetEnumStringSliceValue(c, "op")

		if len(paths) != len(values) || len(values) != len(ops) {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(fmt.Errorf("path, value and op parameters number should be same"), 1)
		}

		var opts clustertemplates.UpdateOpts

		for idx, path := range paths {
			if !strings.HasPrefix(path, "/") {
				return cli.NewExitError(fmt.Errorf("path parameter should be in format /path"), 1)
			}
			var updateValue interface{}
			value := values[idx]
			intValue, err := strconv.Atoi(value)
			if err == nil {
				updateValue = intValue
			} else if path == "/labels" {
				updateValue, err = utils.StringSliceToMap(strings.Split(value, ","))
				if err != nil {
					return cli.NewExitError(fmt.Errorf("wrong labels format. should be in format: label_one=value_one,label_two=value_two"), 1)
				}
			} else {
				updateValue = value
			}
			el := clustertemplates.UpdateOptsElem{
				Path:  path,
				Value: updateValue,
				Op:    types.ClusterUpdateOperation(ops[idx]),
			}
			opts = append(opts, el)
		}

		results, err := clustertemplates.Update(client, clusterTemplateID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var ClusterTemplatesCommands = cli.Command{
	Name:  "template",
	Usage: "Magnum cluster template commands",
	Subcommands: []*cli.Command{
		&clusterTemplateCreateSubCommand,
		&clusterTemplateListSubCommand,
		&clusterTemplateDeleteDubCommand,
		&clusterTemplateGetSubCommand,
		&clusterTemplateUpdateSubCommand,
	},
}
