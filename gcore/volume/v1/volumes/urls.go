package volumes

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL(id)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL()
}

func getURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func resourceActionURL(c *gcorecloud.ServiceClient, id, action string) string {
	return c.ServiceURL(id, action)
}

func attachURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "attach")
}

func detachURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "detach")
}

func retypeURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "retype")
}

func extendURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "extend")
}

func revertURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "revert")
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
