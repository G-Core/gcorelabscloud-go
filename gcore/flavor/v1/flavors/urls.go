package flavors

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL()
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}
