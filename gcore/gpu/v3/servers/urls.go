package servers

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

const (
	clustersPath = "clusters"
	serversPath  = "servers"
)

// ClusterServersURL returns URL for listing servers in a specific GPU cluster
func ClusterServersURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return c.ServiceURL(clustersPath, clusterID, serversPath)
}
