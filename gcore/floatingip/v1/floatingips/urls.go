package floatingips

import (
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

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func deleteURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func assignURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "assign")
}

func unAssignURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "unassign")
}
