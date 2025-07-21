package file_shares

import (
	"fmt"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/file_shares/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	cmeta "github.com/G-Core/gcorelabscloud-go/client/utils/metadata"
	"github.com/G-Core/gcorelabscloud-go/gcore/file_share/v1/file_shares"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"

	"github.com/urfave/cli/v2"
)

var fileShareIDText = "share_id is mandatory argument"

var accessRuleIDText = "rule_id is mandatory argument"

var fileShareCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create file share",
	Category: "file share",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "File share name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "volume-type",
			Usage:    "File share volume type (default_share_type or vast_share_type)",
			Value:    "default_share_type",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "protocol",
			Aliases:  []string{"p"},
			Usage:    "File share protocol",
			Value:    "NFS",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "size",
			Usage:    "File share size GB",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "network",
			Usage:    "File share network id (required for default_share_type)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "subnet",
			Usage:    "File share subnet id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "acl-source-address",
			Usage:    "File share source ip address or subnet",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "acl-access-mode",
			Usage:    "File share access mode (ro/rw)",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "metadata",
			Usage:    "File share tags (deprecated, use --tags). Example: --metadata one=two --metadata three=four",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "tags",
			Usage:    "File share tags. Example: --tags one=two --tags three=four",
			Required: false,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := client.NewFileShareClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}

		var tags map[string]string
		if len(c.StringSlice("tags")) > 0 {
			tags, err = utils.StringSliceToTags(c.StringSlice("tags"))
			if err != nil {
				return cli.Exit(err, 1)
			}
		} else if len(c.StringSlice("metadata")) > 0 {
			tags, err = utils.StringSliceToTags(c.StringSlice("metadata"))
			if err != nil {
				return cli.Exit(err, 1)
			}
		}

		// Validate volume-type
		volumeType := c.String("volume-type")
		if volumeType != "default_share_type" && volumeType != "vast_share_type" {
			return cli.Exit("--volume-type must be either 'default_share_type' or 'vast_share_type'", 1)
		}

		// Validate if user provided network and subnet for vast_share_type, which are automatically set
		// for vast_share_type, so they should not be provided by user.
		if volumeType == "vast_share_type" && (c.String("network") != "" || c.String("subnet") != "") {
			return cli.Exit("--network and/or --subnet should not be provided for vast_share_type", 1)
		}

		opts := file_shares.CreateOpts{
			Name:       c.String("name"),
			VolumeType: c.String("volume-type"),
			Protocol:   c.String("protocol"),
			Size:       c.Int("size"),
			Access:     getAccessRules(c),
			Tags:       tags,
		}

		// Validate if user provided network and subnet for default_share_type, which are required.
		if opts.VolumeType == "default_share_type" {
			if c.String("network") == "" {
				return cli.Exit("--network is required for volume-type=default_share_type (default)", 1)
			}
			opts.Network = &file_shares.FileShareNetworkOpts{
				NetworkID: c.String("network"),
				SubnetID:  c.String("subnet"),
			}
		}

		results, err := file_shares.Create(client, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		if results == nil {
			return cli.Exit(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			fileShareID, err := file_shares.ExtractFileShareIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve file share ID from task info: %w", err)
			}
			fileShare, err := file_shares.Get(client, fileShareID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get file share with ID: %s. Error: %w", fileShareID, err)
			}
			utils.ShowResults(fileShare, c.String("format"))
			return nil, nil
		})
	},
}

func getAccessRules(c *cli.Context) []file_shares.CreateAccessRuleOpts {
	aclIpAdresses := c.StringSlice("acl-source-address")
	aclAccessModes := c.StringSlice("acl-access-mode")
	res := make([]file_shares.CreateAccessRuleOpts, 0, len(aclIpAdresses))
	for index := range aclIpAdresses {

		opts := file_shares.CreateAccessRuleOpts{
			IPAddress:  utils.StringFromIndex(aclIpAdresses, index, ""),
			AccessMode: utils.StringFromIndex(aclAccessModes, index, ""),
		}
		res = append(res, opts)
	}
	return res

}

var fileShareGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get file share information",
	ArgsUsage: "<share_id>",
	Category:  "file share",
	Action: func(c *cli.Context) error {
		fileShareID, err := flags.GetFirstStringArg(c, fileShareIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewFileShareClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		fileShare, err := file_shares.Get(client, fileShareID).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		if fileShare == nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(fileShare, c.String("format"))
		return nil
	},
}

var fileShareUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update file share",
	ArgsUsage: "<share_id>",
	Category:  "file share",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "File share name",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "tags",
			Aliases:  []string{"t"},
			Usage:    "File share key-value tags. Example: --tags key1=value1 --tags key2=value2",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "remove-tags",
			Aliases:  []string{"rt"},
			Usage:    "File share tag names. Example: --remove-tags key1 --remove-tags key2",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		fileShareID, err := flags.GetFirstStringArg(c, fileShareIDText)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		client, err := client.NewFileShareClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}

		tagsToAddOrReplace, err := utils.StringSliceToTags(c.StringSlice("tags"))
		tagsToRemove := c.StringSlice("remove-tags")
		if c.String("name") == "" && len(tagsToAddOrReplace) == 0 && len(tagsToRemove) == 0 {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.Exit("At least one of the flags --name, --tags or --remove-tags must be provided", 1)
		}

		opts := file_shares.UpdateWithTagsOpts{}
		if c.String("name") != "" {
			opts.Name = c.String("name")
		}
		tags := map[string]*string{}
		for tagKey, tagValue := range tagsToAddOrReplace {
			tags[tagKey] = utils.StringToPointer(tagValue)
		}
		for _, tagKey := range tagsToRemove {
			tags[tagKey] = nil // nil value indicates removal of the tag
		}
		opts.Tags = tags
		fileShare, err := file_shares.UpdateWithTags(client, fileShareID, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		if fileShare == nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(fileShare, c.String("format"))
		return nil
	},
}

var fileShareResizeCommand = cli.Command{
	Name:      "resize",
	Usage:     "Resize file share",
	ArgsUsage: "<share_id>",
	Category:  "file share",
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:     "size",
			Aliases:  []string{"s"},
			Usage:    "File share size",
			Required: true,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		fileShareID, err := flags.GetFirstStringArg(c, fileShareIDText)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		client, err := client.NewFileShareClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		opts := file_shares.ExtendOpts{
			Size: c.Int("size"),
		}
		results, err := file_shares.Extend(client, fileShareID, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		if results == nil {
			return cli.Exit(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			fileShare, err := file_shares.Get(client, fileShareID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get file share with ID: %s. Error: %w", fileShareID, err)
			}
			return fileShare, nil
		})
	},
}

var fileShareDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete file share",
	ArgsUsage: "<share_id>",
	Category:  "file share",
	Action: func(c *cli.Context) error {
		fileShareID, err := flags.GetFirstStringArg(c, fileShareIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewFileShareClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		results, err := file_shares.Delete(client, fileShareID).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		if results == nil {
			return cli.Exit(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := file_shares.Get(client, fileShareID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete file share with ID: %s", fileShareID)
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

var fileShareListCommand = cli.Command{
	Name:     "list",
	Usage:    "List file shares",
	Category: "file share",
	Action: func(c *cli.Context) error {
		client, err := client.NewFileShareClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		pages, err := file_shares.List(client).AllPages()
		if err != nil {
			return cli.Exit(err, 1)
		}
		results, err := file_shares.ExtractFileShares(pages)
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var fileShareAccessRuleListCommand = cli.Command{
	Name:      "list",
	Usage:     "List file share access rules",
	ArgsUsage: "<share_id>",
	Category:  "file share access rule",
	Action: func(c *cli.Context) error {
		fileShareID, err := flags.GetFirstStringArg(c, fileShareIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}
		client, err := client.NewFileShareClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		pages, err := file_shares.ListAccessRules(client, fileShareID).AllPages()
		if err != nil {
			return cli.Exit(err, 1)
		}
		results, err := file_shares.ExtractAccessRule(pages)
		if err != nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var fileShareAccessRuleCreateCommand = cli.Command{
	Name:      "create",
	Usage:     "Create file share access rules",
	ArgsUsage: "<share_id>",
	Category:  "file share access rule",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "acl-source-address",
			Usage:    "File share source ip address or subnet",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "acl-access-mode",
			Usage:    "File share access mode (ro/rw)",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		fileShareID, err := flags.GetFirstStringArg(c, fileShareIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return err
		}
		client, err := client.NewFileShareClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		opts := file_shares.CreateAccessRuleOpts{
			IPAddress:  c.String("acl-source-address"),
			AccessMode: c.String("acl-access-mode"),
		}
		accessRule, err := file_shares.CreateAccessRule(client, fileShareID, opts).Extract()
		if err != nil {
			return cli.Exit(err, 1)
		}
		if accessRule == nil {
			return cli.Exit(err, 1)
		}
		utils.ShowResults(accessRule, c.String("format"))
		return nil
	},
}

func GetSecondStringArg(c *cli.Context, errorText string) (string, error) {
	arg := c.Args().Get(1)
	if arg == "" {
		return "", cli.Exit(fmt.Errorf(errorText), 1)
	}
	return arg, nil
}

var fileShareAccessRuleDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete file share access rules",
	ArgsUsage: "<share_id> <rule_id>",
	Category:  "file share access rule",
	Action: func(c *cli.Context) error {
		fileShareID, err := flags.GetFirstStringArg(c, fileShareIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		accessRuleID, err := GetSecondStringArg(c, accessRuleIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewFileShareClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.Exit(err, 1)
		}
		err = file_shares.DeleteAccessRule(client, fileShareID, accessRuleID).ExtractErr()
		if err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	},
}

var Commands = cli.Command{
	Name:  "share",
	Usage: "GCloud file share API",
	Subcommands: []*cli.Command{
		&fileShareListCommand,
		&fileShareGetCommand,
		&fileShareDeleteCommand,
		&fileShareCreateCommand,
		&fileShareUpdateCommand,
		&fileShareResizeCommand,
		{
			Name:     "metadata",
			Usage:    "File share metadata",
			Category: "file share metadata",
			Subcommands: []*cli.Command{
				cmeta.NewMetadataListCommand(
					client.NewFileShareClientV1,
					"Get file share metadata",
					"<share_id>",
					"share_id is mandatory argument",
				),
				cmeta.NewMetadataGetCommand(
					client.NewFileShareClientV1,
					"Show file share metadata by key",
					"<share_id>",
					"share_id is mandatory argument",
				),
				cmeta.NewMetadataDeleteCommand(
					client.NewFileShareClientV1,
					"Delete file share metadata by key",
					"<share_id>",
					"share_id is mandatory argument",
				),
				cmeta.NewMetadataCreateCommand(
					client.NewFileShareClientV1,
					"Create file share metadata. It would update existing keys",
					"<share_id>",
					"share_id is mandatory argument",
				),
				cmeta.NewMetadataUpdateCommand(
					client.NewFileShareClientV1,
					"Update file share metadata. It overriding existing records",
					"<share_id>",
					"share_id is mandatory argument",
				),
				cmeta.NewMetadataReplaceCommand(
					client.NewFileShareClientV1,
					"Replace share metadata. It replace existing records",
					"<share_id>",
					"share_id is mandatory argument",
				),
			},
		},
		{
			Name:     "rule",
			Usage:    "File share access rule",
			Category: "file share access rule",
			Subcommands: []*cli.Command{
				&fileShareAccessRuleListCommand,
				&fileShareAccessRuleCreateCommand,
				&fileShareAccessRuleDeleteCommand,
			},
		},
	},
}
