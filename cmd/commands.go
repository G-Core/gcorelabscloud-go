package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/G-Core/gcorelabscloud-go/client/ais/v1/ais"
	"github.com/G-Core/gcorelabscloud-go/client/apitokens/v1/apitokens"
	"github.com/G-Core/gcorelabscloud-go/client/apptemplates/v1/apptemplates"
	"github.com/G-Core/gcorelabscloud-go/client/faas/v1/functions"
	"github.com/G-Core/gcorelabscloud-go/client/file_shares/v1/file_shares"
	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/flavors/v1/flavors"
	"github.com/G-Core/gcorelabscloud-go/client/floatingips/v1/floatingips"
	gpuimages "github.com/G-Core/gcorelabscloud-go/client/gpu/v3/images"
	"github.com/G-Core/gcorelabscloud-go/client/heat"
	"github.com/G-Core/gcorelabscloud-go/client/images/v1/images"
	"github.com/G-Core/gcorelabscloud-go/client/instances/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/client/k8s/v2/k8s"
	"github.com/G-Core/gcorelabscloud-go/client/keypairs/v2/keypairs"
	"github.com/G-Core/gcorelabscloud-go/client/keystones/v1/keystones"
	"github.com/G-Core/gcorelabscloud-go/client/l7policies/v1/l7policies"
	"github.com/G-Core/gcorelabscloud-go/client/lifecyclepolicy/v1/lifecyclepolicy"
	"github.com/G-Core/gcorelabscloud-go/client/limits/v2/limits"
	"github.com/G-Core/gcorelabscloud-go/client/loadbalancers/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/client/networks/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/client/ports/v1/ports"
	"github.com/G-Core/gcorelabscloud-go/client/projects/v1/projects"
	"github.com/G-Core/gcorelabscloud-go/client/quotas/v2/quotas"
	"github.com/G-Core/gcorelabscloud-go/client/regions/v1/regions"
	"github.com/G-Core/gcorelabscloud-go/client/regionsaccess/v1/regionsaccess"
	"github.com/G-Core/gcorelabscloud-go/client/reservedfixedips/v1/reservedfixedips"
	"github.com/G-Core/gcorelabscloud-go/client/routers/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/client/schedules/v1/schedules"
	"github.com/G-Core/gcorelabscloud-go/client/secrets/v1/secrets"
	"github.com/G-Core/gcorelabscloud-go/client/securitygroups/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/client/servergroups/v1/servergroups"
	"github.com/G-Core/gcorelabscloud-go/client/snapshots/v1/snapshots"
	"github.com/G-Core/gcorelabscloud-go/client/subnets/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/client/tasks/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/client/volumes/v1/volumes"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var AppVersion = "dev"

var commands = []*cli.Command{
	&networks.Commands,
	&tasks.Commands,
	&keypairs.Commands,
	&volumes.Commands,
	&subnets.Commands,
	&flavors.Commands,
	&loadbalancers.Commands,
	&instances.Commands,
	&heat.Commands,
	&securitygroups.Commands,
	&floatingips.Commands,
	&schedules.Commands,
	&ports.Commands,
	&snapshots.Commands,
	&images.Commands,
	&regions.Commands,
	&projects.Commands,
	&keystones.Commands,
	&quotas.Commands,
	&limits.Commands,
	&k8s.Commands,
	&l7policies.Commands,
	&routers.Commands,
	&reservedfixedips.Commands,
	&servergroups.Commands,
	&secrets.Commands,
	&lifecyclepolicy.Commands,
	&regionsaccess.Commands,
	&apptemplates.Commands,
	&apitokens.Commands,
	&file_shares.Commands,
	&ais.Commands,
	&functions.Commands,
	&gpuimages.Commands,
}

type clientCommands struct {
	commands []*cli.Command
	flags    []cli.Flag
	usage    string
}

func buildClientCommands(commands []*cli.Command) clientCommands {
	clientType := os.Getenv("GCLOUD_CLIENT_TYPE")
	tokenClientUsage := fmt.Sprintf("GCloud API client\n%s", flags.TokenClientHelpText)
	platformClientUsage := fmt.Sprintf("GCloud API client\n%s", flags.PlatformClientHelpText)
	apiTokenClientUsage := fmt.Sprintf("GCloud API client\n%s", flags.APITokenClientHelpText)
	switch clientType {
	case flags.ClientTypeToken:
		flags.ClientType = clientType
		return clientCommands{
			commands: commands,
			flags:    flags.TokenClientFlags,
			usage:    tokenClientUsage,
		}
	case flags.ClientTypePlatform:
		flags.ClientType = clientType
		return clientCommands{
			commands: commands,
			flags:    flags.PlatformClientFlags,
			usage:    platformClientUsage,
		}
	case flags.ClientTypeAPIToken:
		flags.ClientType = clientType
		return clientCommands{
			commands: commands,
			flags:    flags.APITokenClientFlags,
			usage:    apiTokenClientUsage,
		}
	}
	mainClientUsage := fmt.Sprintf("GCloud API client\n%s", flags.MainClientHelpText)
	return clientCommands{
		commands: []*cli.Command{
			{
				Name:        "token",
				Usage:       tokenClientUsage,
				Flags:       flags.TokenClientFlags,
				Subcommands: commands,
				Before: func(c *cli.Context) error {
					return c.Set("client-type", "token")
				},
			},
			{
				Name:        "platform",
				Usage:       platformClientUsage,
				Flags:       flags.PlatformClientFlags,
				Subcommands: commands,
				Before: func(c *cli.Context) error {
					return c.Set("client-type", "platform")
				},
			},
			{
				Name:        "api-token",
				Usage:       apiTokenClientUsage,
				Flags:       flags.APITokenClientFlags,
				Subcommands: commands,
				Before: func(c *cli.Context) error {
					return c.Set("client-type", "api-token")
				},
			},
		},
		flags: nil,
		usage: mainClientUsage,
	}
}

func NewApp(args []string) *cli.App {
	flags.AddOutputFlags(commands)
	clientCommands := buildClientCommands(commands)

	app := new(cli.App)
	app.Name = filepath.Base(args[0])
	app.HelpName = filepath.Base(args[0])
	app.Version = AppVersion
	app.EnableBashCompletion = true
	app.Commands = clientCommands.commands
	if clientCommands.flags != nil {
		app.Flags = clientCommands.flags
	}
	if len(clientCommands.usage) > 0 {
		app.Usage = clientCommands.usage
	}
	return app
}

func RunCommand(args []string) {
	app := NewApp(args)
	err := app.Run(args)
	if err != nil {
		logrus.Errorf("Cannot initialize application: %+v", err)
		os.Exit(1)
	}
}
