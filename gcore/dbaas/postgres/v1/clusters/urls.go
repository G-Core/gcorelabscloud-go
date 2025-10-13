package clusters

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

// ClusterURL returns the URL for specific PostgreSQL cluster operations
func ClusterURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return c.ServiceURL(clusterName)
}
