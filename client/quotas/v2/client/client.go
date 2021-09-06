package client

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/common"
	"github.com/urfave/cli/v2"
)

func NewQuotaClientV2(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "quotas", "v2")
}
