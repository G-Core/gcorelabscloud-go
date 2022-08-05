package networks

import (
	"fmt"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}

func resourceActionURL(c *gcorecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL()
}

func getURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func listInstancePortsURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "ports")
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func metadataURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "metadata")
}
func metadataItemURL(c *gcorecloud.ServiceClient, id string, key string) string {
	return resourceActionURL(c, id, fmt.Sprintf("metadata_item?key=%s", key))
}
