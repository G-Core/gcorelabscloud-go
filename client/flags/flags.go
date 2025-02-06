package flags

import (
	"fmt"
	"strconv"

	"github.com/G-Core/gcorelabscloud-go/client/utils"

	"github.com/urfave/cli/v2"
)

const (
	ClientTypeToken    = "token"
	ClientTypePlatform = "platform"
	ClientTypeAPIToken = "api-token"
)

var (
	ClientType  string
	ClientTypes = []string{ClientTypeToken, ClientTypePlatform, ClientTypeAPIToken}
)

var commonFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "api-version",
		Usage:       "API version",
		DefaultText: "In case absent parameter it would take if from environ: GCLOUD_API_VERSION",
		Required:    false,
	},
	&cli.UintFlag{
		Name:        "region",
		DefaultText: "In case absent parameter it would take if from environ: GCLOUD_REGION",
		Usage:       "region ID",
		Required:    false,
	},
	&cli.UintFlag{
		Name:        "project",
		DefaultText: "In case absent parameter it would take if from environ: GCLOUD_PROJECT",
		Usage:       "project ID",
		Required:    false,
	},
	&cli.StringFlag{
		Name:        "api-url",
		Usage:       "Api base url",
		DefaultText: "In case absent parameter it would take if from environ: GCLOUD_API_URL",
		Required:    false,
	},
	&cli.GenericFlag{
		Name: "client-type",
		Value: &utils.EnumValue{
			Enum:        ClientTypes,
			Default:     "token",
			Destination: &ClientType,
		},
		Hidden: true,
		Usage:  "client type. Either use prepared token for gcloud API access, get an access token via gcloud platform or use an api token",
	},
}

var OutputFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:     "debug",
		Aliases:  []string{"d"},
		Usage:    "debug API requests",
		Required: false,
	},
	&cli.GenericFlag{
		Name:    "format",
		Aliases: []string{"f"},
		Value: &utils.EnumValue{
			Enum:    []string{"json", "table", "yaml"},
			Default: "json",
		},
		Usage: "output in json, table or yaml",
	},
}

var apiTokenFlag = []cli.Flag{
	&cli.StringFlag{
		Name:     "api-token",
		Usage:    "api token",
		Required: false,
	},
}

var tokenFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "access",
		Aliases:  []string{"at"},
		Usage:    "access token",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "refresh",
		Aliases:  []string{"rt"},
		Usage:    "refresh token",
		Required: false,
	},
}

var platformFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "auth-url",
		DefaultText: "In case absent parameter it would take if from environ: GCLOUD_AUTH_URL",
		Usage:       "Auth base url",
		Required:    false,
	},
	&cli.StringFlag{
		Name:     "username",
		Aliases:  []string{"u"},
		Usage:    "username",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "password",
		Aliases:  []string{"pass"},
		Usage:    "password",
		Required: false,
	},
}

var WaitCommandFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:     "wait",
		Aliases:  []string{"w"},
		Usage:    "Wait while command is being processed ",
		Value:    false,
		Required: false,
	},
	&cli.IntFlag{
		Name:     "wait-seconds",
		Usage:    "Required amount of time in seconds to wait while command is being processed",
		Value:    3600,
		Required: false,
	},
}

var retryFlags = []cli.Flag{
	&cli.IntFlag{
		Name:     "retry-amount",
		Aliases:  []string{"ra"},
		Usage:    "Amount of retries on a conflicting request",
		Value:    0,
		Required: false,
	},
	&cli.IntFlag{
		Name:     "retry-interval",
		Aliases:  []string{"ri"},
		Usage:    "Required amount of time in seconds to wait between retries",
		Value:    0,
		Required: false,
	},
}

func buildTokenClientFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, commonFlags...)
	flags = append(flags, tokenFlags...)
	return flags
}

func buildPlatformClientFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, commonFlags...)
	flags = append(flags, platformFlags...)
	return flags
}

func buildAPITokenClientFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, commonFlags...)
	flags = append(flags, apiTokenFlag...)
	return flags
}

func buildClientRequestFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, WaitCommandFlags...)
	flags = append(flags, retryFlags...)
	return flags
}

var TokenClientFlags = buildTokenClientFlags()
var PlatformClientFlags = buildPlatformClientFlags()
var APITokenClientFlags = buildAPITokenClientFlags()
var ClientRequestFlags = buildClientRequestFlags()

var TokenClientHelpText = `
   Environment variables example:

   GCLOUD_API_URL=
   GCLOUD_API_VERSION=v1
   GCLOUD_ACCESS_TOKEN=
   GCLOUD_REFRESH_TOKEN=
   GCLOUD_REGION=
   GCLOUD_PROJECT=
`

var PlatformClientHelpText = `
   Environment variables example:

   GCLOUD_AUTH_URL=
   GCLOUD_API_URL=
   GCLOUD_API_VERSION=v1
   GCLOUD_USERNAME=
   GCLOUD_PASSWORD=
   GCLOUD_REGION=
   GCLOUD_PROJECT=
`

var APITokenClientHelpText = `
   Environment variables example:

   GCLOUD_API_URL=
   GCLOUD_API_VERSION=v1
   GCLOUD_API_TOKEN=
   GCLOUD_REGION=
   GCLOUD_PROJECT=
`

var MainClientHelpText = `
   Environment variables example:
	
   GCLOUD_CLIENT_TYPE=[platform,token,api-token]
`

func AddFlags(commands []*cli.Command, flags ...cli.Flag) {
	for _, cmd := range commands {
		subCommands := cmd.Subcommands
		if len(subCommands) != 0 {
			AddFlags(subCommands, flags...)
		} else {
			cmd.Flags = append(cmd.Flags, flags...)
		}
	}
}

func AddOutputFlags(commands []*cli.Command) {
	AddFlags(commands, OutputFlags...)
}

func GetFirstStringArg(c *cli.Context, errorText string) (string, error) {
	arg := c.Args().First()
	if arg == "" {
		return "", cli.NewExitError(fmt.Errorf(errorText), 1)
	}
	return arg, nil
}

// GetNthStringArg returns the nth string argument, or an error if the arg is blank
func GetNthStringArg(c *cli.Context, errorText string, n int) (string, error) {
	arg := c.Args().Get(n)
	if arg == "" {
		return "", cli.Exit(fmt.Errorf(errorText), 1)
	}
	return arg, nil
}

func GetFirstIntArg(c *cli.Context, errorText string) (int, error) {
	arg := c.Args().First()
	if arg == "" {
		return 0, cli.NewExitError(fmt.Errorf(errorText), 1)
	}
	res, err := strconv.Atoi(arg)
	if err != nil {
		return 0, cli.NewExitError(err, 1)
	}
	return res, nil
}
