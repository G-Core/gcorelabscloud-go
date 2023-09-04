package k8s

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/k8s/v2/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/client/utils/k8sconfig"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
)

var (
	clusterNameText = "cluster_name is mandatory argument"
	volumeTypeNames = volumes.VolumeType("").StringList()
)

func getTaskClient(c *cli.Context, sc *gcorecloud.ServiceClient) *gcorecloud.ServiceClient {
	tc := &gcorecloud.ServiceClient{}
	*tc = *sc
	tc.Endpoint = strings.ReplaceAll(tc.Endpoint, "v2", "v1")
	return tc
}

func getPoolCreateOpts(c *cli.Context) ([]pools.CreateOpts, error) {
	names := c.StringSlice("pool-name")
	flavors := c.StringSlice("flavor-id")
	minNodeCounts := c.IntSlice("min-node-count")
	maxNodeCounts := c.IntSlice("max-node-count")
	bootVolumeSizes := c.IntSlice("boot-volume-size")
	bootVolumeTypes := utils.GetEnumStringSliceValue(c, "boot-volume-type")
	autoHealings := c.StringSlice("auto-healing-enabled")
	isPublicIPv4s := c.StringSlice("is-public-ipv4")

	if err := utils.ValidateEqualSlicesLength(names, flavors, minNodeCounts); err != nil {
		return nil, fmt.Errorf("parameters number should be the same for pool-name, flavor-id, min-node-count: %w", err)
	}

	var result []pools.CreateOpts
	for idx, name := range names {
		pool := pools.CreateOpts{
			Name:         name,
			FlavorID:     flavors[idx],
			MinNodeCount: minNodeCounts[idx],
			MaxNodeCount: func(idx int) int {
				if idx < len(maxNodeCounts) {
					return maxNodeCounts[idx]
				}
				return 1
			}(idx),
			BootVolumeSize: func(idx int) int {
				if idx < len(bootVolumeSizes) {
					return bootVolumeSizes[idx]
				}
				return 10
			}(idx),
			BootVolumeType: func(idx int) volumes.VolumeType {
				if idx < len(bootVolumeTypes) {
					return volumes.VolumeType(bootVolumeTypes[idx])
				}
				return "standard"
			}(idx),
			AutoHealingEnabled: func(idx int) bool {
				if idx < len(autoHealings) {
					b, _ := strconv.ParseBool(autoHealings[idx])
					return b
				}
				return false
			}(idx),
			IsPublicIPv4: func(idx int) bool {
				if idx < len(isPublicIPv4s) {
					b, _ := strconv.ParseBool(isPublicIPv4s[idx])
					return b
				}
				return false
			}(idx),
		}
		result = append(result, pool)
	}
	return result, nil
}

type configFileOptions struct {
	save     bool
	filename string
	exists   bool
	merge    bool
	force    bool
}

func (opts configFileOptions) check() error {
	if opts.exists && ((!opts.force && !opts.merge) || (opts.force && opts.merge)) {
		return fmt.Errorf("file %s exists, either --force or --merge shoud be set", opts.filename)
	}
	return nil
}

func getConfigFileOptions(c *cli.Context) (configFileOptions, error) {
	opts := configFileOptions{
		save:     c.Bool("save"),
		filename: c.String("file"),
		exists:   false,
		merge:    c.Bool("merge"),
		force:    c.Bool("force"),
	}
	if opts.save {
		var err error
		opts.exists, err = utils.FileExists(opts.filename)
		if err != nil {
			return opts, err
		}
		err = opts.check()
		if err != nil {
			return opts, err
		}
	}
	return opts, nil
}

var clusterListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "List clusters",
	Category: "cluster",
	Action: func(c *cli.Context) error {
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var clusterGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Get cluster information",
	ArgsUsage: "<cluster_name>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterName, err := flags.GetFirstStringArg(c, clusterNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := clusters.Get(client, clusterName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var clusterCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create cluster",
	Category: "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Cluster name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "fixed-network",
			Usage:    "Fixed network used to allocate network addresses for cluster nodes.",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "fixed-subnet",
			Usage:    "Fixed subnet used to allocate network addresses for cluster nodes.",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "keypair",
			Usage:    "SSH keypair name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "version",
			Usage:    "K8s version",
			Required: true,
		},
		&cli.StringFlag{
			Name:        "pods-ip-pool",
			Usage:       "Pods ip pool in CIDR notation",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "services-ip-pool",
			Usage:       "Services ip pool in CIDR notation",
			DefaultText: "nil",
			Required:    false,
		},

		// pools parameters
		&cli.StringSliceFlag{
			Name:     "pool-name",
			Usage:    "Pool names",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "flavor-id",
			Usage:    "Pool node flavors",
			Required: true,
		},
		&cli.IntSliceFlag{
			Name:     "min-node-count",
			Usage:    "Minimum numbers of pool nodes",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "max-node-count",
			Usage:    "Maximum numbers of pool nodes",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "boot-volume-size",
			Usage:    "Pool boot volume sizes in GB",
			Required: false,
		},
		&cli.GenericFlag{
			Name: "boot-volume-type",
			Value: &utils.EnumStringSliceValue{
				Enum:    volumeTypeNames,
				Default: volumeTypeNames[0],
			},
			Usage:    fmt.Sprintf("Pool noot volume types [%s]", strings.Join(volumeTypeNames, ", ")),
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "auto-healing-enabled",
			Usage:    "Enable/disable auto healing on pools",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "is-public-ipv4",
			Usage:    "Enable public IPv4 addresses on pools",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterName := c.String("name")
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		podsIPPool, err := gcorecloud.ParseCIDRStringOrNil(c.String("pods-ip-pool"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		servicesIPPool, err := gcorecloud.ParseCIDRStringOrNil(c.String("services-ip-pool"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		poolOpts, err := getPoolCreateOpts(c)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		opts := &clusters.CreateOpts{
			Name:           clusterName,
			FixedNetwork:   c.String("fixed-network"),
			FixedSubnet:    c.String("fixed-subnet"),
			PodsIPPool:     podsIPPool,
			ServicesIPPool: servicesIPPool,
			KeyPair:        c.String("keypair"),
			Version:        c.String("version"),
			Pools:          poolOpts,
		}
		results, err := clusters.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		tc := getTaskClient(c, client)
		return utils.WaitTaskAndShowResult(c, tc, results, true, func(task tasks.TaskID) (interface{}, error) {
			cluster, err := clusters.Get(client, clusterName).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot create cluster with name: %s. Error: %w", clusterName, err)
			}
			return cluster, nil
		})
	},
}

var clusterUpgradeSubCommand = cli.Command{
	Name:      "upgrade",
	Usage:     "Upgrade cluster",
	ArgsUsage: "<cluster_name>",
	Category:  "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "version",
			Usage:    "Target k8s version",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterName, err := flags.GetFirstStringArg(c, clusterNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "upgrade")
			return err
		}
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := clusters.UpgradeOpts{
			Version: c.String("version"),
		}
		results, err := clusters.Upgrade(client, clusterName, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		tc := getTaskClient(c, client)
		return utils.WaitTaskAndShowResult(c, tc, results, true, func(task tasks.TaskID) (interface{}, error) {
			cluster, err := clusters.Get(client, clusterName).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot upgrade with name: %s. Error: %w", clusterName, err)
			}
			return cluster, nil
		})
	},
}

var clusterDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete cluster",
	ArgsUsage: "<cluster_name>",
	Category:  "cluster",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		clusterName, err := flags.GetFirstStringArg(c, clusterNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.Delete(client, clusterName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		tc := getTaskClient(c, client)
		return utils.WaitTaskAndShowResult(c, tc, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := clusters.Get(client, clusterName).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete cluster with name: %s. Error: %w", clusterName, err)
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
	Name:      "certificate",
	Usage:     "Get cluster CA certificate",
	ArgsUsage: "<cluster_name>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterName, err := flags.GetFirstStringArg(c, clusterNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "certificate")
			return err
		}
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		certificate, err := clusters.GetCertificate(client, clusterName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(certificate, c.String("format"))
		return nil
	},
}

var clusterConfigSubCommand = cli.Command{
	Name:      "config",
	Usage:     "Get cluster kubeconfig",
	ArgsUsage: "<cluster_name>",
	Category:  "cluster",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "save",
			Aliases:  []string{"s"},
			Usage:    "Save KUBECONFIG to a file",
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
			Usage:    "Merge with existing KUBECONFIG file",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "file",
			Usage:    "KUBECONFIG file",
			EnvVars:  []string{"KUBECONFIG"},
			Value:    "~/.kube/config",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		clusterName, err := flags.GetFirstStringArg(c, clusterNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "config")
			return err
		}
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		options, err := getConfigFileOptions(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "config")
			return cli.NewExitError(err, 1)
		}
		result, err := clusters.GetConfig(client, clusterName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		config := strings.TrimSpace(result.Config)
		if options.save {
			if options.exists {
				if options.force {
					err := k8sconfig.WriteKubeconfigFile(options.filename, []byte(config))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				}
				if options.merge {
					err := k8sconfig.MergeKubeconfigFile(options.filename, []byte(config))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					return nil
				}
			} else {
				err := utils.WriteToFile(options.filename, []byte(config))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			}
		} else {
			fmt.Println(strings.TrimSpace(config))
		}
		return nil
	},
}

var clusterVersionsSubCommand = cli.Command{
	Name:     "versions",
	Usage:    "List supported k8s versions",
	Category: "cluster",
	Action: func(c *cli.Context) error {
		client, err := client.NewK8sClientV2(c)
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
	Usage:     "List cluster instances",
	ArgsUsage: "<cluster_name>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterName, err := flags.GetFirstStringArg(c, clusterNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "instances")
			return err
		}
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.ListInstancesAll(client, clusterName)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var Commands = cli.Command{
	Name:  "cluster",
	Usage: "GCloud k8s cluster API",
	Subcommands: []*cli.Command{
		&clusterListSubCommand,
		&clusterGetSubCommand,
		&clusterCreateSubCommand,
		&clusterUpgradeSubCommand,
		&clusterDeleteSubCommand,
		&clusterCertificateSubCommand,
		&clusterConfigSubCommand,
		&clusterVersionsSubCommand,
		&clusterInstancesSubCommand,
		&poolCommands,
	},
}
