package clusters

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

const (
	clustersPath = "clusters"
)

// ClusterURL returns URL for specific GPU cluster operations
func ClusterURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return c.ServiceURL(clustersPath, clusterID)
}
