package metadata

import (
	"fmt"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func ResourceActionURL(c *gcorecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func MetadataURL(c *gcorecloud.ServiceClient, id string) string {
	return ResourceActionURL(c, id, "metadata")
}
func MetadataItemURL(c *gcorecloud.ServiceClient, id string, key string) string {
	return ResourceActionURL(c, id, fmt.Sprintf("metadata_item?key=%s", key))
}
