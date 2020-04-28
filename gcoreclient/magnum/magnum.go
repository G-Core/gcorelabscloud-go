package magnum

import (
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/magnum/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/magnum/nodegroups"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/magnum/templates"
	"github.com/urfave/cli/v2"
)

var MagnumsCommand = cli.Command{
	Name:  "magnum",
	Usage: "Gcloud Magnum API",
	Subcommands: []*cli.Command{
		&clusters.ClusterCommands,
		&templates.ClusterTemplatesCommands,
		&nodegroups.ClusterNodeGroupCommands,
	},
}
