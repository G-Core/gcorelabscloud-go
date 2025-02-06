package images

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// Common path components
const (
	imagesPath = "images"
)

// UploadBaremetalURL returns URL for uploading baremetal GPU images
func ImageURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL(imagesPath)
}
