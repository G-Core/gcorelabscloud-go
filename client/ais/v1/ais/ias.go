package ais

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/ais/v1/client"
	client2 "github.com/G-Core/gcorelabscloud-go/client/ais/v2/client"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	instance_client "github.com/G-Core/gcorelabscloud-go/client/instances/v1/instances"
	task_client "github.com/G-Core/gcorelabscloud-go/client/tasks/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	cmeta "github.com/G-Core/gcorelabscloud-go/client/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/aiflavors"
	"github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/aiimages"
	ai "github.com/G-Core/gcorelabscloud-go/gcore/ai/v1/ais"
	img_types "github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
)

var (
	aiInstanceIDText          = "instance_id is mandatory argument"
	aiClusterIDText           = "cluster_id is mandatory argument"
	volumeSourceType          = types.VolumeSource("").StringList()
	volumeType                = volumes.VolumeType("").StringList()
	interfaceTypes            = types.InterfaceType("").StringList()
	interfaceFloatingIPSource = types.FloatingIPSource("").StringList()
	visibilityTypes           = img_types.Visibility("").StringList()
	bootableIndex             = 0
)

// Command declarations
var aiClusterRebuildCommand = cli.Command{
	Name: "rebuild",
	Usage: `
	Rebuild GPU AI cluster nodes
	Example: gcoreclient ai rebuild --image 06e62653-1f88-4d38-9aa6-62833e812b4f --user-data "test" --nodes node1 --nodes node2 <cluster_id> -d -w`,
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "image",
			Usage:    "AI cluster image",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "user-data",
			Usage:    "AI cluster user data",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "user-data-file",
			Usage:    "instance user data file",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "nodes",
			Usage:    "List of node IDs to rebuild. Example: --nodes node1 --nodes node2",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "rebuild")
			return err
		}
		client, err := client.NewAIGPUClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		userData, err := instance_client.GetUserData(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "rebuild")
			return cli.NewExitError(err, 1)
		}

		opts := ai.RebuildGPUAIClusterOpts{
			ImageID:  c.String("image"),
			UserData: userData,
			Nodes:    c.StringSlice("nodes"),
		}

		err = gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		results, err := ai.RebuildGPUAICluster(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		taskClient, err := task_client.NewTaskClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		clusterClientV2, err := client2.NewAIClusterClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, taskClient, results, c.Bool("d"), func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(taskClient, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			// on rebuild Task the cluster_id is inside the 'data' section
			clusterID, ok := (*taskInfo.Data)["cluster_id"]
			if !ok {
				return nil, fmt.Errorf("cannot get cluster_id from task data section: %+v", taskInfo.Data)
			}
			cluster, err := ai.Get(clusterClientV2, clusterID.(string)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get AI cluster with ID: %s. Error: %w", clusterID, err)
			}
			return cluster, nil
		})
	},
}

var Commands = cli.Command{
	Name:  "ai",
	Usage: "GCloud AI API",
	Subcommands: []*cli.Command{
		&aiClusterGetCommand,
		&aiClustersListCommand,
		&aiClusterCreateCommand,
		&aiClusterDeleteCommand,
		&aiClusterPowerCycleCommand,
		&aiClusterRebuildCommand,
		&aiClusterRebootCommand,
		&aiClusterSuspendCommand,
		&aiClusterResumeCommand,
		&aiClusterResizeCommand,
		{
			Name:     "interface",
			Usage:    "AI cluster interface action",
			Category: "cluster",
			Subcommands: []*cli.Command{
				&aiClusterListInterfacesCommand,
				&aiClusterAttachInterfacesCommand,
				&aiClusterDetachInterfacesCommand,
			},
		},
		{
			Name:     "port",
			Usage:    "AI cluster port action",
			Category: "cluster",
			Subcommands: []*cli.Command{
				&aiClusterListPortsCommand,
			},
		},
		{
			Name:     "instance",
			Usage:    "AI instances action",
			Category: "cluster",
			Subcommands: []*cli.Command{
				&aiInstancePowerCycleCommand,
				&aiInstanceRebootCommand,
				&aiInstanceGetConsoleCommand,
			},
		},
		{
			Name:     "securitygroup",
			Usage:    "AI cluster security groups",
			Category: "cluster",
			Subcommands: []*cli.Command{
				&aiClusterAssignSecurityGroupsCommand,
				&aiClusterUnAssignSecurityGroupsCommand,
			},
		},
		{
			Name:     "image",
			Usage:    "AI cluster available images",
			Category: "cluster",
			Subcommands: []*cli.Command{
				&aiClusterAvailableImagesCommand,
			},
		},
		{
			Name:     "flavor",
			Usage:    "AI cluster available flavors",
			Category: "cluster",
			Subcommands: []*cli.Command{
				&aiClusterAvailableFlavorsCommand,
			},
		},
		{
			Name:     "metadata",
			Usage:    "AI cluster metadata",
			Category: "AI cluster metadata",
			Subcommands: []*cli.Command{
				cmeta.NewMetadataListCommand(
					client.NewAIClusterClientV1,
					"Get AI cluster metadata",
					"<cluster_id>",
					"cluster_id is mandatory argument",
				),
				cmeta.NewMetadataGetCommand(
					client2.NewAIClusterClientV2,
					"Show AI cluster metadata by key",
					"<cluster_id>",
					"cluster_id is mandatory argument",
				),
				cmeta.NewMetadataDeleteCommand(
					client2.NewAIClusterClientV2,
					"Delete AI cluster metadata by key",
					"<cluster_id>",
					"cluster_id is mandatory argument",
				),
				cmeta.NewMetadataCreateCommand(
					client2.NewAIClusterClientV2,
					"Create AI cluster metadata. It would update existing keys",
					"<cluster_id>",
					"cluster_id is mandatory argument",
				),
				cmeta.NewMetadataUpdateCommand(
					client2.NewAIClusterClientV2,
					"Update AI cluster metadata. It overriding existing records",
					"<cluster_id>",
					"cluster_id is mandatory argument",
				),
				cmeta.NewMetadataReplaceCommand(
					client2.NewAIClusterClientV2,
					"Replace AI cluster metadata. It replace existing records",
					"<cluster_id>",
					"cluster_id is mandatory argument",
				),
			},
		},
	},
}

var aiClustersListCommand = cli.Command{
	Name:     "list",
	Usage:    "List ai clusters",
	Category: "cluster",
	Action: func(c *cli.Context) error {
		client, err := client2.NewAIClusterClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := ai.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var aiClusterListInterfacesCommand = cli.Command{
	Name:      "list",
	Usage:     "List ai cluster interfaces",
	ArgsUsage: "<cluster_id>",
	Category:  "interface",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := ai.ListInterfacesAll(client, clusterID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var aiClusterAttachInterfacesCommand = cli.Command{
	Name:      "attach",
	Usage:     "attach interface to AI instance",
	ArgsUsage: "<cluster_id>",
	Category:  "interface",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "type",
			Aliases:  []string{"t"},
			Usage:    "interface type",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "network",
			Aliases:  []string{"n"},
			Usage:    "interface network id",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "subnet",
			Aliases:  []string{"s"},
			Usage:    "interface subnet id",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "port",
			Aliases:  []string{"p"},
			Usage:    "interface port id",
			Required: false,
		},
	},
	),
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, aiInstanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "attach")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := ai.AttachInterfaceOpts{
			Type:      types.InterfaceType(c.String("type")),
			NetworkID: c.String("network"),
			SubnetID:  c.String("subnet"),
			PortID:    c.String("port"),
		}

		results, err := ai.AttachAIInstanceInterface(client, instanceID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var aiClusterDetachInterfacesCommand = cli.Command{
	Name:      "detach",
	Usage:     "detach interface to AI instance",
	ArgsUsage: "<cluster_id>",
	Category:  "interface",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "port",
			Aliases:  []string{"p"},
			Usage:    "interface port id",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "ip-address",
			Aliases:  []string{"ip"},
			Usage:    "interface ip address id",
			Required: true,
		},
	},
	),
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, aiInstanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "attach")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := ai.DetachInterfaceOpts{
			PortID:    c.String("port"),
			IpAddress: c.String("ip-address"),
		}

		results, err := ai.DetachAIInstanceInterface(client, instanceID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var aiClusterListPortsCommand = cli.Command{
	Name:      "list",
	Usage:     "List ai cluster ports",
	ArgsUsage: "<cluster_id>",
	Category:  "port",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := ai.ListPortsAll(client, clusterID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var aiClusterAssignSecurityGroupsCommand = cli.Command{
	Name:      "add",
	Usage:     "Add AI cluster security group",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "security group name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "add")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := instances.SecurityGroupOpts{Name: c.String("name")}

		err = ai.AssignSecurityGroup(client, clusterID, opts).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var aiClusterUnAssignSecurityGroupsCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete AI cluster security group",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "security group name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := instances.SecurityGroupOpts{Name: c.String("name")}

		err = ai.UnAssignSecurityGroup(client, instanceID, opts).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

func StringSliceToMetadata(slice []string) (map[string]string, error) {
	if len(slice) == 0 {
		return nil, nil
	}
	m := make(map[string]string, len(slice))
	for _, s := range slice {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return m, fmt.Errorf("wrong label format: %s", s)
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}

var aiClusterCreateCommand = cli.Command{
	Name: "create",
	Usage: `
	Create AI cluster
	Example: gcoreclient ai create --flavor g2a-ai-fake-v1pod-8 --image 256f6681-49bf-449f-8078-6ff6ea771ef4 --keypair sshkey --it external --volume-type standard --volume-source image  --volume-image-id 256f6681-49bf-449f-8078-6ff6ea771ef4 --volume-size  20 --name aicluster -d -w`,
	Category: "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "AI cluster name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "flavor",
			Usage:    "AI cluster flavor",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "image",
			Usage:    "AI cluster image",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "keypair",
			Aliases:  []string{"k"},
			Usage:    "AI cluster ssh keypair",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "AI cluster password",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "username",
			Aliases:  []string{"u"},
			Usage:    "AI cluster username",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "user-data",
			Usage:    "AI cluster user data",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "user-data-file",
			Usage:    "instance user data file",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "volume-source",
			Aliases: []string{"vs"},
			Value: &utils.EnumStringSliceValue{
				Enum: volumeSourceType,
			},
			Usage:    fmt.Sprintf("instance volume source. output in %s", strings.Join(volumeSourceType, ", ")),
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "volume-boot-index",
			Usage:    "instance volume boot index",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "volume-size",
			Usage:    "instance volume size",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "volume-type",
			Aliases: []string{"vt"},
			Value: &utils.EnumStringSliceValue{
				Enum: volumeType,
			},
			Usage:    fmt.Sprintf("instance volume types. output in %s", strings.Join(volumeType, ", ")),
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "volume-name",
			Usage:    "instance volume name",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "volume-image-id",
			Usage:    "instance volume image id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "volume-snapshot-id",
			Usage:    "instance volume snapshot id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "volume-volume-id",
			Usage:    "instance volume volume id",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "interface-type",
			Aliases: []string{"it"},
			Value: &utils.EnumStringSliceValue{
				Enum: interfaceTypes,
			},
			Usage:    fmt.Sprintf("instance interface type. output in %s", strings.Join(interfaceTypes, ", ")),
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "interface-network-id",
			Usage:    "instance interface network id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "interface-subnet-id",
			Usage:    "instance interface subnet id",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "interface-floating-source",
			Aliases: []string{"ifs"},
			Value: &utils.EnumStringSliceValue{
				Enum: interfaceFloatingIPSource,
			},
			Usage:    fmt.Sprintf("instance floating ip source. output in %s", strings.Join(interfaceFloatingIPSource, ", ")),
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "interface-floating-ip",
			Usage:    "instance interface existing floating ip. Required when --interface-floating-source set as `existing`",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "security-group",
			Usage:    "instance security group",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "metadata",
			Usage:    "instance metadata. Example: --metadata one=two --metadata three=four",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "gpu",
			Usage:    "create gpu cluster else not gpu cluster",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "instances-count",
			Usage:    "number of instances to create",
			Required: false,
			Value:    1,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterClient, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		if c.Bool("gpu") {
			clusterClient, err = client.NewAIGPUClusterClientV1(c)
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
		}
		if c.Int("instances-count") < 1 {
			return cli.Exit("instances-count must be greater or equal than 1", 1)
		}

		clusterClientV2, err := client2.NewAIClusterClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		userData, err := instance_client.GetUserData(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		instanceVolumes, err := instance_client.GetInstanceVolumes(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		// todo add security group mapping
		instanceInterfaces, err := instance_client.GetInterfaces(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		securityGroups := instance_client.GetSecurityGroups(c)

		metadata, err := StringSliceToMetadata(c.StringSlice("metadata"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		opts := ai.CreateOpts{
			Flavor:         c.String("flavor"),
			Name:           c.String("name"),
			ImageID:        c.String("image"),
			Volumes:        instanceVolumes,
			Interfaces:     instanceInterfaces,
			SecurityGroups: securityGroups,
			Keypair:        c.String("keypair"),
			Password:       c.String("password"),
			Username:       c.String("username"),
			UserData:       userData,
			Metadata:       metadata,
			InstancesCount: c.Int("instances-count"),
		}

		err = gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		results, err := ai.Create(clusterClient, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		taskClient, err := task_client.NewTaskClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, taskClient, results, c.Bool("d"), func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(taskClient, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			// on create cluster, the cluster_id is the same as the Task id
			clusterID := taskInfo.ID
			cluster, err := ai.Get(clusterClientV2, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get AI cluster with ID: %s. Error: %w", clusterID, err)
			}
			return cluster, nil
		})
	},
}

var aiClusterResizeCommand = cli.Command{
	Name: "resize",
	Usage: `
	Resize AI cluster
	Example: token ai resize --flavor g2a-ai-fake-v1pod-8 --image 06e62653-1f88-4d38-9aa6-62833e812b4f --keypair sshkey --it any_subnet --interface-network-id 518ba531-496b-4676-8ea4-68e2ed3b2e4b --interface-floating-source new --volume-type standard --volume-source image  --volume-image-id 06e62653-1f88-4d38-9aa6-62833e812b4f --volume-size  20  e673bba0-fcef-44d9-904c-824546b608ec -d -w`,
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "flavor",
			Usage:    "AI cluster flavor",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "image",
			Usage:    "AI cluster image",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "keypair",
			Aliases:  []string{"k"},
			Usage:    "AI cluster ssh keypair",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "AI cluster password",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "username",
			Aliases:  []string{"u"},
			Usage:    "AI cluster username",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "user-data",
			Usage:    "AI cluster user data",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "user-data-file",
			Usage:    "instance user data file",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "volume-source",
			Aliases: []string{"vs"},
			Value: &utils.EnumStringSliceValue{
				Enum: volumeSourceType,
			},
			Usage:    fmt.Sprintf("instance volume source. output in %s", strings.Join(volumeSourceType, ", ")),
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "volume-boot-index",
			Usage:    "instance volume boot index",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "volume-size",
			Usage:    "instance volume size",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "volume-type",
			Aliases: []string{"vt"},
			Value: &utils.EnumStringSliceValue{
				Enum: volumeType,
			},
			Usage:    fmt.Sprintf("instance volume types. output in %s", strings.Join(volumeType, ", ")),
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "volume-name",
			Usage:    "instance volume name",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "volume-image-id",
			Usage:    "instance volume image id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "volume-snapshot-id",
			Usage:    "instance volume snapshot id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "volume-volume-id",
			Usage:    "instance volume volume id",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "interface-type",
			Aliases: []string{"it"},
			Value: &utils.EnumStringSliceValue{
				Enum: interfaceTypes,
			},
			Usage:    fmt.Sprintf("instance interface type. output in %s", strings.Join(interfaceTypes, ", ")),
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "interface-network-id",
			Usage:    "instance interface network id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "interface-subnet-id",
			Usage:    "instance interface subnet id",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "interface-floating-source",
			Aliases: []string{"ifs"},
			Value: &utils.EnumStringSliceValue{
				Enum: interfaceFloatingIPSource,
			},
			Usage:    fmt.Sprintf("instance floating ip source. output in %s", strings.Join(interfaceFloatingIPSource, ", ")),
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "interface-floating-ip",
			Usage:    "instance interface existing floating ip. Required when --interface-floating-source set as `existing`",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "security-group",
			Usage:    "instance security group",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "metadata",
			Usage:    "instance metadata. Example: --metadata one=two --metadata three=four",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		userData, err := instance_client.GetUserData(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		instanceVolumes, err := instance_client.GetInstanceVolumes(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		// todo add security group mapping
		instanceInterfaces, err := instance_client.GetInterfaces(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		securityGroups := instance_client.GetSecurityGroups(c)

		metadata, err := StringSliceToMetadata(c.StringSlice("metadata"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		opts := ai.ResizeAIClusterOpts{
			Flavor:         c.String("flavor"),
			ImageID:        c.String("image"),
			Volumes:        instanceVolumes,
			Interfaces:     instanceInterfaces,
			SecurityGroups: securityGroups,
			Keypair:        c.String("keypair"),
			Password:       c.String("password"),
			Username:       c.String("username"),
			UserData:       userData,
			Metadata:       metadata,
		}

		err = gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		results, err := ai.Resize(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, c.Bool("d"), func(task tasks.TaskID) (interface{}, error) {
			cluster, err := ai.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get AI cluster with ID: %s. Error: %w", clusterID, err)
			}
			return cluster, nil
		})
	},
}

var aiClusterGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get ai cluster information",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client2.NewAIClusterClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		instance, err := ai.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var aiClusterDeleteCommand = cli.Command{
	Name:  "delete",
	Usage: "Delete AI cluster",
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:     "volume-id",
			Usage:    "instance volume id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "floating-ip",
			Usage:    "delete selected cluster floating ips",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "delete-floating-ips",
			Usage:    "delete all instance floating ips",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "reserved-fixed-ip",
			Usage:    "delete selected instance reserved fixed ips",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client2.NewAIClusterClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := ai.DeleteOpts{
			Volumes:          c.StringSlice("volume-id"),
			DeleteFloatings:  c.Bool("delete-floating-ips"),
			FloatingIPs:      c.StringSlice("floating-ip"),
			ReservedFixedIPs: c.StringSlice("reserved-fixed-ip"),
		}

		err = gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		results, err := ai.Delete(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		taskClient, err := task_client.NewTaskClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, taskClient, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := ai.Get(client, clusterID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete AI cluster with ID: %s", clusterID)
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

var aiClusterPowerCycleCommand = cli.Command{
	Name:      "powercycle",
	Usage:     "Stop and start AI cluster. Aka hard reboot",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "powercycle")
			return err
		}
		client, err := client2.NewAIClusterClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		aiInstances, err := ai.PowerCycleAICluster(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(aiInstances, c.String("format"))
		return nil
	},
}

var aiInstancePowerCycleCommand = cli.Command{
	Name:      "powercycle",
	Usage:     "Stop and start AI instance. Aka hard reboot",
	ArgsUsage: "<instance_id>",
	Category:  "instances",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, aiInstanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "powercycle")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instance, err := ai.PowerCycleAIInstance(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var aiClusterRebootCommand = cli.Command{
	Name:      "reboot",
	Usage:     "Reboot AI cluster instaces",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "reboot")
			return err
		}
		client, err := client2.NewAIClusterClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instances, err := ai.RebootAICluster(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instances, c.String("format"))
		return nil
	},
}

var aiInstanceRebootCommand = cli.Command{
	Name:      "reboot",
	Usage:     "Reboot AI instance",
	ArgsUsage: "<instance_id>",
	Category:  "instances",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, aiInstanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "reboot")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instance, err := ai.RebootAIInstance(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var aiInstanceGetConsoleCommand = cli.Command{
	Name:      "console",
	Usage:     "Get AI instance console",
	ArgsUsage: "<instance_id>",
	Category:  "instances",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, aiInstanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "console")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		console, err := ai.GetInstanceConsole(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(console, c.String("format"))
		return nil
	},
}

var aiClusterSuspendCommand = cli.Command{
	Name:      "suspend",
	Usage:     "Suspend AI cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		cluserID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "suspend")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instance, err := ai.Suspend(client, cluserID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var aiClusterResumeCommand = cli.Command{
	Name:      "resume",
	Usage:     "Resume AI cluser",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, aiClusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "resume")
			return err
		}
		client, err := client.NewAIClusterClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instances, err := ai.Resume(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instances, c.String("format"))
		return nil
	},
}

var aiClusterAvailableImagesCommand = cli.Command{
	Name:     "list",
	Usage:    "List images available for AI cluser",
	Category: "image",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "private",
			Aliases:  []string{"p"},
			Usage:    "only private images. any value to show private images",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "visibility",
			Aliases:  []string{"v"},
			Usage:    fmt.Sprintf("image visibility type. output in %s", strings.Join(visibilityTypes, ", ")),
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "gpu",
			Usage:    "only gpu images or not gpu images",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		imageClient, err := client.NewAIImageClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		if c.Bool("gpu") {
			imageClient, err = client.NewAIGPUImageClientV1(c)
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
		}
		opts := aiimages.AIImageListOpts{
			Visibility: c.String("visibility"),
			Private:    c.String("private"),
		}
		images, err := aiimages.ListAll(imageClient, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(images, c.String("format"))
		return nil
	},
}

var aiClusterAvailableFlavorsCommand = cli.Command{
	Name:     "list",
	Usage:    "List flavors available for AI cluser",
	Category: "flavor",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "disabled",
			Usage:    "show disabled flavors",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "capacity",
			Usage:    "show flavor capacity",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "price",
			Usage:    "show flavor price",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewAIFlavorClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := aiflavors.AIFlavorListOpts{
			Disabled:        c.Bool("disabled"),
			IncludeCapacity: c.Bool("capacity"),
			IncludePrices:   c.Bool("price"),
		}
		flavors, err := aiflavors.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(flavors, c.String("format"))
		return nil
	},
}
