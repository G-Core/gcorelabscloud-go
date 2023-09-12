package aiimages

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func listAIImagesURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL()
}
