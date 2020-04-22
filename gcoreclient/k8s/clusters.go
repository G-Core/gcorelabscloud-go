package k8s

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/k8s/v1/clusters"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var (
	clusterIDText = "cluster_id is mandatory argument"

//	clusterUpdateTypes     = types.ClusterUpdateOperation("").StringList()
//	k8sClusterVersionTypes = types.K8sClusterVersion("").StringList()
)

var clusterListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Magnum list clusters",
	Category: "cluster",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "k8s", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.ListAll(client, clusters.ListOpts{})
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var clusterGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Magnum get cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstStringArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "k8s", "", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := clusters.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var ClusterCommands = cli.Command{
	Name:  "cluster",
	Usage: "k8s cluster commands",
	Subcommands: []*cli.Command{
		&clusterListSubCommand,
		&clusterGetSubCommand,
	},
}
