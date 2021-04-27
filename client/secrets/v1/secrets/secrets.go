package secrets

import (
	"fmt"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/secrets/v1/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/secret/v1/secrets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/urfave/cli/v2"
)

var secretsIDText = "secrets_id is mandatory argument"

var SecretCommands = cli.Command{
	Name:  "secrets",
	Usage: "GCloud secrets v1 API",
	Subcommands: []*cli.Command{
		&secretListCommand,
		&secretGetCommand,
		&secretDeleteCommand,
		&secretCreateCommand,
	},
}

var secretListCommand = cli.Command{
	Name:  "list",
	Usage: "List secret",
	Action: func(c *cli.Context) error {
		client, err := client.NewSecretClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := secrets.ListAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var secretGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get secret information",
	ArgsUsage: "<secret_id>",
	Action: func(c *cli.Context) error {
		secretID, err := flags.GetFirstStringArg(c, secretsIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := client.NewSecretClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		secret, err := secrets.Get(client, secretID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(secret, c.String("format"))
		return nil
	},
}

var secretDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete secret by ID",
	ArgsUsage: "<secret_id>",
	Action: func(c *cli.Context) error {
		secretID, err := flags.GetFirstStringArg(c, secretsIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := client.NewSecretClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := secrets.Delete(client, secretID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		if result != nil {
			utils.ShowResults(result, c.String("format"))
		}
		return nil
	},
}

var secretCreateCommand = cli.Command{
	Name:  "create",
	Usage: "Create secret",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Secret name",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage: `Secret type.
			symmetric - Used for storing byte arrays such as keys suitable for symmetric encryption;
			public - Used for storing the public key of an asymmetric keypair;
			private - Used for storing the private key of an asymmetric keypair;
			passphrase - Used for storing plain text passphrases;
			certificate - Used for storing cryptographic certificates such as X.509 certificates;
			opaque - Used for backwards compatibility with previous versions of the API`,
			Required: true,
		},
		&cli.StringFlag{
			Name:     "payload",
			Aliases:  []string{"p"},
			Usage:    "Secret payload",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "payload-content-encoding",
			Aliases:  []string{"pce"},
			Usage:    "The encoding used for the payload to be able to include it in the JSON request. Currently only base64 is supported.",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "payload-content-type",
			Aliases:  []string{"pct"},
			Usage:    "The media type for the content of the payload.",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "algorithm",
			Aliases: []string{"a"},
			Usage:   "Metadata provided by a user or system for informational purposes. ",
		},
		&cli.StringFlag{
			Name:    "mode",
			Aliases: []string{"m"},
			Usage:   "Metadata provided by a user or system for informational purposes.",
		},
		&cli.IntFlag{
			Name:    "bit-len",
			Aliases: []string{"bl"},
			Usage:   "Metadata provided by a user or system for informational purposes. Value must be greater than zero.",
		},
		&cli.StringFlag{
			Name:    "expiration-time",
			Aliases: []string{"e"},
			Usage:   fmt.Sprintf("Date when the secret will expire. Format `%s`", gcorecloud.RFC3339NoZ),
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := client.NewSecretClientV1(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := secrets.CreateOpts{
			Name:                   c.String("name"),
			Type:                   secrets.SecretType(c.String("type")),
			Payload:                c.String("payload"),
			PayloadContentEncoding: c.String("payload-content-encoding"),
			PayloadContentType:     c.String("payload-content-type"),
		}

		algorithm := c.String("algorithm")
		if len(algorithm) != 0 {
			opts.Algorithm = &algorithm
		}
		bitLen := c.Int("bit-len")
		if bitLen > 0 {
			opts.BitLength = &bitLen
		}
		rawExpTime := c.String("expiration-time")
		if len(rawExpTime) != 0 {
			expirationTime, err := time.Parse(gcorecloud.RFC3339NoZ, rawExpTime)
			if err != nil {
				_ = cli.ShowSubcommandHelp(c)
				return cli.NewExitError(err, 1)
			}
			opts.Expiration = &expirationTime
		}
		mode := c.String("mode")
		if len(mode) != 0 {
			opts.Mode = &mode
		}

		results, err := secrets.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			secretID, err := secrets.ExtractSecretIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve secret ID from task info: %w", err)
			}
			secret, err := secrets.Get(client, secretID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get secret with ID: %s. Error: %w", secretID, err)
			}
			utils.ShowResults(secret, c.String("format"))
			return nil, nil
		})
	},
}
