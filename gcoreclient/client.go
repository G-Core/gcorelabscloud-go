package main

import (
	"fmt"
	"os"

	"github.com/G-Core/gcorelabscloud-go/gcoreclient/k8s"

	"github.com/G-Core/gcorelabscloud-go/gcoreclient/limits"

	"github.com/G-Core/gcorelabscloud-go/gcoreclient/flags"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/flavors"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/floatingips/floatingips"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/heat"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/images"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/instances"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/keypairs"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/keystones"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/loadbalancers/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/magnum"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/networks"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/projects"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/quotas"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/regions"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/securitygroups/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/snapshots"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/volumes"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	&networks.NetworkCommands,
	&tasks.TaskCommands,
	&keypairs.KeypairCommands,
	&volumes.VolumeCommands,
	&subnets.SubnetCommands,
	&flavors.FlavorCommands,
	&loadbalancers.LoadBalancerCommands,
	&instances.InstanceCommands,
	&magnum.MagnumsCommand,
	&heat.HeatsCommand,
	&securitygroups.SecurityGroupCommands,
	&floatingips.FloatingIPCommands,
	&snapshots.SnapshotCommands,
	&images.ImageCommands,
	&regions.RegionCommands,
	&projects.ProjectCommands,
	&keystones.KeystoneCommands,
	&quotas.QuotaCommands,
	&limits.LimitCommands,
	&k8s.ClusterCommands,
	&k8s.ClusterPoolCommands,
}

func buildClientCommands(commands []*cli.Command) ([]*cli.Command, []cli.Flag, string) {
	clientType := os.Getenv("GCLOUD_CLIENT_TYPE")
	tokenClientUsage := fmt.Sprintf("GCloud API client\n%s", flags.TokenClientHelpText)
	passwordClientUsage := fmt.Sprintf("GCloud API client\n%s", flags.PasswordClientHelpText)
	if clientType == "client" {
		return commands, flags.TokenClientFlags, tokenClientUsage
	} else if clientType == "password" {
		return commands, flags.PasswordClientFlags, passwordClientUsage
	}
	return []*cli.Command{
		{
			Name:        "token",
			Aliases:     nil,
			Usage:       tokenClientUsage,
			Subcommands: commands,
			Flags:       flags.TokenClientFlags,
		},
		{
			Name:  "password",
			Usage: passwordClientUsage,
			Flags: flags.PasswordClientFlags,
			Before: func(c *cli.Context) error {
				return c.Set("client-type", "password")
			},
			Subcommands: commands,
		},
	}, nil, ""
}

func main() {

	flags.AddOutputFlags(commands)

	commands, appFlags, usage := buildClientCommands(commands)
	app := cli.NewApp()
	app.Version = "v0.2.7"
	app.EnableBashCompletion = true
	app.Commands = commands
	if appFlags != nil {
		app.Flags = appFlags
	}
	if len(usage) > 0 {
		app.Usage = usage
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Errorf("Cannot initialize application: %+v", err)
		os.Exit(1)
	}
}
