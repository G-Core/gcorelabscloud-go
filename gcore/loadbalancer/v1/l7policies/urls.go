package l7policies

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
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

func rulesrootURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "rules")
}

func rulescreateURL(c *gcorecloud.ServiceClient, id string) string {
	return rulesrootURL(c, id)
}

func rulesgetURL(c *gcorecloud.ServiceClient, plid string, rlid string) string {
	return c.ServiceURL(plid, "rules", rlid)
}
