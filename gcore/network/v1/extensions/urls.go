package extensions

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL("extensions", id)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("extensions")
}

func getURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}
