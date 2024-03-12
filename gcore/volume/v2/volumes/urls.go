package volumes

import gcorecloud "github.com/G-Core/gcorelabscloud-go"

func resourceActionURL(c *gcorecloud.ServiceClient, id, action string) string {
	return c.ServiceURL(id, action)
}

func attachURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "attach")
}

func detachURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "detach")
}
