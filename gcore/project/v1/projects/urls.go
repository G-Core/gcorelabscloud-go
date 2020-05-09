package projects

import (
	"strconv"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id int) string {
	return c.BaseServiceURL("projects", strconv.Itoa(id))
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("projects")
}

func getURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func deleteURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}
