package heat

import (
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/heat/resources"
	"github.com/G-Core/gcorelabscloud-go/gcoreclient/heat/stacks"
	"github.com/urfave/cli/v2"
)

var HeatsCommand = cli.Command{
	Name:  "heat",
	Usage: "Gcloud Heat API",
	Subcommands: []*cli.Command{
		&resources.ResourceCommands,
		&stacks.StackCommands,
	},
}
