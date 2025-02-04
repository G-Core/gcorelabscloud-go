package images

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

// Common path components
const (
	imagesPath = "images"
)

// UploadBaremetalURL returns URL for uploading baremetal GPU images
func UploadBaremetalURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL(imagesPath)
}

// UploadVirtualURL returns URL for uploading virtual GPU images
func UploadVirtualURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL(imagesPath)
}
