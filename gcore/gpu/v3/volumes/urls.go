package volumes

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

const (
	clustersPath = "clusters"
	volumesPath  = "volumes"
)

// listURL returns URL for listing GPU virtual cluster volumes
func listURL(c *gcorecloud.ServiceClient, projectID int, regionID int, clusterID string) string {
	return c.ServiceURL(clustersPath, clusterID, volumesPath)
}
