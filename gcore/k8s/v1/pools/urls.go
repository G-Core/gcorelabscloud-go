package pools

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceURL(c *gcorecloud.ServiceClient, clusterID, id string) string {
	return c.ServiceURL(clusterID, "pools", id)
}

func resourceActionURL(c *gcorecloud.ServiceClient, clusterID, id, action string) string {
	return c.ServiceURL(clusterID, "pools", id, action)
}

func rootURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return c.ServiceURL(clusterID, "pools")
}

func getURL(c *gcorecloud.ServiceClient, clusterID string, id string) string {
	return resourceURL(c, clusterID, id)
}

func listURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return rootURL(c, clusterID)
}

func createURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return rootURL(c, clusterID)
}

func updateURL(c *gcorecloud.ServiceClient, clusterID string, id string) string {
	return resourceURL(c, clusterID, id)
}

func deleteURL(c *gcorecloud.ServiceClient, clusterID string, id string) string {
	return resourceURL(c, clusterID, id)
}

func instancesURL(c *gcorecloud.ServiceClient, clusterID string, id string) string {
	return resourceActionURL(c, clusterID, id, "instances")
}

func volumesURL(c *gcorecloud.ServiceClient, clusterID string, id string) string {
	return resourceActionURL(c, clusterID, id, "volumes")
}
