package main

import (
	"fmt"
	"os"

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
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/securitygroups/securitygrouprules"
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
	&securitygrouprules.SecurityGroupRuleCommands,
	&floatingips.FloatingIPCommands,
	&snapshots.SnapshotCommands,
	&images.ImageCommands,
	&regions.RegionCommands,
	&projects.ProjectCommands,
	&keystones.KeystoneCommands,
	&quotas.QuotaCommands,
	&limits.LimitCommands,
}

func buildClientCommands(commands []*cli.Command) []*cli.Command {
	clientType := os.Getenv("GCLOUD_CLIENT_TYPE")
	if clientType == "client" {
		flags.AddFlags(commands, flags.TokenClientFlags...)
		return commands
	} else if clientType == "password" {
		flags.AddFlags(commands, flags.PasswordClientFlags...)
		return commands
	}
	return []*cli.Command{
		{
			Name:        "token",
			Aliases:     nil,
			Usage:       fmt.Sprintf("GCloud API client\n%s", flags.TokenClientHelpText),
			Subcommands: commands,
			Flags:       flags.TokenClientFlags,
		},
		{
			Name:  "password",
			Usage: fmt.Sprintf("GCloud API client\n%s", flags.PasswordClientHelpText),
			Flags: flags.PasswordClientFlags,
			Before: func(c *cli.Context) error {
				return c.Set("client-type", "password")
			},
			Subcommands: commands,
		},
	}
}

func main() {

	flags.AddOutputFlags(commands)

	app := cli.NewApp()
	app.Version = "v0.2.7"
	app.EnableBashCompletion = true
	app.Commands = buildClientCommands(commands)
	err := app.Run(os.Args)
	if err != nil {
		logrus.Errorf("Cannot initialize application: %+v", err)
		os.Exit(1)
	}
}
