package client

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/common"

	"github.com/urfave/cli/v2"
)

func NewFloatingIPClientV1(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "floatingips", "v1")
}

func NewAvailableFloatingIPClientV1(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "availablefloatingips", "v1")
}
