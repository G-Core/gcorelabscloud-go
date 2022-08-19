package metadata

import (
	"encoding/json"
	"fmt"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	"github.com/urfave/cli/v2"
	"strings"
)

type ClientConstructor func(c *cli.Context) (*gcorecloud.ServiceClient, error)

func showResults(c *cli.Context, data interface{}) error {
	if c.Bool("test") {
		err := json.NewEncoder(c.App.Writer).Encode(data)
		if err != nil {
			return err
			//return cli.NewExitError(err, 1)
		}
	} else {
		utils.ShowResults(data, c.String("format"))
	}
	return nil
}
func StringSliceToMap(slice []string) (map[string]string, error) {
	if len(slice) == 0 {
		return nil, nil
	}
	metadata := make(map[string]string)
	for _, s := range slice {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("wrong label format: %s", s)
		}
		metadata[parts[0]] = parts[1]
	}
	return metadata, nil
}

func NewMetadataListCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     usage,
		ArgsUsage: argsUsage,
		Category:  "metadata",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "test",
				Aliases:  []string{"t"},
				Usage:    "Test flag",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			resourceID, err := flags.GetFirstStringArg(c, errorText)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "list")
				return err
			}
			client, err := cc(c)
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
			metadata, err := metadata.MetadataListAll(client, resourceID)
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			showResults(c, metadata)
			return nil
		},
	}
}

func NewMetadataGetCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return &cli.Command{
		Name:      "show",
		Usage:     usage,
		ArgsUsage: argsUsage,
		Category:  "metadata",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "metadata",
				Aliases:  []string{"m"},
				Usage:    "Metadata key",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "test",
				Aliases:  []string{"t"},
				Usage:    "Test flag",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			resourceID, err := flags.GetFirstStringArg(c, errorText)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "show")
				return err
			}
			client, err := cc(c)
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}

			metadata, err := metadata.MetadataGet(client, resourceID, c.String("metadata")).Extract()
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			showResults(c, metadata)
			return nil
		}}
}

func NewMetadataDeleteCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     usage,
		ArgsUsage: argsUsage,
		Category:  "metadata",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "metadata",
				Aliases:  []string{"m"},
				Usage:    "Metadata key",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			resourceID, err := flags.GetFirstStringArg(c, errorText)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "delete")
				return err
			}
			client, err := cc(c)
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
			err = metadata.MetadataDelete(client, resourceID, c.String("metadata")).ExtractErr()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func NewMetadataCreateCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     usage,
		ArgsUsage: argsUsage,
		Category:  "metadata",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "metadata",
				Aliases:  []string{"m"},
				Usage:    "Example: --metadata one=two --metadata three=four",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			resourceID, err := flags.GetFirstStringArg(c, errorText)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "create")
				return err
			}
			client, err := cc(c)
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
			opts, err := StringSliceToMap(c.StringSlice("metadata"))
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
			err = metadata.MetadataCreateOrUpdate(client, resourceID, opts).ExtractErr()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func NewMetadataUpdateCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return &cli.Command{
		Name:      "update",
		Usage:     usage,
		ArgsUsage: argsUsage,
		Category:  "metadata",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "metadata",
				Aliases:  []string{"m"},
				Usage:    "Example: --metadata one=two --metadata three=four",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			resourceID, err := flags.GetFirstStringArg(c, errorText)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "update")
				return err
			}
			client, err := cc(c)
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
			opts, err := StringSliceToMap(c.StringSlice("metadata"))
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
			err = metadata.MetadataCreateOrUpdate(client, resourceID, opts).ExtractErr()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		}}
}

func NewMetadataReplaceCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return &cli.Command{
		Name:      "replace",
		Usage:     usage,
		ArgsUsage: argsUsage,
		Category:  "metadata",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "metadata",
				Aliases:  []string{"m"},
				Usage:    "Example: --metadata one=two --metadata three=four",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			resourceID, err := flags.GetFirstStringArg(c, errorText)
			if err != nil {
				_ = cli.ShowCommandHelp(c, "replace")
				return err
			}
			client, err := cc(c)
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
			opts, err := StringSliceToMap(c.StringSlice("metadata"))
			if err != nil {
				_ = cli.ShowAppHelp(c)
				return cli.NewExitError(err, 1)
			}
			err = metadata.MetadataReplace(client, resourceID, opts).ExtractErr()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		}}
}
