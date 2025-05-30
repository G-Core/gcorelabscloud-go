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

// ClusterServerURL returns URL for accessing a single server in a specific GPU cluster
func ClusterServerURL(c *gcorecloud.ServiceClient, clusterID, serverID string) string {
	return c.ServiceURL(clustersPath, clusterID, serversPath, serverID)
}
