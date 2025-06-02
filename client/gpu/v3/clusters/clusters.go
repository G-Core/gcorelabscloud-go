package clusters

import (
	"fmt"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/client"
	taskclient "github.com/G-Core/gcorelabscloud-go/client/tasks/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/urfave/cli/v2"
	"k8s.io/utils/pointer"

	"strings"
)

func showClusterAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "show")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	clusterDetails := clusters.Get(gpuClient, clusterID)
	if clusterDetails.Err != nil {
		return cli.Exit(clusterDetails.Err, 1)
	}

	utils.ShowResults(clusterDetails.Body, c.String("format"))
	return nil
}

func showVirtualClusterAction(c *cli.Context) error {
	return showClusterAction(c, client.NewGPUVirtualClientV3)
}

func showBaremetalClusterAction(c *cli.Context) error {
	return showClusterAction(c, client.NewGPUBaremetalClientV3)
}

func deleteClusterAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "delete")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	results, err := clusters.Delete(gpuClient, clusterID).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, results, false, func(task tasks.TaskID) (interface{}, error) {
		_, err := clusters.Get(gpuClient, clusterID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete GPU cluster with ID: %s. Error: %w", clusterID, err)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
}

func deleteVirtualClusterAction(c *cli.Context) error {
	return deleteClusterAction(c, client.NewGPUVirtualClientV3)
}

func deleteBaremetalClusterAction(c *cli.Context) error {
	return deleteClusterAction(c, client.NewGPUBaremetalClientV3)
}

func resizeVirtualClusterAction(c *cli.Context) error {
	return resizeClusterAction(c, client.NewGPUVirtualClientV3)
}

func softRebootVirtualClusterAction(c *cli.Context) error {
	return softRebootClusterAction(c, client.NewGPUVirtualClientV3)
}

func hardRebootVirtualClusterAction(c *cli.Context) error {
	return hardRebootClusterAction(c, client.NewGPUVirtualClientV3)
}

func startVirtualClusterAction(c *cli.Context) error {
	return startClusterAction(c, client.NewGPUVirtualClientV3)
}

func stopVirtualClusterAction(c *cli.Context) error {
	return stopClusterAction(c, client.NewGPUVirtualClientV3)
}

func updateTagsVirtualClusterAction(c *cli.Context) error {
	return updateTagsClusterAction(c, client.NewGPUVirtualClientV3)
}

func softRebootClusterAction(c *cli.Context, newClient func(ctx *cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "softreboot")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	results, err := clusters.SoftReboot(gpuClient, clusterID).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, results, false,
		waitForClusterOperation(gpuClient, clusterID))
}

func hardRebootClusterAction(c *cli.Context, newClient func(ctx *cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "hardreboot")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	results, err := clusters.HardReboot(gpuClient, clusterID).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, results, false,
		waitForClusterOperation(gpuClient, clusterID))
}

func startClusterAction(c *cli.Context, newClient func(ctx *cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "start")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	results, err := clusters.Start(gpuClient, clusterID).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, results, false,
		waitForClusterOperation(gpuClient, clusterID))
}

func stopClusterAction(c *cli.Context, newClient func(ctx *cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "stop")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	results, err := clusters.Stop(gpuClient, clusterID).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, results, false,
		waitForClusterOperation(gpuClient, clusterID))
}

func updateTagsClusterAction(c *cli.Context, newClient func(ctx *cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "updatetags")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	tags, err := utils.StringSliceToTags(c.StringSlice("tags"))
	results, err := clusters.UpdateTags(gpuClient, clusterID, tags).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, results, false,
		waitForClusterOperation(gpuClient, clusterID))
}

func waitForClusterOperation(gpuClient *gcorecloud.ServiceClient, clusterID string) func(task tasks.TaskID) (interface{}, error) {
	return func(task tasks.TaskID) (interface{}, error) {
		cluster, err := clusters.Get(gpuClient, clusterID).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot perform GPU cluster operation with ID: %s. Error: %w", clusterID, err)
		}
		return cluster, nil
	}
}

func resizeClusterAction(c *cli.Context, newClient func(ctx *cli.Context) (*gcorecloud.ServiceClient, error)) error {
	clusterID := c.Args().First()
	if clusterID == "" {
		_ = cli.ShowCommandHelp(c, "resize")
		return cli.Exit("cluster ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	// Validate servers count
	if c.Int("servers-count") <= 0 {
		return cli.Exit("`servers-count` must be greater than 0", 1)
	}

	results, err := clusters.Resize(gpuClient, clusterID, c.Int("servers-count")).Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, results, false,
		waitForClusterOperation(gpuClient, clusterID))
}

func createVirtualClusterAction(c *cli.Context) error {
	return createClusterAction(c, client.NewGPUVirtualClientV3)
}

func createClusterAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	// Validate mutually exclusive flags
	if c.IsSet("user-data") && (c.IsSet("server-username") || c.IsSet("server-password")) {
		return cli.Exit("`user-data` cannot be used together with `server-username` or `server-password`", 1)
	}

	// build create cluster options from CLI flags
	serverSettings, err := getServerSettings(c)
	if err != nil {
		_ = cli.ShowCommandHelp(c, "create")
		return cli.Exit(err, 1)
	}
	tags, err := utils.StringSliceToTags(c.StringSlice("tags"))
	if err != nil {
		_ = cli.ShowCommandHelp(c, "create")
		return cli.Exit(err, 1)
	}
	// Validate servers count
	if c.Int("servers-count") <= 0 {
		return cli.Exit("`servers-count` must be greater than 0", 1)
	}
	opts := clusters.CreateClusterOpts{
		Name:            c.String("name"),
		Flavor:          c.String("flavor"),
		ServersCount:    c.Int("servers-count"),
		Tags:            tags,
		ServersSettings: serverSettings,
	}

	// create the cluster and extract the task result
	result := clusters.Create(gpuClient, opts)
	if result.Err != nil {
		return cli.Exit(result.Err, 1)
	}
	taskResults, err := result.Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}
	utils.ShowResults(taskResults, c.String("format"))
	return nil
}

func getServerSettings(c *cli.Context) (clusters.ServerSettingsOpts, error) {
	interfaceOpts, err := getInterfaceOpts(c)
	if err != nil {
		return clusters.ServerSettingsOpts{}, err
	}
	volumeOpts, err := getVolumeOpts(c)
	if err != nil {
		return clusters.ServerSettingsOpts{}, err
	}
	credentialOpts := clusters.ServerCredentialsOpts{
		Username:   c.String("server-username"),
		Password:   c.String("server-password"),
		SSHKeyName: c.String("ssh-key-name"),
	}

	serverSettings := clusters.ServerSettingsOpts{
		Interfaces:     []clusters.InterfaceOpts{interfaceOpts},
		Volumes:        []clusters.VolumeOpts{volumeOpts},
		Credentials:    &credentialOpts,
		SecurityGroups: c.StringSlice("security-groups"),
		UserData:       StringPtrExcludeEmpty(c, "user-data"),
	}
	return serverSettings, nil
}

func StringPtrExcludeEmpty(c *cli.Context, name string) *string {
	if c.IsSet(name) && c.String(name) != "" {
		return pointer.StringPtr(c.String(name))
	}
	return nil
}

func getInterfaceOpts(c *cli.Context) (clusters.InterfaceOpts, error) {
	interfaceType := utils.GetEnumStringSliceValue(c, "interface-type")[0]
	interfaceName := StringPtrExcludeEmpty(c, "interface-name")

	sourceSlice := utils.GetEnumStringSliceValue(c, "interface-floating-source")
	var floatingIP *clusters.FloatingIPOpts
	if len(sourceSlice) > 0 {
		floatingIP = &clusters.FloatingIPOpts{Source: sourceSlice[0]}
	}

	switch clusters.InterfaceType(interfaceType) {
	case clusters.External:
		ipFamilySlice := utils.GetEnumStringSliceValue(c, "interface-ip-family")
		var ipFamily clusters.IPFamilyType
		if len(ipFamilySlice) > 0 {
			ipFamily = clusters.IPFamilyType(ipFamilySlice[0])
		}
		interfaceOpts := clusters.ExternalInterfaceOpts{
			Name:     interfaceName,
			Type:     interfaceType,
			IPFamily: ipFamily,
		}
		return interfaceOpts, nil
	case clusters.Subnet:
		interfaceOpts := clusters.SubnetInterfaceOpts{
			Name:       interfaceName,
			NetworkID:  c.String("interface-network-id"),
			Type:       interfaceType,
			SubnetID:   c.String("interface-subnet-id"),
			FloatingIP: floatingIP,
		}
		return interfaceOpts, nil
	case clusters.AnySubnet:
		ipFamilySlice := utils.GetEnumStringSliceValue(c, "interface-ip-family")
		var ipFamily clusters.IPFamilyType
		if len(ipFamilySlice) > 0 {
			ipFamily = clusters.IPFamilyType(ipFamilySlice[0])
		}
		interfaceOpts := clusters.AnySubnetInterfaceOpts{
			Name:       interfaceName,
			NetworkID:  c.String("interface-network-id"),
			Type:       interfaceType,
			FloatingIP: floatingIP,
			IPAddress:  StringPtrExcludeEmpty(c, "interface-ip-address"),
			IPFamily:   ipFamily,
		}
		return interfaceOpts, nil
	}
	return nil, fmt.Errorf("unexpected interface-type: %v", interfaceType)
}

func getVolumeOpts(c *cli.Context) (clusters.VolumeOpts, error) {
	volumeType := utils.GetEnumStringSliceValue(c, "volume-type")[0]
	volume := clusters.VolumeOpts{
		Source:              clusters.Image,
		Name:                c.String("volume-name"),
		BootIndex:           0,
		DeleteOnTermination: c.Bool("volume-delete-on-termination"),
		Size:                c.Int("volume-size"),
		Type:                clusters.VolumeType(volumeType),
		ImageID:             c.String("volume-image-id"),
	}
	if c.IsSet("volume-tags") {
		tags, err := utils.StringSliceToTags(c.StringSlice("volume-tags"))
		if err != nil {
			return clusters.VolumeOpts{}, err
		}
		volume.Tags = tags
	}
	return volume, nil
}

func createClusterFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "name of the cluster",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "flavor",
			Usage:    "flavor ID of the cluster",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "tags",
			Aliases:  []string{"t"},
			Usage:    "cluster key-value tags. Example: --tags key1=value1 --tags key2=value2",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "servers-count",
			Aliases:  []string{"sc"},
			Usage:    "number of servers of the cluster",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "security-groups",
			Aliases:  []string{"sg"},
			Usage:    "security groups IDs of the cluster. Example: --security-groups b4849ffa-89f2-45a1-951f-0ae5b7809d98 --security-groups d478ae29-dedc-4869-82f0-96104425f565",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "user-data",
			Aliases:  []string{"ud"},
			Usage:    "user data for the cluster (Base64 encoded string)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "server-username",
			Aliases:  []string{"u"},
			Usage:    "username for the servers in the cluster",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "server-password",
			Aliases:  []string{"p"},
			Usage:    "password for the servers in the cluster",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "ssh-key-name",
			Aliases:  []string{"k"},
			Usage:    "(ssh) keypair name for the servers in the cluster",
			Required: false,
		},
		&cli.BoolFlag{
			Name:  "volume-delete-on-termination",
			Usage: "delete volume on termination",
		},
		&cli.StringFlag{
			Name:     "volume-name",
			Aliases:  []string{"vn"},
			Usage:    "name of the volume",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "volume-size",
			Usage:    "size of the volume",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "volume-image-id",
			Usage:    "image ID of the volume",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:  "volume-tags",
			Usage: "tags for the volume",
		},
		&cli.GenericFlag{
			Name:    "volume-type",
			Aliases: []string{"vt"},
			Value: &utils.EnumStringSliceValue{
				Enum: clusters.VolumeTypesStringList(),
			},
			Usage: fmt.Sprintf("volume types. One of %s",
				strings.Join(clusters.VolumeTypesStringList(), ", ")),
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "interface-type",
			Aliases: []string{"it"},
			Value: &utils.EnumStringSliceValue{
				Enum: clusters.InterfaceTypeStringList(),
			},
			Usage: fmt.Sprintf("interface type. One of %s",
				strings.Join(clusters.InterfaceTypeStringList(), ", ")),
			Required: true,
		},
		&cli.StringFlag{
			Name:     "interface-name",
			Usage:    "name of the interface",
			Required: false,
		},
		&cli.GenericFlag{
			Name: "interface-ip-family",
			Value: &utils.EnumStringSliceValue{
				Enum: clusters.IPFamilyTypeListStringList(),
			},
			Usage: fmt.Sprintf("IP family of the interface. One of %s",
				strings.Join(clusters.IPFamilyTypeListStringList(), ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:  "interface-network-id",
			Usage: "network ID of the interface",
		},
		&cli.StringFlag{
			Name:  "interface-subnet-id",
			Usage: "subnet ID of the interface",
		},
		&cli.StringFlag{
			Name:  "interface-ip-address",
			Usage: "IP address of the interface",
		},
		&cli.GenericFlag{
			Name:    "interface-floating-source",
			Aliases: []string{"ifs"},
			Value: &utils.EnumStringSliceValue{
				Enum: clusters.FloatingIPSourceStringList(),
			},
			Usage: fmt.Sprintf("floating ip source. One of %s",
				strings.Join(clusters.FloatingIPSourceStringList(), ", ")),
			Required: false,
		},
	}
}

// listClustersAction handles the common logic for listing both virtual and baremetal clusters
func listClustersAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}
	opts := &clusters.ListOpts{}
	clusterList, err := clusters.ListAll(gpuClient, opts)
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(clusterList, c.String("format"))
	return nil
}

func listVirtualClustersAction(c *cli.Context) error {
	return listClustersAction(c, client.NewGPUVirtualClientV3)
}

func listBaremetalClustersAction(c *cli.Context) error {
	return listClustersAction(c, client.NewGPUBaremetalClientV3)
}

// BaremetalCommands returns commands for managing baremetal GPU clusters
func BaremetalCommands() *cli.Command {
	return &cli.Command{
		Name:        "clusters",
		Usage:       "Manage baremetal GPU clusters",
		Description: "Commands for managing baremetal GPU clusters",
		Subcommands: []*cli.Command{
			{
				Name:        "show",
				Usage:       "Show baremetal GPU cluster details",
				Description: "Show details of a specific baremetal GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      showBaremetalClusterAction,
			},
			{
				Name:        "delete",
				Usage:       "Delete baremetal GPU cluster",
				Description: "Delete a specific baremetal GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      deleteBaremetalClusterAction,
				Flags:       flags.WaitCommandFlags,
			},
			{
				Name:        "list",
				Usage:       "List baremetal GPU clusters",
				Description: "List all baremetal GPU clusters",
				Category:    "clusters",
				ArgsUsage:   " ",
				Action:      listBaremetalClustersAction,
			},
		},
	}
}

// VirtualCommands returns commands for managing virtual GPU clusters
func VirtualCommands() *cli.Command {
	return &cli.Command{
		Name:        "clusters",
		Usage:       "Manage virtual GPU clusters",
		Description: "Commands for managing virtual GPU clusters",
		Subcommands: []*cli.Command{
			{
				Name:        "show",
				Usage:       "Show virtual GPU cluster details",
				Description: "Show details of a specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      showVirtualClusterAction,
			},
			{
				Name:        "delete",
				Usage:       "Delete virtual GPU cluster",
				Description: "Delete a specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      deleteVirtualClusterAction,
				Flags:       flags.WaitCommandFlags,
			},
			{
				Name:        "create",
				Usage:       "Create a new virtual GPU cluster",
				Description: "Create a new virtual GPU cluster with the specified options",
				Category:    "clusters",
				Flags:       append(createClusterFlags(), flags.WaitCommandFlags...),
				Action:      createVirtualClusterAction,
			},
			{
				Name:        "list",
				Usage:       "List virtual GPU clusters",
				Description: "List all virtual GPU clusters",
				Category:    "clusters",
				ArgsUsage:   " ",
				Action:      listVirtualClustersAction,
			},
			{
				Name:        "resize",
				Usage:       "Resize virtual GPU cluster",
				Description: "Resize a specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      resizeVirtualClusterAction,
				Flags: append([]cli.Flag{
					&cli.IntFlag{
						Name:     "servers-count",
						Aliases:  []string{"sc"},
						Usage:    "number of servers of the cluster",
						Required: true,
					},
				},
					flags.WaitCommandFlags...),
			},
			{
				Name:        "softreboot",
				Usage:       "Soft reboot virtual GPU cluster",
				Description: "Soft reboot of specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      softRebootVirtualClusterAction,
				Flags:       flags.WaitCommandFlags,
			},
			{
				Name:        "hardreboot",
				Usage:       "Hard reboot virtual GPU cluster",
				Description: "Hard reboot of specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      hardRebootVirtualClusterAction,
				Flags:       flags.WaitCommandFlags,
			},
			{
				Name:        "start",
				Usage:       "Start (power on) virtual GPU cluster",
				Description: "Power on a specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      startVirtualClusterAction,
				Flags:       flags.WaitCommandFlags,
			},
			{
				Name:        "stop",
				Usage:       "Stop (power off) virtual GPU cluster",
				Description: "Power off a specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      stopVirtualClusterAction,
				Flags:       flags.WaitCommandFlags,
			},
			{
				Name:        "updatetags",
				Usage:       "Update tags of virtual GPU cluster",
				Description: "Updates the tags of specific virtual GPU cluster",
				Category:    "clusters",
				ArgsUsage:   "<cluster_id>",
				Action:      updateTagsVirtualClusterAction,
				Flags: append([]cli.Flag{
					&cli.StringSliceFlag{
						Name:     "tags",
						Aliases:  []string{"t"},
						Usage:    "cluster key-value tags. Example: --tags key1=value1 --tags key2=value2",
						Required: false,
					},
				}, flags.WaitCommandFlags...),
			},
		},
	}
}
