package k8s

import (
	"fmt"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/k8s/v2/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/urfave/cli/v2"
)

var (
	poolNameText = "pool_name is mandatory argument"
)

var poolListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "List cluster pools",
	Category: "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-name",
			Aliases:  []string{"c"},
			Usage:    "Cluster name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		clusterName := c.String("cluster-name")
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := pools.ListAll(client, clusterName)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var poolGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Get cluster pool information",
	ArgsUsage: "<pool_name>",
	Category:  "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-name",
			Aliases:  []string{"c"},
			Usage:    "Cluster name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		poolName, err := flags.GetFirstStringArg(c, poolNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return cli.NewExitError(err, 1)
		}
		clusterName := c.String("cluster-name")
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := pools.Get(client, clusterName, poolName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var poolCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create cluster pool",
	Category: "pool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-name",
			Aliases:  []string{"c"},
			Usage:    "Cluster name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Pool name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "flavor-id",
			Usage:    "Pool node flavor",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "min-node-count",
			Usage:    "Minimum number of pool nodes",
			Value:    1,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "max-node-count",
			Usage:    "Maximum number of pool nodes",
			Value:    1,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "boot-volume-size",
			Usage:    "Boot volume size in GB",
			Value:    10,
			Required: false,
		},
		&cli.GenericFlag{
			Name: "boot-volume-type",
			Value: &utils.EnumValue{
				Enum:    volumeTypeNames,
				Default: volumeTypeNames[0],
			},
			Usage:    fmt.Sprintf("Boot volume type [%s]", strings.Join(volumeTypeNames, ", ")),
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "auto-healing-enabled",
			Usage:    "Enable/disable auto healing",
			Value:    false,
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "is-public-ipv4",
			Usage:    "Enable public IPv4 address",
			Value:    false,
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterName := c.String("cluster-name")
		poolName := c.String("name")
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := pools.CreateOpts{
			Name:               poolName,
			FlavorID:           c.String("flavor-id"),
			MinNodeCount:       c.Int("min-node-count"),
			MaxNodeCount:       c.Int("max-node-count"),
			BootVolumeSize:     c.Int("boot-volume-size"),
			BootVolumeType:     volumes.VolumeType(c.String("boot-volume-type")),
			AutoHealingEnabled: c.Bool("auto-healing-enabled"),
			IsPublicIPv4:       c.Bool("is-public-ipv4"),
		}
		results, err := pools.Create(client, clusterName, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		tc := getTaskClient(c, client)
		return utils.WaitTaskAndShowResult(c, tc, results, true, func(task tasks.TaskID) (interface{}, error) {
			pool, err := pools.Get(client, clusterName, poolName).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot create pool with name: %s. Error: %w", poolName, err)
			}
			return pool, err
		})
	},
}

var poolUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Update cluster pool",
	ArgsUsage: "<pool_name>",
	Category:  "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-name",
			Aliases:  []string{"c"},
			Usage:    "Cluster name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "min-node-count",
			Usage:    "Minimum number of pool nodes",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "max-node-count",
			Usage:    "Maximum number of pool nodes",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "auto-healing-enabled",
			Usage:    "Enable/disable auto healing",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		poolName, err := flags.GetFirstStringArg(c, poolNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}
		clusterName := c.String("cluster-name")
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := pools.UpdateOpts{
			AutoHealingEnabled: c.Bool("auto-healing-enabled"),
			MinNodeCount:       c.Int("min-node-count"),
			MaxNodeCount:       c.Int("max-node-count"),
		}
		result, err := pools.Update(client, clusterName, poolName, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var poolResizeSubCommand = cli.Command{
	Name:      "resize",
	Usage:     "Resize cluster pool",
	ArgsUsage: "<pool_name>",
	Category:  "pool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-name",
			Aliases:  []string{"c"},
			Usage:    "Cluster name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "node-count",
			Usage:    "Node count",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		poolName, err := flags.GetFirstStringArg(c, poolNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "resize")
			return cli.NewExitError(err, 1)
		}
		clusterName := c.String("cluster-name")
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := pools.ResizeOpts{
			NodeCount: c.Int("node-count"),
		}
		results, err := pools.Resize(client, clusterName, poolName, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		tc := getTaskClient(c, client)
		return utils.WaitTaskAndShowResult(c, tc, results, true, func(task tasks.TaskID) (interface{}, error) {
			pool, err := pools.Get(client, clusterName, poolName).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot resize pool with name: %s. Error: %w", poolName, err)
			}
			return pool, err
		})
	},
}

var poolDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete cluster pool",
	ArgsUsage: "<pool_name>",
	Category:  "pool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-name",
			Aliases:  []string{"c"},
			Usage:    "Cluster name",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		poolName, err := flags.GetFirstStringArg(c, poolNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return cli.NewExitError(err, 1)
		}
		clusterName := c.String("cluster-name")
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := pools.Delete(client, clusterName, poolName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		tc := getTaskClient(c, client)
		return utils.WaitTaskAndShowResult(c, tc, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := pools.Get(client, clusterName, poolName).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete pool with name: %s. Error: %w", poolName, err)
			}
			switch err.(type) {
			case gcorecloud.ErrDefault404:
				return nil, nil
			default:
				return nil, err
			}
		})
	},
}

var poolInstancesSubCommand = cli.Command{
	Name:      "instances",
	Usage:     "List cluster pool instances",
	ArgsUsage: "<pool_name>",
	Category:  "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-name",
			Aliases:  []string{"c"},
			Usage:    "Cluster name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		poolName, err := flags.GetFirstStringArg(c, poolNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "instances")
			return cli.NewExitError(err, 1)
		}
		clusterName := c.String("cluster-name")
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := pools.ListInstancesAll(client, clusterName, poolName)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var poolCommands = cli.Command{
	Name:  "pool",
	Usage: "GCloud k8s cluster pool commands",
	Subcommands: []*cli.Command{
		&poolListSubCommand,
		&poolGetSubCommand,
		&poolCreateSubCommand,
		&poolUpdateSubCommand,
		&poolResizeSubCommand,
		&poolDeleteSubCommand,
		&poolInstancesSubCommand,
	},
}
