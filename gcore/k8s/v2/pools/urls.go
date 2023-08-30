package pools

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func rootURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return c.ServiceURL(clusterName, "pools")
}

func resourceURL(c *gcorecloud.ServiceClient, clusterName, poolName string) string {
	return c.ServiceURL(clusterName, "pools", poolName)
}

func resourceActionURL(c *gcorecloud.ServiceClient, clusterName, poolName, action string) string {
	return c.ServiceURL(clusterName, "pools", poolName, action)
}

func listURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return rootURL(c, clusterName)
}

func createURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return rootURL(c, clusterName)
}

func getURL(c *gcorecloud.ServiceClient, clusterName, poolName string) string {
	return resourceURL(c, clusterName, poolName)
}

func updateURL(c *gcorecloud.ServiceClient, clusterName, poolName string) string {
	return resourceURL(c, clusterName, poolName)
}

func deleteURL(c *gcorecloud.ServiceClient, clusterName, poolName string) string {
	return resourceURL(c, clusterName, poolName)
}

func resizeURL(c *gcorecloud.ServiceClient, clusterName, poolName string) string {
	return resourceActionURL(c, clusterName, poolName, "resize")
}

func instancesURL(c *gcorecloud.ServiceClient, clusterName, poolName string) string {
	return resourceActionURL(c, clusterName, poolName, "instances")
}
