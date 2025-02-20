package clusters

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

const (
	clustersPath = "clusters"
)

// ClustersURL returns URL for GPU clusters operations
func ClustersURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL(clustersPath)
}

// ClusterURL returns URL for specific GPU cluster operations
func ClusterURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return c.ServiceURL(clustersPath, clusterID)
}
