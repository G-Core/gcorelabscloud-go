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

func replaceURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func rulesRootURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id, "rules")
}

func rulesCreateURL(c *gcorecloud.ServiceClient, id string) string {
	return rulesRootURL(c, id)
}

func rulesGetURL(c *gcorecloud.ServiceClient, plid string, rlid string) string {
	return c.ServiceURL(plid, "rules", rlid)
}

func rulesListURL(c *gcorecloud.ServiceClient, id string) string {
	return rulesRootURL(c, id)
}

func rulesDeleteURL(c *gcorecloud.ServiceClient, plid string, rlid string) string {
	return c.ServiceURL(plid, "rules", rlid)
}

func rulesReplaceURL(c *gcorecloud.ServiceClient, plid string, rlid string) string {
	return c.ServiceURL(plid, "rules", rlid)
}
