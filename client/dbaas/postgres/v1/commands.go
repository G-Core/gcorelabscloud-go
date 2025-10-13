package postgres

import (
	"github.com/G-Core/gcorelabscloud-go/client/dbaas/postgres/v1/clusters"
	"github.com/urfave/cli/v2"
)

// Commands returns the aggregated list of Managed PostgreSQL commands
var Commands = cli.Command{
	Name:        "postgres",
	Usage:       "Manage PostgreSQL databases",
	Description: "Parent command for Managed PostgreSQL operations",
	Category:    "dbaas",
	Subcommands: []*cli.Command{
		&clusters.Commands,
	},
}
