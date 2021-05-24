package instances

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/client/instances/v1/client"
	client2 "github.com/G-Core/gcorelabscloud-go/client/instances/v2/client"
	"github.com/G-Core/gcorelabscloud-go/gcore/baremetal/v1/bminstances"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/urfave/cli/v2"
)

var (
	instanceIDText            = "instance_id is mandatory argument"
	volumeSourceType          = types.VolumeSource("").StringList()
	volumeType                = volumes.VolumeType("").StringList()
	interfaceTypes            = types.InterfaceType("").StringList()
	interfaceFloatingIPSource = types.FloatingIPSource("").StringList()
	bootableIndex             = 0
)

var Commands = cli.Command{
	Name:  "instance",
	Usage: "GCloud instances API",
	Subcommands: []*cli.Command{
		&instanceGetCommand,
		&instanceListCommand,
		&instanceCreateCommandV2,
		&instanceDeleteCommand,
		&instanceStartCommand,
		&instanceStopCommand,
		&instancePowerCycleCommand,
		&instanceRebootCommand,
		&instanceSuspendCommand,
		&instanceResumeCommand,
		&instanceResizeCommand,
		&instanceCreateBaremetalCommand,
		{
			Name:  "interface",
			Usage: "Instance interfaces",
			Subcommands: []*cli.Command{
				&instanceListInterfacesCommand,
			},
		},
		{
			Name:  "securitygroup",
			Usage: "Instance security groups",
			Subcommands: []*cli.Command{
				&instanceListSecurityGroupsCommand,
				&instanceAssignSecurityGroupsCommand,
				&instanceUnAssignSecurityGroupsCommand,
			},
		},
		{
			Name:  "metadata",
			Usage: "Instance metadata",
			Subcommands: []*cli.Command{
				&metadataListCommand,
				&metadataGetCommand,
				&metadataCreateCommand,
				&metadataUpdateCommand,
				&metadataDeleteCommand,
			},
		},
	},
}

func StringSliceToMetadataSetOpts(slice []string) (*instances.MetadataSetOpts, error) {
	if len(slice) == 0 {
		return nil, nil
	}
	var m instances.MetadataSetOpts
	for _, s := range slice {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return &m, fmt.Errorf("wrong label format: %s", s)
		}
		m.Metadata = append(m.Metadata, instances.MetadataOpts{Key: parts[0], Value: parts[1]})
	}
	return &m, nil
}

func StringSliceToAppConfigSetOpts(slice []string) (map[string]interface{}, error) {
	if len(slice) == 0 {
		return nil, nil
	}
	m := make(map[string]interface{}, len(slice))
	for _, s := range slice {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return m, fmt.Errorf("wrong label format: %s", s)
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}

func getUserData(c *cli.Context) (string, error) {
	userData := ""
	userDataFile := c.String("user-data-file")
	userDataContent := c.String("user-data")

	if userDataFile != "" {
		fileContent, err := ioutil.ReadFile(userDataFile)
		if err != nil {
			return "", err
		}
		userData = base64.StdEncoding.EncodeToString(fileContent)
	} else if userDataContent != "" {
		userData = base64.StdEncoding.EncodeToString([]byte(userDataContent))
	}
	return userData, nil
}

func getInstanceVolumes(c *cli.Context) ([]instances.CreateVolumeOpts, error) {
	volumeSources := utils.GetEnumStringSliceValue(c, "volume-source")
	volumeTypes := utils.GetEnumStringSliceValue(c, "volume-type")
	volumeBootIndexes := c.IntSlice("volume-boot-index")
	volumeSizes := c.IntSlice("volume-size")
	volumeNames := c.StringSlice("volume-name")
	volumeImageIDs := c.StringSlice("volume-image-id")
	volumeVolumeIDs := c.StringSlice("volume-volume-id")
	volumeSnapshotIDs := c.StringSlice("volume-snapshot-id")

	res := make([]instances.CreateVolumeOpts, 0, len(volumeSources))

	for idx, s := range volumeSources {
		source := types.VolumeSource(s)
		bootIndex := utils.IntFromIndex(volumeBootIndexes, idx, 0)
		if !source.Bootable() {
			bootIndex = -1
		}
		opts := instances.CreateVolumeOpts{
			Source:    source,
			BootIndex: bootIndex,
			Size:      utils.IntFromIndex(volumeSizes, idx, 0),
			TypeName: func(idx int) volumes.VolumeType {
				if idx < len(volumeTypes) {
					return volumes.VolumeType(volumeTypes[idx])
				}
				return volumes.Standard
			}(idx),
			Name:       utils.StringFromIndex(volumeNames, idx, ""),
			ImageID:    utils.StringFromIndex(volumeImageIDs, idx, ""),
			SnapshotID: utils.StringFromIndex(volumeSnapshotIDs, idx, ""),
			VolumeID:   utils.StringFromIndex(volumeVolumeIDs, idx, ""),
		}
		err := gcorecloud.TranslateValidationError(opts.Validate())

		if err != nil {
			return nil, err
		}

		res = append(res, opts)

	}

	// adjust boot order number so that they goes consistently
	sort.Slice(res, func(i, j int) bool {
		return res[i].BootIndex < res[j].BootIndex
	})

	minOrder := 0
	bootableExists := false

	for _, opts := range res {
		if opts.BootIndex < 0 {
			continue
		}
		if opts.BootIndex > minOrder {
			opts.BootIndex = minOrder
		}
		minOrder++
		if opts.BootIndex == bootableIndex {
			bootableExists = true
		}
	}

	if !bootableExists {
		return res, fmt.Errorf("set bootable volume with source from next list: %s, %s", types.Image, types.Snapshot)
	}

	return res, nil

}

func getInterfaces(c *cli.Context) ([]instances.InterfaceOpts, error) {
	interfaceTypes := utils.GetEnumStringSliceValue(c, "interface-type")
	interfaceNetworkIDs := c.StringSlice("interface-network-id")
	interfaceSubnetIDs := c.StringSlice("interface-subnet-id")
	interfaceFloatingSources := utils.GetEnumStringSliceValue(c, "interface-floating-source")
	interfaceFloatingIPs := c.StringSlice("interface-floating-ip")

	res := make([]instances.InterfaceOpts, 0, len(interfaceTypes))

	for idx, t := range interfaceTypes {
		interfaceType := types.InterfaceType(t)
		var fIP *instances.CreateNewInterfaceFloatingIPOpts = nil
		if interfaceType == types.SubnetInterfaceType {
			source := types.FloatingIPSource(utils.StringFromIndex(interfaceFloatingSources, idx, ""))
			if source != "" {
				fIP = &instances.CreateNewInterfaceFloatingIPOpts{
					Source:             types.FloatingIPSource(utils.StringFromIndex(interfaceFloatingSources, idx, "")),
					ExistingFloatingID: utils.StringFromIndex(interfaceFloatingIPs, idx, ""),
				}
			}
		}

		opts := instances.InterfaceOpts{
			Type:       interfaceType,
			NetworkID:  utils.StringFromIndex(interfaceNetworkIDs, idx, ""),
			SubnetID:   utils.StringFromIndex(interfaceSubnetIDs, idx, ""),
			FloatingIP: fIP,
		}

		err := gcorecloud.TranslateValidationError(opts.Validate())

		if err != nil {
			return nil, err
		}

		res = append(res, opts)

	}

	return res, nil

}

func getBaremetalInterfaces(c *cli.Context) ([]bminstances.InterfaceOpts, error) {
	interfaceTypes := utils.GetEnumStringSliceValue(c, "interface-type")
	interfaceNetworkIDs := c.StringSlice("interface-network-id")
	interfaceSubnetIDs := c.StringSlice("interface-subnet-id")
	interfaceFloatingSources := utils.GetEnumStringSliceValue(c, "interface-floating-source")
	interfaceFloatingIPs := c.StringSlice("interface-floating-ip")

	res := make([]bminstances.InterfaceOpts, 0, len(interfaceTypes))

	for idx, t := range interfaceTypes {
		interfaceType := types.InterfaceType(t)
		var fIP *bminstances.CreateNewInterfaceFloatingIPOpts = nil
		if interfaceType == types.SubnetInterfaceType {
			source := types.FloatingIPSource(utils.StringFromIndex(interfaceFloatingSources, idx, ""))
			if source != "" {
				fIP = &bminstances.CreateNewInterfaceFloatingIPOpts{
					Source:             types.FloatingIPSource(utils.StringFromIndex(interfaceFloatingSources, idx, "")),
					ExistingFloatingID: utils.StringFromIndex(interfaceFloatingIPs, idx, ""),
				}
			}
		}

		opts := bminstances.InterfaceOpts{
			Type:       interfaceType,
			NetworkID:  utils.StringFromIndex(interfaceNetworkIDs, idx, ""),
			SubnetID:   utils.StringFromIndex(interfaceSubnetIDs, idx, ""),
			FloatingIP: fIP,
		}

		err := gcorecloud.TranslateValidationError(opts.Validate())

		if err != nil {
			return nil, err
		}

		res = append(res, opts)

	}

	return res, nil

}

func getSecurityGroups(c *cli.Context) []gcorecloud.ItemID {
	securityGroups := c.StringSlice("security-group")
	res := make([]gcorecloud.ItemID, len(securityGroups))
	for _, s := range securityGroups {
		res = append(res, gcorecloud.ItemID{ID: s})
	}
	return res
}

var instanceListCommand = cli.Command{
	Name:     "list",
	Usage:    "List instances",
	Category: "instance",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "exclude-security-group",
			Aliases:  []string{"e"},
			Usage:    "exclude instances with specified security group name",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "available-floating",
			Aliases:  []string{"a"},
			Usage:    "show only instances which are able to handle floating address",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "filter out instances by name",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "flavor_id",
			Aliases:  []string{"fid"},
			Usage:    "filter out instances by flavor id",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "limit",
			Aliases:  []string{"l"},
			Usage:    "limit the number of returned instances",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "offset",
			Aliases:  []string{"o"},
			Usage:    "offset value is used to exclude the first set of records from the result",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := instances.ListOpts{
			ExcludeSecGroup:   c.String("exclude-security-group"),
			AvailableFloating: c.Bool("available-floating"),
			Name:              c.String("name"),
			FlavorID:          c.String("flavor_id"),
			Limit:             c.Int("limit"),
			Offset:            c.Int("offset"),
		}
		results, err := instances.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var instanceListInterfacesCommand = cli.Command{
	Name:      "list",
	Usage:     "List instance interfaces",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := instances.ListInterfacesAll(client, instanceID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var instanceListSecurityGroupsCommand = cli.Command{
	Name:      "list",
	Usage:     "List instance security groups",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := instances.ListSecurityGroupsAll(client, instanceID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var instanceAssignSecurityGroupsCommand = cli.Command{
	Name:      "add",
	Usage:     "Add instance security group",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "security group name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "add")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := instances.SecurityGroupOpts{Name: c.String("name")}

		err = instances.AssignSecurityGroup(client, instanceID, opts).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var instanceUnAssignSecurityGroupsCommand = cli.Command{
	Name:      "delete",
	Usage:     "Add instance security group",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "security group name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := instances.SecurityGroupOpts{Name: c.String("name")}

		err = instances.UnAssignSecurityGroup(client, instanceID, opts).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var instanceCreateCommandV2 = cli.Command{
	Name: "create",
	Usage: `
	Create instance. 
	Example: gcoreclient token instance create --flavor g1-standard-1-2 --name test1 --keypair keypair --volume-source image --volume-type standard --volume-image-id --interface-type subnet --interface-network-id 95ea2073-c5eb-448a-aed5-78045f88f24a --interface-subnet-id b7fd6e0a-36a5-4f6a-9dc4-90a39eb9833f --volume-size 2 --metadata one=two -d -w`,
	Category: "instance",
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "instance name",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "name-template",
			Usage:    "instance name templates",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "flavor",
			Usage:    "instance flavor",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "keypair",
			Aliases:  []string{"k"},
			Usage:    "instance ssh keypair",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
			Usage:    "instance password",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "username",
			Aliases:  []string{"u"},
			Usage:    "instance username",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "user-data",
			Usage:    "instance user data",
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
			Required: true,
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
			Required: true,
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
		&cli.StringFlag{
			Name:  "server-group",
			Usage: "instance server group",
		},
		&cli.StringSliceFlag{
			Name:     "metadata",
			Usage:    "instance metadata. Example: --metadata one=two --metadata three=four",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clientV2, err := client2.NewInstanceClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		clientV1, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		userData, err := getUserData(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		instanceVolumes, err := getInstanceVolumes(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		instanceInterfaces, err := getInterfaces(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		securityGroups := getSecurityGroups(c)

		metadata, err := StringSliceToMetadataSetOpts(c.StringSlice("metadata"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		opts := instances.CreateOpts{
			Flavor:         c.String("flavor"),
			Names:          c.StringSlice("name"),
			NameTemplates:  c.StringSlice("name-template"),
			Volumes:        instanceVolumes,
			Interfaces:     instanceInterfaces,
			SecurityGroups: securityGroups,
			Keypair:        c.String("keypair"),
			Password:       c.String("password"),
			Username:       c.String("username"),
			UserData:       userData,
			Metadata:       metadata,
			ServerGroupID:  c.String("server-group"),
		}

		fmt.Printf("%#v", opts)

		err = gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		results, err := instances.Create(clientV2, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, clientV1, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(clientV1, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			instanceID, err := instances.ExtractInstanceIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve instance ID from task info: %w", err)
			}
			instance, err := instances.Get(clientV1, instanceID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get instance with ID: %s. Error: %w", instanceID, err)
			}
			return instance, nil
		})
	},
}

var instanceGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get instance information",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var instanceDeleteCommand = cli.Command{
	Name:  "delete",
	Usage: "Delete instance",
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:     "volume-id",
			Usage:    "instance volume id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "floating-ip",
			Usage:    "instance floating ip",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "delete-floating-ips",
			Usage:    "delete all instance floating ips",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := instances.DeleteOpts{
			Volumes:         c.StringSlice("volume-id"),
			DeleteFloatings: c.Bool("delete-floating-ips"),
			FloatingIPs:     c.StringSlice("floating-ip"),
		}

		err = gcorecloud.TranslateValidationError(opts.Validate())
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		results, err := instances.Delete(client, instanceID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := instances.Get(client, instanceID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete instance with ID: %s", instanceID)
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

var instanceStartCommand = cli.Command{
	Name:      "start",
	Usage:     "Start instance",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "start")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instance, err := instances.Start(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var instanceStopCommand = cli.Command{
	Name:      "stop",
	Usage:     "Stop instance",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "stop")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instance, err := instances.Stop(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var instancePowerCycleCommand = cli.Command{
	Name:      "powercycle",
	Usage:     "Stop and start instance. Aka hard reboot",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "powercycle")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instance, err := instances.PowerCycle(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var instanceRebootCommand = cli.Command{
	Name:      "reboot",
	Usage:     "Reboot instance",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "reboot")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instance, err := instances.Reboot(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var instanceSuspendCommand = cli.Command{
	Name:      "suspend",
	Usage:     "Suspend instance",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "suspend")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instance, err := instances.Suspend(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var instanceResumeCommand = cli.Command{
	Name:      "resume",
	Usage:     "Resume instance",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "resume")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		instance, err := instances.Resume(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(instance, c.String("format"))
		return nil
	},
}

var instanceResizeCommand = cli.Command{
	Name:      "resize",
	Usage:     "Resize instance",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:     "flavor",
			Usage:    "instance flavor id",
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "resize")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := instances.ChangeFlavorOpts{FlavorID: c.String("flavor")}

		results, err := instances.Resize(client, instanceID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			instance, err := instances.Get(client, instanceID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get instance with ID: %s. Error: %w", instanceID, err)
			}
			return instance, nil
		})
	},
}

var metadataListCommand = cli.Command{
	Name:      "list",
	Usage:     "Get instance metadata",
	ArgsUsage: "<instance_id>",
	Category:  "metadata",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		metadata, err := instances.MetadataListAll(client, instanceID)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(metadata, c.String("format"))
		return nil
	},
}

var metadataGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Show instance metadata by key",
	ArgsUsage: "<instance_id>",
	Category:  "metadata",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "metadata",
			Aliases:  []string{"m"},
			Usage:    "instance metadata key",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		metadata, err := instances.MetadataGet(client, instanceID, c.String("metadata")).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(metadata, c.String("format"))
		return nil
	},
}

var metadataDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete instance metadata by key",
	ArgsUsage: "<instance_id>",
	Category:  "metadata",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "metadata",
			Aliases:  []string{"m"},
			Usage:    "instance metadata key",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = instances.MetadataDelete(client, instanceID, c.String("metadata")).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var metadataCreateCommand = cli.Command{
	Name:      "create",
	Usage:     "Create instance metadata. It would update existing keys",
	ArgsUsage: "<instance_id>",
	Category:  "metadata",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "metadata",
			Aliases:  []string{"m"},
			Usage:    "instance metadata. Example: --metadata one=two --metadata three=four",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts, err := StringSliceToMetadataSetOpts(c.StringSlice("metadata"))
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = instances.MetadataCreate(client, instanceID, *opts).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var metadataUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update instance metadata. It replace existing records",
	ArgsUsage: "<instance_id>",
	Category:  "metadata",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "metadata",
			Aliases:  []string{"m"},
			Usage:    "instance metadata. Example: --metadata one=two --metadata three=four",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstStringArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := client.NewInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts, err := StringSliceToMetadataSetOpts(c.StringSlice("metadata"))
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = instances.MetadataUpdate(client, instanceID, *opts).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var instanceCreateBaremetalCommand = cli.Command{
	Name: "create_baremetal",
	Usage: `
	Create baremetal instance. 
	Example: gcoreclient instance create_baremetal --flavor bm1-infrastructure-small --name test1 --keypair keypair --image-id 1ee7ccee-5003-48c9-8ae0-d96063af75b2 --interface-type external`,
	Category: "instance",
	Flags: append([]cli.Flag{
		&cli.StringSliceFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Baremetal instance name",
		},
		&cli.StringSliceFlag{
			Name:  "name-template",
			Usage: "Baremetal instance name templates",
		},
		&cli.StringFlag{
			Name:     "flavor",
			Usage:    "Baremetal instance flavor",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "keypair",
			Aliases: []string{"k"},
			Usage:   "Baremetal instance ssh keypair",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Usage:   "Baremetal instance password",
		},
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Usage:   "Baremetal instance username",
		},
		&cli.StringFlag{
			Name:  "user-data",
			Usage: "Baremetal instance user data",
		},
		&cli.StringFlag{
			Name:  "image-id",
			Usage: "Baremetal instance volume image id",
		},
		&cli.StringFlag{
			Name:  "apptemplate-id",
			Usage: "Baremetal instance volume apptemplate id",
		},
		&cli.GenericFlag{
			Name:    "interface-type",
			Aliases: []string{"it"},
			Value: &utils.EnumStringSliceValue{
				Enum: interfaceTypes,
			},
			Usage:    fmt.Sprintf("Baremetal instance interface type. output in %s", strings.Join(interfaceTypes, ", ")),
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:  "interface-network-id",
			Usage: "Baremetal instance interface network id",
		},
		&cli.StringSliceFlag{
			Name:  "interface-subnet-id",
			Usage: "Baremetal instance interface subnet id",
		},
		&cli.GenericFlag{
			Name:    "interface-floating-source",
			Aliases: []string{"ifs"},
			Value: &utils.EnumStringSliceValue{
				Enum: interfaceFloatingIPSource,
			},
			Usage: fmt.Sprintf("Baremetal instance floating ip source. output in %s", strings.Join(interfaceFloatingIPSource, ", ")),
		},
		&cli.StringSliceFlag{
			Name:  "interface-floating-ip",
			Usage: "Baremetal instance interface existing floating ip. Required when --interface-floating-source set as `existing`",
		},
		&cli.StringSliceFlag{
			Name:  "appconfig",
			Usage: "Baremetal instance appconfig. Example: --appconfig one=two --appconfig three=four",
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clientV1, err := client.NewBmInstanceClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		userData, err := getUserData(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create_baremetal")
			return cli.NewExitError(err, 1)
		}

		instanceInterfaces, err := getBaremetalInterfaces(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create_baremetal")
			return cli.NewExitError(err, 1)
		}

		appCfg, err := StringSliceToAppConfigSetOpts(c.StringSlice("appconfig"))
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create_baremetal")
			return cli.NewExitError(err, 1)
		}

		opts := bminstances.CreateOpts{
			Flavor:        c.String("flavor"),
			Names:         c.StringSlice("name"),
			NameTemplates: c.StringSlice("name-template"),
			ImageID:       c.String("image-id"),
			AppTemplateID: c.String("apptemplate-id"),
			Interfaces:    instanceInterfaces,
			Keypair:       c.String("keypair"),
			Password:      c.String("password"),
			Username:      c.String("username"),
			UserData:      userData,
			AppConfig:     appCfg,
		}

		results, err := bminstances.Create(clientV1, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, clientV1, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(clientV1, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			instanceID, err := instances.ExtractInstanceIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve instance ID from task info: %w", err)
			}
			instance, err := instances.Get(clientV1, instanceID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get instance with ID: %s. Error: %w", instanceID, err)
			}
			return instance, nil
		})
	},
}
