package flags

import (
	"fmt"
	"strconv"

	"github.com/G-Core/gcorelabscloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
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
			Enum:    []string{"token", "password"},
			Default: "token",
		},
		Hidden: true,
		Usage:  "client type as token or password",
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

var passwordFlags = []cli.Flag{
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

func buildTokenClientFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, commonFlags...)
	flags = append(flags, tokenFlags...)
	return flags
}

func buildPasswordClientFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, commonFlags...)
	flags = append(flags, passwordFlags...)
	return flags
}

var TokenClientFlags = buildTokenClientFlags()
var PasswordClientFlags = buildPasswordClientFlags()

var TokenClientHelpText = `
   Environment variables example:

   GCLOUD_API_URL=
   GCLOUD_API_VERSION=v1
   GCLOUD_ACCESS_TOKEN=
   GCLOUD_REFRESH_TOKEN=
   GCLOUD_REGION=
   GCLOUD_PROJECT=
`

var PasswordClientHelpText = `
   Environment variables example:

   GCLOUD_AUTH_URL=
   GCLOUD_API_URL=
   GCLOUD_API_VERSION=v1
   GCLOUD_USERNAME=
   GCLOUD_PASSWORD=
   GCLOUD_REGION=
   GCLOUD_PROJECT=
`

func AddFlags(commands []*cli.Command, flags ...cli.Flag) {
	for _, cmd := range commands {
		sunCommands := cmd.Subcommands
		if len(sunCommands) != 0 {
			AddFlags(sunCommands, flags...)
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
