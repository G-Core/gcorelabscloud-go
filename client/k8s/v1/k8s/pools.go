package k8s

import (
	"fmt"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/k8s/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"
)

var (
	poolIDText = "pool_id is mandatory argument"
)

var poolListSubCommand = cli.Command{
	Name:      "list",
	Usage:     "K8s list pools",
	Category:  "pool",
	ArgsUsage: "<cluster_id>",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return cli.NewExitError(err, 1)
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := pools.ListAll(client, clusterID, pools.ListOpts{})
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var poolDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "K8s delete pool",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-id",
			Aliases:  []string{"c"},
			Usage:    "K8s cluster ID",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		poolID, err := flags.GetFirstStringArg(c, poolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return cli.NewExitError(err, 1)
		}
		clusterID := c.String("cluster-id")
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := pools.Delete(client, clusterID, poolID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := pools.Get(client, clusterID, poolID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete pool with ID: %s", poolID)
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

var poolUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "K8s update pool",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-id",
			Aliases:  []string{"c"},
			Usage:    "K8s cluster ID",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "cluster pool name",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "min-node-count",
			Usage:    "minimum number of pool nodes",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "max-node-count",
			Usage:    "maximum number of pool nodes",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		poolID, err := flags.GetFirstStringArg(c, poolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}
		clusterID := c.String("cluster-id")
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := pools.UpdateOpts{
			Name:         c.String("name"),
			MinNodeCount: c.Int("min-node-count"),
			MaxNodeCount: c.Int("max-node-count"),
		}

		result, err := pools.Update(client, clusterID, poolID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var poolGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "K8s pool show",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-id",
			Aliases:  []string{"c"},
			Usage:    "K8s cluster ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		poolID, err := flags.GetFirstStringArg(c, poolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return cli.NewExitError(err, 1)
		}
		clusterID := c.String("cluster-id")
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := pools.Get(client, clusterID, poolID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var poolCreateSubCommand = cli.Command{
	Name:      "create",
	Usage:     "K8s create pool",
	Category:  "pool",
	ArgsUsage: "<cluster_id>",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "cluster pool name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "node-count",
			Usage:    "nodes count",
			Aliases:  []string{"c"},
			Value:    1,
			Required: false,
		},
		&cli.StringFlag{
			Name:     "flavor-id",
			Usage:    "node flavor",
			Required: true,
		},
		&cli.IntFlag{
			Name:        "docker-volume-size",
			Usage:       "The size in GB for the local storage on each server for the Docker daemon to cache the images and host the containers",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.GenericFlag{
			Name: "docker-volume-type",
			Value: &utils.EnumValue{
				Enum:    volumeTypeNames,
				Default: volumeTypeNames[0],
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(volumeTypeNames, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "min-node-count",
			Usage:    "minimum number of pool nodes",
			Value:    1,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "max-node-count",
			Usage:    "maximum number of pool nodes",
			Value:    5,
			Required: false,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		dockerVolumeType := volumes.VolumeType(c.String("docker-volume-type"))

		opts := pools.CreateOpts{
			Name:             c.String("name"),
			FlavorID:         c.String("flavor-id"),
			NodeCount:        c.Int("node-count"),
			DockerVolumeSize: c.Int("docker-volume-size"),
			MinNodeCount:     c.Int("min-node-count"),
			MaxNodeCount:     c.Int("max-node-count"),
			DockerVolumeType: dockerVolumeType,
		}

		results, err := pools.Create(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			poolID, err := pools.ExtractClusterPoolIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve pool ID from task info: %w", err)
			}
			pool, err := pools.Get(client, clusterID, poolID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get pool with ID: %s. Error: %w", poolID, err)
			}
			utils.ShowResults(pool, c.String("format"))
			return nil, nil
		})
	},
}

var poolInstancesSubCommand = cli.Command{
	Name:      "instances",
	Usage:     "K8s cluster pool instances",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-id",
			Aliases:  []string{"c"},
			Usage:    "K8s cluster ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		poolID, err := flags.GetFirstStringArg(c, poolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "instances")
			return cli.NewExitError(err, 1)
		}
		clusterID := c.String("cluster-id")
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		results, err := pools.InstancesAll(client, clusterID, poolID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var poolVolumesSubCommand = cli.Command{
	Name:      "volumes",
	Usage:     "K8s cluster pool volumes",
	ArgsUsage: "<pool_id>",
	Category:  "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "cluster-id",
			Aliases:  []string{"c"},
			Usage:    "K8s cluster ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		poolID, err := flags.GetFirstStringArg(c, poolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "volumes")
			return cli.NewExitError(err, 1)
		}
		clusterID := c.String("cluster-id")
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		results, err := pools.VolumesAll(client, clusterID, poolID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var ClusterPoolCommands = cli.Command{
	Name:  "pool",
	Usage: "Gcloud K8s pool commands",
	Subcommands: []*cli.Command{
		&poolListSubCommand,
		&poolDeleteSubCommand,
		&poolGetSubCommand,
		&poolCreateSubCommand,
		&poolUpdateSubCommand,
		&poolInstancesSubCommand,
		&poolVolumesSubCommand,
	},
}
