package keypairs

import (
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/keypairs/v2/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/keypair/v2/keypairs"

	"github.com/urfave/cli/v2"
)

var keyPairIDText = "keypair_id is mandatory argument"

var keypairListCommand = cli.Command{
	Name:     "list",
	Usage:    "List keypairs",
	Category: "keypair2",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "user-id",
			Aliases:  []string{"u"},
			Usage:    "User ID",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "project-id",
			Usage:    "Project ID",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewKeypairClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := keypairs.ListOpts{
			UserID:    c.String("user-id"),
			ProjectID: c.Int("project-id"),
		}
		results, err := keypairs.ListAll(client, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var keypairGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get keypair information",
	ArgsUsage: "<keypair_id>",
	Category:  "keypair2",
	Action: func(c *cli.Context) error {
		keypairID, err := flags.GetFirstStringArg(c, keyPairIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewKeypairClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := keypairs.Get(client, keypairID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(task, c.String("format"))
		return nil
	},
}

var keypairDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete keypair by ID",
	ArgsUsage: "<keypair_id>",
	Category:  "keypair2",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		keypairID, err := flags.GetFirstStringArg(c, keyPairIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewKeypairClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = keypairs.Delete(client, keypairID).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var keypairCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create keypair",
	Category: "keypair2",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Keypair name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "ssh-public-key",
			Usage:    "Keypair SSH public key file",
			Aliases:  []string{"k"},
			Required: false,
		},
		&cli.IntFlag{
			Name:     "project-id",
			Usage:    "Project ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := client.NewKeypairClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		sshPublicKeyFile := c.String("ssh-public-key")
		var sshKeyContent string
		if sshPublicKeyFile != "" {
			data, err := utils.ReadFile(sshPublicKeyFile)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "create")
				return cli.NewExitError(err, 1)
			}
			sshKeyContent = string(data)
		}
		opts := keypairs.CreateOpts{
			Name:      c.String("name"),
			PublicKey: sshKeyContent,
			ProjectID: c.Int("project-id"),
		}
		result, err := keypairs.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var KeypairCommands = cli.Command{
	Name:  "keypair2",
	Usage: "GCloud keypairs V2 API",
	Subcommands: []*cli.Command{
		&keypairListCommand,
		&keypairGetCommand,
		&keypairDeleteCommand,
		&keypairCreateCommand,
	},
}
