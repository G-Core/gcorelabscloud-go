package client

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/client/common"

	"github.com/urfave/cli/v2"
)

func NewK8sClientV2(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "k8s", "v2")
}

func NewK8sClustersClientV2(c *cli.Context) (*gcorecloud.ServiceClient, error) {
	return common.BuildClient(c, "k8s/clusters", "v2")
}
