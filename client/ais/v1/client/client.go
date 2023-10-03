package client

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/common"

	"github.com/urfave/cli/v2"
)

func NewAIClusterClientV1(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "ai/clusters", "v1")
}

func NewAIImageClientV1(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "ai/images", "v1")
}
func NewAIFlavorClientV1(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "ai/flavors", "v1")
}




