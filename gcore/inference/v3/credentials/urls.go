package credentials

import (
	"strconv"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, name string) string {
	return c.BaseServiceURL("inference", strconv.Itoa(c.ProjectID), "registry_credentials", name)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("inference", strconv.Itoa(c.ProjectID), "registry_credentials")
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *gcorecloud.ServiceClient, name string) string {
	return resourceURL(c, name)
}

func deleteURL(c *gcorecloud.ServiceClient, name string) string {
	return resourceURL(c, name)
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *gcorecloud.ServiceClient, name string) string {
	return resourceURL(c, name)
}
