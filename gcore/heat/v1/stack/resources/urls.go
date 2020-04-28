package resources

import (
	"github.com/G-Core/gcorelabscloud-go"
)

func resourceActionURL(c *gcorecloud.ServiceClient, stackID, resourceName, action string) string {
	return c.ServiceURL("stacks", stackID, "resources", resourceName, action)
}

func resourceURL(c *gcorecloud.ServiceClient, stackID, resourceName string) string {
	return c.ServiceURL("stacks", stackID, "resources", resourceName)
}

func rootURL(c *gcorecloud.ServiceClient, stackID string) string {
	return c.ServiceURL("stacks", stackID, "resources")
}

func MetadataURL(c *gcorecloud.ServiceClient, stackID, resourceName string) string {
	return resourceActionURL(c, stackID, resourceName, "metadata")
}

func SignalURL(c *gcorecloud.ServiceClient, stackID, resourceName string) string {
	return resourceActionURL(c, stackID, resourceName, "signal")
}

func listURL(c *gcorecloud.ServiceClient, stackID string) string {
	return rootURL(c, stackID)
}

func getURL(c *gcorecloud.ServiceClient, stackID, resourceName string) string {
	return resourceURL(c, stackID, resourceName)
}

func markUnhealthyURL(c *gcorecloud.ServiceClient, stackID, resourceName string) string {
	return resourceURL(c, stackID, resourceName)
}
