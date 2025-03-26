package volumes

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

const (
	gpuVirtualPath = "gpu/virtual"
	clustersPath   = "clusters"
	volumesPath    = "volumes"
)

// listURL returns URL for listing GPU virtual cluster volumes
func listURL(c *gcorecloud.ServiceClient, projectID int, regionID int, clusterID string) string {
	return c.ServiceURL(gpuVirtualPath, fmt.Sprintf("%d", projectID), fmt.Sprintf("%d", regionID), clustersPath, clusterID, volumesPath)
}
