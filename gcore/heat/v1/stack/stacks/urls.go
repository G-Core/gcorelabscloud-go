package stacks

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, stackID string) string {
	return c.ServiceURL("stacks", stackID)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("stacks")
}

func getURL(c *gcorecloud.ServiceClient, stackID string) string {
	return resourceURL(c, stackID)
}

func updateURL(c *gcorecloud.ServiceClient, stackID string) string {
	return resourceURL(c, stackID)
}

func deleteURL(c *gcorecloud.ServiceClient, stackID string) string {
	return resourceURL(c, stackID)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}
