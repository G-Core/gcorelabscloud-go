package client

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/common"
	"github.com/urfave/cli/v2"
)

func NewPostgresClientV1(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "dbaas/postgres/clusters", "v1")
}
