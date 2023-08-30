package k8s

import (
	"github.com/G-Core/gcorelabscloud-go/client/k8s/v2/client"
	"github.com/G-Core/gcorelabscloud-go/client/utils"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v2/clusters"
	"github.com/urfave/cli/v2"
)

// var (
// 	clusterNameText = "cluster_name is mandatory argument"
// )

var clusterVersionsSubCommand = cli.Command{
	Name:     "versions",
	Usage:    "List available k8s cluster versions",
	Category: "cluster",
	Action: func(c *cli.Context) error {
		client, err := client.NewK8sClientV2(c)
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		results, err := clusters.VersionsAll(client)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))

		return nil
	},
}

var Commands = cli.Command{
	Name:  "cluster",
	Usage: "Gcloud k8s cluster commands",
	Subcommands: []*cli.Command{
		// &clusterCreateSubCommand,
		// &clusterListSubCommand,
		// &clusterGetSubCommand,
		// &clusterDeleteSubCommand,
		// &clusterCertificateSubCommand,
		// &clusterConfigSubCommand,
		// &clusterInstancesSubCommand,
		// &clusterUpgradeSubCommand,
		&clusterVersionsSubCommand,
		// &poolCommands,
	},
}
