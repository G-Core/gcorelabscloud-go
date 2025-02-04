package client

import (
	"github.com/urfave/cli/v2"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/common"
)

// NewGPUBaremetalClientV3 creates a new GPU baremetal client
func NewGPUBaremetalClientV3(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "gpu/baremetal", "v3")
}

// NewGPUVirtualClientV3 creates a new GPU virtual client
func NewGPUVirtualClientV3(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "gpu/virtual", "v3")
}
