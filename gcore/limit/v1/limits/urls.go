package limits

import (
	"strconv"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id int) string {
	return c.BaseServiceURL("limits_request", strconv.Itoa(id))
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("limits_request")
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func getURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func deleteURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func statusURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}
