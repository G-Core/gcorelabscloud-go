package gpu

import (
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/clusters"
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/flavors"
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/images"
	"github.com/G-Core/gcorelabscloud-go/client/gpu/v3/volumes"
	"github.com/urfave/cli/v2"
)

var baremetalCommands = cli.Command{
	Name:        "baremetal",
	Usage:       "Manage baremetal GPU resources",
	Description: "Commands for managing baremetal GPU resources",
	Subcommands: []*cli.Command{
		images.BaremetalCommands(),
		flavors.BaremetalCommands(),
		clusters.BaremetalCommands(),
	},
}

var virtualCommands = cli.Command{
	Name:        "virtual",
	Usage:       "Manage virtual GPU resources",
	Description: "Commands for managing virtual GPU resources",
	Subcommands: []*cli.Command{
		images.VirtualCommands(),
		flavors.VirtualCommands(),
		volumes.VirtualCommands(),
		clusters.VirtualCommands(),
	},
}

// Commands returns the aggregated list of GPU commands
var Commands = cli.Command{
	Name:        "gpu",
	Usage:       "Manage GPU resources",
	Description: "Parent command for GPU-related operations",
	Category:    "gpu",
	Subcommands: []*cli.Command{
		&baremetalCommands,
		&virtualCommands,
	},
}
