package clusters

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

//const (
//	clustersPath = "clusters"
//)

//// ClustersURL returns the URL for PostgreSQL clusters operations
//func ClustersURL(c *gcorecloud.ServiceClient) string {
//	return c.ServiceURL()
//}

// ClusterURL returns the URL for specific PostgreSQL cluster operations
func ClusterURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return c.ServiceURL(clusterName)
}
