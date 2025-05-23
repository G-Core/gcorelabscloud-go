package images

import (
	"strings"

	"github.com/urfave/cli/v2"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/client"
	taskclient "github.com/G-Core/gcorelabscloud-go/client/tasks/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/gpu/v3/images"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
)

func stringToMap(slice []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, s := range slice {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, cli.Exit("invalid metadata format", 1)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}

var imageUploadFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "url",
		Usage:    "Image URL",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "name",
		Usage:    "Image name",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "ssh-key",
		Usage:    "Permission to use SSH key in instances (allow/deny)",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "cow-format",
		Usage:    "When True, image cannot be deleted unless all volumes created from it are deleted",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "architecture",
		Usage:    "Image architecture type (aarch64/x86_64)",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "os-distro",
		Usage:    "OS Distribution (Debian/CentOS/Ubuntu/etc)",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "os-type",
		Usage:    "Operating system type (linux/windows)",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "os-version",
		Usage:    "OS version",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "hw-firmware-type",
		Usage:    "Type of firmware for booting the guest (bios/uefi)",
		Required: false,
	},
	&cli.StringSliceFlag{
		Name:     "metadata",
		Usage:    "Image metadata (key=value)",
		Required: false,
	},
}

var baremetalCommand = cli.Command{
	Name:        "baremetal",
	Usage:       "Manage baremetal GPU resources",
	Description: "Commands for managing baremetal GPU resources",
	Subcommands: []*cli.Command{
		{
			Name:        "images",
			Usage:       "Manage baremetal GPU images",
			Description: "Commands for managing baremetal GPU images",
			Subcommands: []*cli.Command{
				{
					Name:        "upload",
					Usage:       "Upload baremetal GPU image",
					Description: "Upload a new baremetal GPU image with the specified URL and name",
					Category:    "images",
					ArgsUsage:   " ",
					Flags:       append(imageUploadFlags, flags.WaitCommandFlags...),
					Action:      uploadBaremetalImageAction,
				},
				{
					Name:        "list",
					Usage:       "List baremetal GPU images",
					Description: "List all baremetal GPU images",
					Category:    "images",
					ArgsUsage:   " ",
					Action:      listBaremetalImagesAction,
				},
				{
					Name:        "show",
					Usage:       "Show baremetal GPU image details",
					Description: "Show details of a specific baremetal GPU image",
					Category:    "images",
					ArgsUsage:   "<image_id>",
					Action:      showBaremetalImageAction,
				},
				{
					Name:        "delete",
					Usage:       "Delete baremetal GPU image",
					Description: "Delete baremetal GPU image by ID",
					Category:    "images",
					ArgsUsage:   "<image_id>",
					Flags:       flags.WaitCommandFlags,
					Action:      deleteBaremetalImageAction,
				},
			},
		},
	},
}

var virtualCommand = cli.Command{
	Name:        "virtual",
	Usage:       "Manage virtual GPU resources",
	Description: "Commands for managing virtual GPU resources",
	Subcommands: []*cli.Command{
		{
			Name:        "images",
			Usage:       "Manage virtual GPU images",
			Description: "Commands for managing virtual GPU images",
			Subcommands: []*cli.Command{
				{
					Name:        "upload",
					Usage:       "Upload virtual GPU image",
					Description: "Upload a new virtual GPU image with the specified URL and name",
					Category:    "images",
					ArgsUsage:   " ",
					Flags:       append(imageUploadFlags, flags.WaitCommandFlags...),
					Action:      uploadVirtualImageAction,
				},
				{
					Name:        "list",
					Usage:       "List virtual GPU images",
					Description: "List all virtual GPU images",
					Category:    "images",
					ArgsUsage:   " ",
					Action:      listVirtualImagesAction,
				},
				{
					Name:        "show",
					Usage:       "Show virtual GPU image details",
					Description: "Show details of a specific virtual GPU image",
					Category:    "images",
					ArgsUsage:   "<image_id>",
					Action:      showVirtualImageAction,
				},
				{
					Name:        "delete",
					Usage:       "Delete virtual GPU image",
					Description: "Delete virtual GPU image by ID",
					Category:    "images",
					ArgsUsage:   "<image_id>",
					Flags:       flags.WaitCommandFlags,
					Action:      deleteVirtualImageAction,
				},
			},
		},
	},
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// uploadImageAction handles the common logic for both virtual and baremetal image uploads
func uploadImageAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	if c.Args().Len() > 0 {
		return cli.ShowCommandHelp(c, "")
	}

	// Only validate if not showing help
	if !c.Bool("help") && (c.String("url") == "" || c.String("name") == "") {
		_ = cli.ShowCommandHelp(c, "")
		return cli.Exit("Required flags 'url' and 'name' must be set", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	sshKey := images.SshKeyType(c.String("ssh-key"))
	cowFormat := c.Bool("cow-format")
	osType := images.ImageOsType(c.String("os-type"))
	hwType := images.ImageHwFirmwareType(c.String("hw-firmware-type"))

	opts := images.ImageOpts{
		URL:            c.String("url"),
		Name:           c.String("name"),
		SshKey:         &sshKey,
		CowFormat:      &cowFormat,
		Architecture:   stringPtr(c.String("architecture")),
		OsDistro:       stringPtr(c.String("os-distro")),
		OsType:         &osType,
		OsVersion:      stringPtr(c.String("os-version")),
		HwFirmwareType: &hwType,
	}

	if c.IsSet("metadata") {
		metadata, err := stringToMap(c.StringSlice("metadata"))
		if err != nil {
			return cli.Exit(err, 1)
		}
		metadataInterface := make(map[string]interface{})
		for k, v := range metadata {
			metadataInterface[k] = v
		}
		opts.Metadata = metadataInterface
	}

	results := images.UploadImage(gpuClient, opts)
	if results.Err != nil {
		return cli.Exit(results.Err, 1)
	}

	taskID, err := results.Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, taskID, true, func(task tasks.TaskID) (interface{}, error) {
		return task, nil
	})
}

func uploadBaremetalImageAction(c *cli.Context) error {
	return uploadImageAction(c, client.NewGPUBaremetalClientV3)
}

func uploadVirtualImageAction(c *cli.Context) error {
	return uploadImageAction(c, client.NewGPUVirtualClientV3)
}

// listImagesAction handles the common logic for listing both virtual and baremetal images
func listImagesAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	pages, err := images.List(gpuClient).AllPages()
	if err != nil {
		return cli.Exit(err, 1)
	}

	imageList, err := images.ExtractImages(pages)
	if err != nil {
		return cli.Exit(err, 1)
	}

	utils.ShowResults(imageList, c.String("format"))
	return nil
}

// deleteImageAction handles the common logic for deleting both virtual and baremetal images
func deleteImageAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	imageID := c.Args().First()
	if imageID == "" {
		_ = cli.ShowCommandHelp(c, "delete")
		return cli.Exit("image ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	results := images.Delete(gpuClient, imageID)
	if results.Err != nil {
		return cli.Exit(results.Err, 1)
	}

	taskID, err := results.Extract()
	if err != nil {
		return cli.Exit(err, 1)
	}

	taskClient, err := taskclient.NewTaskClientV1(c)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return utils.WaitTaskAndShowResult(c, taskClient, taskID, true, func(task tasks.TaskID) (interface{}, error) {
		return task, nil
	})
}

func listBaremetalImagesAction(c *cli.Context) error {
	return listImagesAction(c, client.NewGPUBaremetalClientV3)
}

func listVirtualImagesAction(c *cli.Context) error {
	return listImagesAction(c, client.NewGPUVirtualClientV3)
}

func deleteBaremetalImageAction(c *cli.Context) error {
	return deleteImageAction(c, client.NewGPUBaremetalClientV3)
}

func deleteVirtualImageAction(c *cli.Context) error {
	return deleteImageAction(c, client.NewGPUVirtualClientV3)
}

func showImageAction(c *cli.Context, newClient func(*cli.Context) (*gcorecloud.ServiceClient, error)) error {
	imageID := c.Args().First()
	if imageID == "" {
		_ = cli.ShowCommandHelp(c, "show")
		return cli.Exit("image ID is required", 1)
	}

	gpuClient, err := newClient(c)
	if err != nil {
		_ = cli.ShowAppHelp(c)
		return cli.Exit(err, 1)
	}

	imageDetails := images.Get(gpuClient, imageID)
	if imageDetails.Err != nil {
		return cli.Exit(imageDetails.Err, 1)
	}

	utils.ShowResults(imageDetails.Body, c.String("format"))
	return nil
}

func showVirtualImageAction(c *cli.Context) error {
	return showImageAction(c, client.NewGPUVirtualClientV3)
}

func showBaremetalImageAction(c *cli.Context) error {
	return showImageAction(c, client.NewGPUBaremetalClientV3)
}

// BaremetalCommands returns commands for managing baremetal GPU images
func BaremetalCommands() *cli.Command {
	return &cli.Command{
		Name:        "images",
		Usage:       "Manage baremetal GPU images",
		Description: "Commands for managing baremetal GPU images",
		Subcommands: []*cli.Command{
			{
				Name:        "upload",
				Usage:       "Upload baremetal GPU image",
				Description: "Upload a new baremetal GPU image with the specified URL and name",
				Category:    "images",
				ArgsUsage:   " ",
				Flags:       append(imageUploadFlags, flags.WaitCommandFlags...),
				Action:      uploadBaremetalImageAction,
			},
			{
				Name:        "list",
				Usage:       "List baremetal GPU images",
				Description: "List all baremetal GPU images",
				Category:    "images",
				ArgsUsage:   " ",
				Action:      listBaremetalImagesAction,
			},
			{
				Name:        "show",
				Usage:       "Show baremetal GPU image details",
				Description: "Show details of a specific baremetal GPU image",
				Category:    "images",
				ArgsUsage:   "<image_id>",
				Action:      showBaremetalImageAction,
			},
			{
				Name:        "delete",
				Usage:       "Delete baremetal GPU image",
				Description: "Delete baremetal GPU image by ID",
				Category:    "images",
				ArgsUsage:   "<image_id>",
				Flags:       flags.WaitCommandFlags,
				Action:      deleteBaremetalImageAction,
			},
		},
	}
}

// VirtualCommands returns commands for managing virtual GPU images
func VirtualCommands() *cli.Command {
	return &cli.Command{
		Name:        "images",
		Usage:       "Manage virtual GPU images",
		Description: "Commands for managing virtual GPU images",
		Subcommands: []*cli.Command{
			{
				Name:        "upload",
				Usage:       "Upload virtual GPU image",
				Description: "Upload a new virtual GPU image with the specified URL and name",
				Category:    "images",
				ArgsUsage:   " ",
				Flags:       append(imageUploadFlags, flags.WaitCommandFlags...),
				Action:      uploadVirtualImageAction,
			},
			{
				Name:        "list",
				Usage:       "List virtual GPU images",
				Description: "List all virtual GPU images",
				Category:    "images",
				ArgsUsage:   " ",
				Action:      listVirtualImagesAction,
			},
			{
				Name:        "show",
				Usage:       "Show virtual GPU image details",
				Description: "Show details of a specific virtual GPU image",
				Category:    "images",
				ArgsUsage:   "<image_id>",
				Action:      showVirtualImageAction,
			},
			{
				Name:        "delete",
				Usage:       "Delete virtual GPU image",
				Description: "Delete virtual GPU image by ID",
				Category:    "images",
				ArgsUsage:   "<image_id>",
				Flags:       flags.WaitCommandFlags,
				Action:      deleteVirtualImageAction,
			},
		},
	}
}
