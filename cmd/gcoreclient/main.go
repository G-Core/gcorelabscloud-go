package main

import (
	"fmt"
	"os"

	"github.com/G-Core/gcorelabscloud-go/client/flags"
	"github.com/G-Core/gcorelabscloud-go/client/flavors/v1/flavors"
	"github.com/G-Core/gcorelabscloud-go/client/floatingips/v1/floatingips"
	"github.com/G-Core/gcorelabscloud-go/client/heat"
	"github.com/G-Core/gcorelabscloud-go/client/images/v1/images"
	"github.com/G-Core/gcorelabscloud-go/client/instances/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/client/k8s/v1/k8s"
	"github.com/G-Core/gcorelabscloud-go/client/keypairs/v2/keypairs"
	"github.com/G-Core/gcorelabscloud-go/client/keystones/v1/keystones"
	"github.com/G-Core/gcorelabscloud-go/client/l7policies/v1/l7policies"
	"github.com/G-Core/gcorelabscloud-go/client/limits/v1/limits"
	"github.com/G-Core/gcorelabscloud-go/client/loadbalancers/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/client/networks/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/client/ports/v1/ports"
	"github.com/G-Core/gcorelabscloud-go/client/projects/v1/projects"
	"github.com/G-Core/gcorelabscloud-go/client/quotas/v1/quotas"
	"github.com/G-Core/gcorelabscloud-go/client/regions/v1/regions"
	"github.com/G-Core/gcorelabscloud-go/client/reservedfixedips/v1/reservedfixedips"
	"github.com/G-Core/gcorelabscloud-go/client/routers/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/client/securitygroups/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/client/servergroups/v1/servergroups"
	"github.com/G-Core/gcorelabscloud-go/client/snapshots/v1/snapshots"
	"github.com/G-Core/gcorelabscloud-go/client/subnets/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/client/tasks/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/client/volumes/v1/volumes"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var AppVersion = "v0.2.11"

var commands = []*cli.Command{
	&networks.NetworkCommands,
	&tasks.TaskCommands,
	&keypairs.KeypairCommands,
	&volumes.VolumeCommands,
	&subnets.SubnetCommands,
	&flavors.FlavorCommands,
	&loadbalancers.LoadBalancerCommands,
	&instances.InstanceCommands,
	&heat.HeatsCommand,
	&securitygroups.SecurityGroupCommands,
	&floatingips.FloatingIPCommands,
	&ports.PortCommands,
	&snapshots.SnapshotCommands,
	&images.ImageCommands,
	&regions.RegionCommands,
	&projects.ProjectCommands,
	&keystones.KeystoneCommands,
	&quotas.QuotaCommands,
	&limits.LimitCommands,
	&k8s.ClusterCommands,
	&k8s.ClusterPoolCommands,
	&l7policies.L7PolicyCommands,
	&routers.RouterCommands,
	&reservedfixedips.ReservedFixedIPCommands,
	&servergroups.ServerGroupsCommands,
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
		},
		flags: nil,
		usage: mainClientUsage,
	}
}

func main() {
	flags.AddOutputFlags(commands)
	clientCommands := buildClientCommands(commands)
	app := cli.NewApp()
	app.Version = AppVersion
	app.EnableBashCompletion = true
	app.Commands = clientCommands.commands
	if clientCommands.flags != nil {
		app.Flags = clientCommands.flags
	}
	if len(clientCommands.usage) > 0 {
		app.Usage = clientCommands.usage
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Errorf("Cannot initialize application: %+v", err)
		os.Exit(1)
	}
}
