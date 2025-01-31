package images

import (
	"fmt"
)

// resourcePath is a base path for GPU image endpoints
func resourcePath(projectID, regionID int) string {
	return fmt.Sprintf("/v3/gpu/%d/%d", projectID, regionID)
}

// uploadBaremetalURL returns URL for uploading baremetal GPU images
func uploadBaremetalURL(projectID, regionID int) string {
	return resourcePath(projectID, regionID) + "/baremetal/images"
}

// uploadVirtualURL returns URL for uploading virtual GPU images
func uploadVirtualURL(projectID, regionID int) string {
	return resourcePath(projectID, regionID) + "/virtual/images"
}
