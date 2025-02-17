package tasks

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.BaseServiceURL("tasks", id)
}

func getURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listActiveURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("active")
}
