package quotas

import (
	"strconv"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id int) string {
	return c.BaseServiceURL("client_quotas", strconv.Itoa(id))
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.BaseServiceURL("client_quotas")
}

func replaceURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func getURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}

func quotaURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func updateURL(c *gcorecloud.ServiceClient, id int) string {
	return resourceURL(c, id)
}
