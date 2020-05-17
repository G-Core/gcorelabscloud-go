package heat

import (
	"github.com/G-Core/gcorelabscloud-go/client/heat/v1/resources"
	"github.com/G-Core/gcorelabscloud-go/client/heat/v1/stacks"
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
