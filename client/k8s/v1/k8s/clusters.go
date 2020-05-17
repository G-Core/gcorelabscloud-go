package k8s

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/k8s/v1/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/clusters"

	"github.com/urfave/cli/v2"
)

var (
	clusterIDText = "cluster_id is mandatory argument"
)

func getPools(c *cli.Context) ([]pools.CreateOpts, error) {

	poolNames := c.StringSlice("pool-name")
	poolFlavors := c.StringSlice("flavor-id")
	poolNodesCount := c.IntSlice("node-count")
	poolMinNodesCount := c.IntSlice("min-node-count")
	poolMaxNodesCount := c.IntSlice("max-node-count")
	poolDockerVolumeSizes := c.IntSlice("docker-volume-size")

	if err := utils.ValidateEqualSlicesLength(poolNames, poolFlavors, poolNodesCount, poolMinNodesCount, poolMaxNodesCount); err != nil {
		return nil, fmt.Errorf("parameters number should be same for pool names, flavors, node-count, min-node-count and max_node_count: %w", err)
	}

	var result []pools.CreateOpts

	for idx, name := range poolNames {
		pool := pools.CreateOpts{
			Name:      name,
			FlavorID:  poolFlavors[idx],
			NodeCount: poolNodesCount[idx],
			DockerVolumeSize: func(idx int) int {
				if idx < len(poolDockerVolumeSizes) {
					return poolDockerVolumeSizes[idx]
				}
				return 0
			}(idx),
			MinNodeCount: poolMinNodesCount[idx],
			MaxNodeCount: poolMaxNodesCount[idx],
		}

		result = append(result, pool)

	}

	return result, nil

}

var clusterListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "K8s list clusters",
	Category: "cluster",
	Action: func(c *cli.Context) error {
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.ListAll(client, clusters.ListOpts{})
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var clusterGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "K8s get cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := clusters.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var clusterConfigSubCommand = cli.Command{
	Name:      "config",
	Usage:     "K8s get cluster config",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "save",
			Aliases:  []string{"s"},
			Usage:    "Save k8s config in file",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "force",
			Usage:    "Force rewrite KUBECONFIG file",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "merge",
			Aliases:  []string{"m"},
			Usage:    "Merge into existing KUBECONFIG file",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "file",
			Aliases:  []string{"c"},
			Usage:    "KUBECONFIG file",
			Value:    "~/.kube/config",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "config")
			return err
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		result, err := clusters.GetConfig(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if c.Bool("save") {
			filename := c.String("file")
			exists, err := utils.FileExists(filename)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "config")
				return cli.NewExitError(err, 1)
			}
			if exists {
				merge := c.Bool("merge")
				force := c.Bool("force")
				if (!force && !merge) || (force && merge) {
					_ = cli.ShowCommandHelp(c, "config")
					return cli.NewExitError(fmt.Errorf("either --force or --merge shoud be set"), 1)
				}
				if force {
					err := utils.WriteKubeconfigFile(filename, []byte(result.Config))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				}
				if merge {
					err := utils.MergeKubeconfigFile(filename, []byte(result.Config))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				}
			} else {
				err := utils.WriteToFile(filename, []byte(result.Config))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			}
		} else {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var clusterCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "K8s create cluster",
	Category: "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "cluster name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "master-node-count",
			Usage:    "master nodes count",
			Value:    1,
			Required: false,
		},
		&cli.StringFlag{
			Name:        "keypair",
			Aliases:     []string{"k"},
			Usage:       "The name of the SSH keypair",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:     "fixed-subnet",
			Usage:    "Fixed subnet that are using to allocate network address for nodes in cluster.",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "fixed-network",
			Usage:    "Fixed network that are using to allocate network address for nodes in cluster.",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "version",
			Value:    "",
			Usage:    "K8s cluster version",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "auto-healing-enabled",
			Usage:    "cluster auto healing",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "master-lb-floating-ip-enabled",
			Usage:    "use load balancer for K8s API",
			Required: false,
		},

		// pools parameters
		&cli.StringSliceFlag{
			Name:     "pool-name",
			Aliases:  []string{"p"},
			Usage:    "cluster pool names",
			Required: true,
		},
		&cli.IntSliceFlag{
			Name:     "node-count",
			Usage:    "pool nodes counts",
			Required: true,
		},
		&cli.IntSliceFlag{
			Name:     "docker-volume-size",
			Usage:    "docker volume size for pool nodes",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "flavor-id",
			Usage:    "pool node flavors",
			Required: true,
		},
		&cli.IntSliceFlag{
			Name:     "min-node-count",
			Usage:    "minimum number of pool nodes",
			Required: true,
		},
		&cli.IntSliceFlag{
			Name:     "max-node-count",
			Usage:    "maximum number of pool nodes",
			Required: true,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		clusterPools, err := getPools(c)

		if err != nil {
			return cli.NewExitError(err, 1)
		}

		masterCount := c.Int("master-node-count")
		masterLbFloatingIPEnabled := c.Bool("master-lb-floating-ip-enabled")
		if masterCount > 1 {
			masterLbFloatingIPEnabled = true
		}

		opts := clusters.CreateOpts{
			Name:                      c.String("name"),
			FixedNetwork:              c.String("fixed-network"),
			FixedSubnet:               c.String("fixed-subnet"),
			MasterCount:               masterCount,
			KeyPair:                   c.String("keypair"),
			AutoHealingEnabled:        c.Bool("auto-healing-enabled"),
			MasterLBFloatingIPEnabled: masterLbFloatingIPEnabled,
			Version:                   c.String("version"),
			Pools:                     clusterPools,
		}

		results, err := clusters.Create(client, opts).Extract()
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
			clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
			}
			cluster, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get cluster with ID: %s. Error: %w", clusterID, err)
			}
			utils.ShowResults(cluster, c.String("format"))
			return nil, nil
		})
	},
}

var clusterResizeSubCommand = cli.Command{
	Name:      "resize",
	Usage:     "K8s resize cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:     "node-count",
			Aliases:  []string{"n"},
			Usage:    "cluster nodes count",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:        "nodes-to-remove",
			Usage:       "cluster nodes chose to remove",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "pool",
			Aliases:     []string{"p"},
			Usage:       "cluster pool",
			DefaultText: "nil",
			Required:    false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "resize")
			return err
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		nodes := c.StringSlice("nodes-to-remove")
		if len(nodes) == 0 {
			nodes = nil
		}

		opts := clusters.ResizeOpts{
			NodeCount:     c.Int("node-count"),
			NodesToRemove: nodes,
			Pool:          c.String("pool"),
		}

		results, err := clusters.Resize(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
			}
			cluster, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get cluster with ID: %s. Error: %w", clusterID, err)
			}
			return cluster, nil
		})

	},
}

var clusterUpgradeSubCommand = cli.Command{
	Name:      "upgrade",
	Usage:     "K8s upgrade cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "version",
			Usage:    "cluster version",
			Required: true,
		},
		&cli.StringFlag{
			Name:        "pool",
			Usage:       "cluster pool",
			DefaultText: "nil",
			Required:    false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "upgrade")
			return err
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := clusters.UpgradeOpts{
			Pool:    c.String("pool"),
			Version: c.String("version"),
		}

		results, err := clusters.Upgrade(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
			}
			cluster, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get cluster with ID: %s. Error: %w", clusterID, err)
			}
			utils.ShowResults(cluster, c.String("format"))
			return nil, nil
		})

	},
}

var clusterVersionsSubCommand = cli.Command{
	Name:     "versions",
	Usage:    "K8s cluster versions",
	Category: "cluster",
	Action: func(c *cli.Context) error {
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		results, err := clusters.VersionsAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var clusterInstancesSubCommand = cli.Command{
	Name:      "instances",
	Usage:     "K8s cluster instances",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "instances")
			return err
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		results, err := clusters.InstancesAll(client, clusterID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var clusterVolumesSubCommand = cli.Command{
	Name:      "volumes",
	Usage:     "K8s cluster volumes",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "volumes")
			return err
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		results, err := clusters.VolumesAll(client, clusterID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var clusterDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "K8s delete cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.Delete(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := clusters.Get(client, clusterID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete cluster with ID: %s", clusterID)
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

var clusterCertificateSubCommand = cli.Command{
	Name:      "get",
	Usage:     "K8s get cluster CA certificate",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "get")
			return err
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		certificate, err := clusters.Certificate(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(certificate, c.String("format"))
		return nil

	},
}

var clusterSignCertificateSubCommand = cli.Command{
	Name:      "sign",
	Usage:     "K8s sign cluster CSR certificate request",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "csr",
			Aliases:  []string{"r"},
			Usage:    "cluster certificate sign request file",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "get")
			return err
		}

		data, err := utils.ReadFile(c.String("csr"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "get")
			return err
		}
		opts := clusters.ClusterSignCertificateOpts{
			CSR: string(data),
		}
		client, err := client.NewK8sClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		signedCertificate, err := clusters.SignCertificate(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(signedCertificate, c.String("format"))
		return nil

	},
}

var ClusterCommands = cli.Command{
	Name:  "cluster",
	Usage: "k8s cluster commands",
	Subcommands: []*cli.Command{
		&clusterListSubCommand,
		&clusterGetSubCommand,
		&clusterCreateSubCommand,
		&clusterUpgradeSubCommand,
		&clusterResizeSubCommand,
		&clusterDeleteSubCommand,
		&clusterConfigSubCommand,
		&clusterVersionsSubCommand,
		&clusterVolumesSubCommand,
		&clusterInstancesSubCommand,
		{
			Name:  "certificate",
			Usage: "K8s sign  certificates",
			Subcommands: []*cli.Command{
				&clusterCertificateSubCommand,
				&clusterSignCertificateSubCommand,
			},
		},
	},
}
