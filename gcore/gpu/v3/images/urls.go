package images

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

const (
	imagesPath = "images"
)

// ImagesURL returns URL for GPU images operations
func ImagesURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL(imagesPath)
}

// ImageURL returns URL for specific GPU image operations
func ImageURL(c *gcorecloud.ServiceClient, imageID string) string {
	return c.ServiceURL(imagesPath, imageID)
}
