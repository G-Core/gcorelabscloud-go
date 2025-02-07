package flavors

import (
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, name string) string {
	return c.BaseServiceURL("inference", "flavors", name)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("inference", "flavors")
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *gcorecloud.ServiceClient, name string) string {
	return resourceURL(c, name)
}
