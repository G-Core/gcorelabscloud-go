package lbpools

import (
	"github.com/G-Core/gcorelabscloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}

func resourceActionURL(c *gcorecloud.ServiceClient, id string, action string) string {
	return c.ServiceURL(id, action)
}

func resourceActionDetailURL(c *gcorecloud.ServiceClient, id string, action string, actorID string) string {
	return c.ServiceURL(id, action, actorID)
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

func updateURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func deleteURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func createMemberURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "member")
}

func deleteMemberURL(c *gcorecloud.ServiceClient, id string, memberID string) string {
	return resourceActionDetailURL(c, id, "member", memberID)
}
