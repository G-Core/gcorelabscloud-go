package metadata

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata"
	v1 "github.com/G-Core/gcorelabscloud-go/gcore/utils/metadata/v1/metadata"
)

type ClientConstructor func(c *cli.Context) (*gcorecloud.ServiceClient, error)

// Operations are the SDK calls that the metadata commands make
type Operations struct {
	ListAll        func(client *gcorecloud.ServiceClient, id string) ([]metadata.Metadata, error)
	Get            func(client *gcorecloud.ServiceClient, id string, key string) metadata.MetadataResult
	Delete         func(client *gcorecloud.ServiceClient, id string, key string) metadata.MetadataActionResult
	CreateOrUpdate func(client *gcorecloud.ServiceClient, id string, opts map[string]string) metadata.MetadataActionResult
	Replace        func(client *gcorecloud.ServiceClient, id string, opts map[string]string) metadata.MetadataActionResult
}

// V1 calls the /metadata and /metadata_item endpoints of a resource
var V1 = Operations{
	ListAll:        v1.MetadataListAll,
	Get:            v1.MetadataGet,
	Delete:         v1.MetadataDelete,
	CreateOrUpdate: v1.MetadataCreateOrUpdate,
	Replace:        v1.MetadataReplace,
}

func NewMetadataListCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return V1.NewListCommand(cc, usage, argsUsage, errorText)
}

func NewMetadataGetCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return V1.NewGetCommand(cc, usage, argsUsage, errorText)
}

func NewMetadataDeleteCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return V1.NewDeleteCommand(cc, usage, argsUsage, errorText)
}

func NewMetadataCreateCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return V1.NewCreateCommand(cc, usage, argsUsage, errorText)
}

func NewMetadataUpdateCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return V1.NewUpdateCommand(cc, usage, argsUsage, errorText)
}

func NewMetadataReplaceCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
	return V1.NewReplaceCommand(cc, usage, argsUsage, errorText)
}

func showResults(c *cli.Context, data interface{}) error {
	if c.Bool("test") {
		err := json.NewEncoder(c.App.Writer).Encode(data)
		if err != nil {
			return err
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

func (ops Operations) NewListCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
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
			items, err := ops.ListAll(client, resourceID)
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			err = showResults(c, items)
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
}

func (ops Operations) NewGetCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
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

			item, err := ops.Get(client, resourceID, c.String("metadata")).Extract()
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			return showResults(c, item)
		}}
}

func (ops Operations) NewDeleteCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
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
			err = ops.Delete(client, resourceID, c.String("metadata")).ExtractErr()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func (ops Operations) NewCreateCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
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
			err = ops.CreateOrUpdate(client, resourceID, opts).ExtractErr()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func (ops Operations) NewUpdateCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
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
			err = ops.CreateOrUpdate(client, resourceID, opts).ExtractErr()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		}}
}

func (ops Operations) NewReplaceCommand(cc ClientConstructor, usage string, argsUsage string, errorText string) *cli.Command {
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
			err = ops.Replace(client, resourceID, opts).ExtractErr()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		}}
}
