package client

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/common"

	"github.com/urfave/cli/v2"
)

func NewFileShareClientV1(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "file_shares", "v1")
}

func NewFileShareClientV3(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "file_shares", "v3")
}
