package client

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/common"

	"github.com/urfave/cli/v2"
)

func NewSnapshotClientV1(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "snapshots", "v1")
}
