package main

import (
	"fmt"
	"os"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/k8s"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/limits"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flavors"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/floatingips/floatingips"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/heat"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/images"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/instances"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/keypairs"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/keystones"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/loadbalancers/loadbalancers"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/magnum"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/networks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/projects"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/quotas"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/regions"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/securitygroups/securitygroups"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/snapshots"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/subnets"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/volumes"
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
