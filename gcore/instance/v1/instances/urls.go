package instances

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

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

func resourceActionURL(c *gcorecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func interfacesListURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "interfaces")
}
