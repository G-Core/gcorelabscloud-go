package clusters

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("")
}

func resourceURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return c.ServiceURL(clusterName)
}

func resourceActionURL(c *gcorecloud.ServiceClient, clusterName, action string) string {
	return c.ServiceURL(clusterName, action)
}

func checkLimitsURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("check_limits")
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return resourceURL(c, clusterName)
}

func deleteURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return resourceURL(c, clusterName)
}

func certificatesURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return resourceActionURL(c, clusterName, "certificates")
}

func configURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return resourceActionURL(c, clusterName, "config")
}

func instancesURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return resourceActionURL(c, clusterName, "instances")
}

func upgradeURL(c *gcorecloud.ServiceClient, clusterName string) string {
	return resourceActionURL(c, clusterName, "upgrade")
}

func versionsURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("k8s", "versions")
}
