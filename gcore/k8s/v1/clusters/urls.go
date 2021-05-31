package clusters

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func versionsURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("k8s", "versions")
}

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("clusters")
}

func resourceActionURL(c *gcorecloud.ServiceClient, id, action string) string {
	return c.ServiceURL("clusters", id, action)
}

func configURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id, "config")
}

func resizeURL(c *gcorecloud.ServiceClient, clusterID, poolID string) string {
	return c.ServiceURL("pools", clusterID, poolID, "resize")
}

func upgradeURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "upgrade")
}

func instancesURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "instances")
}

func certificatesURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "certificates")
}

func volumesURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "volumes")
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
