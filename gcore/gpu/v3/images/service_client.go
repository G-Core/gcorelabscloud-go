package images

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/common"
	"github.com/urfave/cli/v2"
)

// NewGPUImageClientV3 creates a new GPU image client
func NewGPUImageClientV3(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "gpu/images", "v3")
}
